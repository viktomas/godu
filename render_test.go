package main

import (
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/mattn/go-runewidth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/viktomas/godu/commands"
	"github.com/viktomas/godu/files"
)

// setContent is one SetContent call made against the screen during a draw.
type setContent struct {
	x, y int
	r    rune
}

// recordingScreen wraps a tcell Screen and records every SetContent call.
// The final cell buffer cannot distinguish several rendering bugs (a wide
// rune's continuation column holds a plain space whether it was skipped or
// explicitly overwritten), so the draw *sequence* is the reliable oracle.
type recordingScreen struct {
	tcell.SimulationScreen
	sets []setContent
}

func (r *recordingScreen) SetContent(x, y int, mainc rune, combc []rune, style tcell.Style) {
	r.sets = append(r.sets, setContent{x, y, mainc})
	r.SimulationScreen.SetContent(x, y, mainc, combc, style)
}

// renderScreen drives the real printOptions pipeline (BoxLayout with the
// back/current/preview CellView panes and the status bar) against a headless
// tcell SimulationScreen, wrapped in a recordingScreen. It returns the
// recorded draw sequence and the resulting cell grid + styles. This
// exercises the same drawing code that runs against a real terminal, so it
// catches regressions that unit tests on GetCell alone cannot (blank
// screens, dropped or overwritten wide runes, ...).
func renderScreen(t *testing.T, state commands.State, w, h int) ([]setContent, [][]rune, [][]tcell.Style) {
	t.Helper()
	s := &recordingScreen{SimulationScreen: tcell.NewSimulationScreen("UTF-8")}
	require.NoError(t, s.Init())
	t.Cleanup(func() { s.Fini() })
	s.SetSize(w, h)
	printOptions(state, s)
	contents, cw, ch := s.GetContents()
	runes := make([][]rune, ch)
	styles := make([][]tcell.Style, ch)
	for y := 0; y < ch; y++ {
		runes[y] = make([]rune, cw)
		styles[y] = make([]tcell.Style, cw)
		for x := 0; x < cw; x++ {
			cell := contents[y*cw+x]
			if len(cell.Runes) > 0 {
				runes[y][x] = cell.Runes[0]
			} else {
				runes[y][x] = ' '
			}
			styles[y][x] = cell.Style
		}
	}
	return s.sets, runes, styles
}

// findRow returns the index of the first grid row containing substr, or -1.
func findRow(grid [][]rune, substr string) int {
	for y, row := range grid {
		if strings.Contains(string(row), substr) {
			return y
		}
	}
	return -1
}

// drawnRunes returns how often r was drawn at row y (any row if y < 0).
func drawnRunes(sets []setContent, r rune, y int) int {
	count := 0
	for _, s := range sets {
		if s.r == r && (y < 0 || s.y == y) {
			count++
		}
	}
	return count
}

func emptyMarked() map[*files.File]struct{} {
	return map[*files.File]struct{}{}
}

func TestRenderDrawsFolderTreeAndStatusBar(t *testing.T) {
	root := files.NewTestFolder(".",
		files.NewTestFile("big.bin", 5*1024*1024),
		files.NewTestFolder("docs", files.NewTestFile("manual.pdf", 2*1024*1024)),
	)
	state := commands.State{Folder: root, Selected: 0, MarkedFiles: emptyMarked()}
	sets, grid, _ := renderScreen(t, state, 120, 12)

	if len(sets) == 0 {
		t.Fatal("nothing was drawn, screen would be blank")
	}
	treeRow := findRow(grid, "5M big.bin")
	if treeRow == -1 {
		t.Fatalf("folder tree did not render, screen rows:\n%s", screenshot(grid))
	}
	// The tree must occupy the widget area above the status bar, not be
	// pushed off-screen (a blank screen was godu's most visible regression).
	if treeRow >= 10 {
		t.Fatalf("tree rendered at row %d, expected above the status bar", treeRow)
	}
	if !strings.Contains(string(grid[11]), "Total size:") {
		t.Errorf("status bar missing or wrong, last row: %q", string(grid[11]))
	}
}

// TestRenderWritesEveryRuneOfEveryName guards against draw loops that drop
// content: tcell v2.8/v2.9 dropped wide-rune names from CellView panes
// entirely, and godu's zero-width bug let the rune after a '\n' in a
// filename go undrawn. Every non-space rune of every displayed name (plain
// ASCII, wide CJK, and after a zero-width rune) must be drawn at least once.
func TestRenderWritesEveryRuneOfEveryName(t *testing.T) {
	root := files.NewTestFolder(".",
		files.NewTestFile("new\nline.txt", 1),
		files.NewTestFolder("日本語フォルダ",
			files.NewTestFile("中文文件名字很长很长很长很长.txt", 1000),
		),
	)
	// Selecting the CJK folder puts its long-named file into the preview
	// pane, so all three names render at once.
	state := commands.State{Folder: root, Selected: 1, MarkedFiles: emptyMarked()}
	sets, grid, _ := renderScreen(t, state, 120, 12)

	names := []string{"new\nline.txt", "日本語フォルダ/", "中文文件名字很长很长很长很长.txt"}
	for _, name := range names {
		for _, r := range name {
			if r == ' ' {
				continue
			}
			if drawnRunes(sets, r, -1) == 0 {
				t.Errorf("rune %q of %q was never drawn, screen rows:\n%s", r, name, screenshot(grid))
			}
		}
	}
	// The rune following the zero-width '\n' must be drawn at a column the
	// '\n' does not occupy (they shared a column in the bug).
	newlineCol := -1
	for _, s := range sets {
		if s.r == '\n' {
			newlineCol = s.x
		}
	}
	if newlineCol == -1 {
		t.Fatal("the '\\n' rune itself was never drawn")
	}
	for _, s := range sets {
		if s.r == 'l' && s.x == newlineCol {
			t.Errorf("'l' after '\\n' drawn into the same column %d — the rune following a zero-width rune must get its own column", newlineCol)
		}
	}
}

// TestRenderDoesNotOverwriteWideRuneContinuation guards against the tcell
// v2.10+ CellView regression: its draw loop asks the model for the second
// column of a wide rune and draws a space there, clobbering the glyph on
// the terminal. Once a wide rune is drawn, nothing may write into its
// continuation column afterwards.
func TestRenderDoesNotOverwriteWideRuneContinuation(t *testing.T) {
	root := files.NewTestFolder(".",
		files.NewTestFolder("日本語フォルダ",
			files.NewTestFile("中文文件名字很长很长很长很长.txt", 1000),
		),
	)
	// Selected line exercises the reverse-video style path too.
	state := commands.State{Folder: root, Selected: 0, MarkedFiles: emptyMarked()}
	sets, _, _ := renderScreen(t, state, 120, 12)

	for i, s := range sets {
		if runewidth.RuneWidth(s.r) <= 1 {
			continue
		}
		for _, later := range sets[i+1:] {
			if later.x == s.x+1 && later.y == s.y {
				t.Fatalf("rune %q at (%d,%d) had its continuation column overwritten by %q — wide glyphs will be clobbered on the terminal",
					s.r, s.x, s.y, later.r)
			}
		}
	}
}

// TestRenderZeroWidthRunesDoNotEatRunes covers filenames containing
// zero-width runes such as '\n'. Before the GetCell width fix the draw of
// such a line hung tcell <= v2.11 forever (blank screen at 100% CPU); the
// draw here completing at all is part of the assertion, and the rendered
// line must contain the runes after the newline.
func TestRenderZeroWidthRunesDoNotEatRunes(t *testing.T) {
	root := files.NewTestFolder(".",
		files.NewTestFile("new\nline.txt", 1),
	)
	state := commands.State{Folder: root, Selected: 0, MarkedFiles: emptyMarked()}
	_, grid, _ := renderScreen(t, state, 120, 12)

	row := findRow(grid, "new")
	if row == -1 {
		t.Fatalf("file with zero-width rune in its name did not render, screen rows:\n%s", screenshot(grid))
	}
	if !strings.Contains(string(grid[row]), "line.txt") {
		t.Errorf("runes after the zero-width rune were not drawn, row: %q", string(grid[row]))
	}
}

// TestVisualStateAlwaysAdvancesDrawCursor encodes the exact invariant that
// the pre-v2.12 tcell CellView draw loop relies on: after drawing a cell it
// advances the cursor by the returned width (post-increment x++ plus
// x += wid - 1). If GetCell ever returns width 0 again, that loop spins
// forever at 100% CPU and godu shows a blank screen. Every cell must
// advance the cursor by at least one column.
func TestVisualStateAlwaysAdvancesDrawCursor(t *testing.T) {
	root := files.NewTestFolder(".",
		files.NewTestFile("new\nline.txt", 1),
		files.NewTestFile("中文文件名字很长很长很长很长.txt", 1000),
	)
	state := commands.State{Folder: root, Selected: 0, MarkedFiles: emptyMarked()}
	vs := newVisualState(state, 5)
	xbound, ybound := vs.GetBounds()

	for y := 0; y < ybound; y++ {
		x, steps := 0, 0
		for x <= xbound {
			_, _, _, width := vs.GetCell(x, y)
			if width < 1 {
				t.Fatalf("GetCell(%d, %d) returned width %d; tcell's draw loop advances by the returned width and would hang forever", x, y, width)
			}
			x += width
			steps++
			if steps > xbound+2 {
				t.Fatalf("draw cursor did not cross row %d after %d steps", y, steps)
			}
		}
	}
}

func TestRenderMarkedFileIsHighlighted(t *testing.T) {
	marked := files.NewTestFile("precious.txt", 1024)
	root := files.NewTestFolder(".",
		marked,
		files.NewTestFile("other.txt", 2048),
	)
	state := commands.State{
		Folder:      root,
		Selected:    0,
		MarkedFiles: map[*files.File]struct{}{marked: {}},
	}
	_, grid, styles := renderScreen(t, state, 120, 12)

	row := findRow(grid, "precious.txt")
	if row == -1 {
		t.Fatalf("marked file not rendered, screen rows:\n%s", screenshot(grid))
	}
	if !strings.HasPrefix(strings.TrimSpace(string(grid[row])), "*") {
		t.Errorf("marked file missing '*' marker, row: %q", string(grid[row]))
	}
	fg, _, _ := styles[row][strings.Index(string(grid[row]), "precious")].Decompose()
	if fg != tcell.ColorGreen {
		t.Errorf("marked file not rendered in green, got %v", fg)
	}
	if !strings.Contains(string(grid[11]), "Selected size:") {
		t.Errorf("status bar missing selected size, last row: %q", string(grid[11]))
	}
}

func TestRenderShowsAdjacentPanes(t *testing.T) {
	sub := files.NewTestFolder("sub", files.NewTestFile("inner.txt", 3))
	docs := files.NewTestFolder("docs",
		files.NewTestFile("manual.pdf", 2*1024*1024),
		sub,
	)
	files.NewTestFolder(".",
		files.NewTestFile("big.bin", 5*1024*1024),
		docs,
	)
	// Standing inside "docs" with "sub" selected: the back pane shows the
	// parent folder, the middle pane shows "docs", the preview pane shows
	// the contents of the selected subfolder.
	state := commands.State{Folder: docs, Selected: 1, MarkedFiles: emptyMarked()}
	_, grid, _ := renderScreen(t, state, 120, 12)

	assert.NotEqual(t, -1, findRow(grid, "big.bin"), "back pane should show parent folder contents")
	assert.NotEqual(t, -1, findRow(grid, "sub/"), "middle pane should show current folder contents")
	assert.NotEqual(t, -1, findRow(grid, "inner.txt"), "preview pane should show selected folder contents")
}

// TestRenderUsesCellViewModel ensures the panes keep going through the
// views.CellView/CellModel contract — the layer where the tcell regressions
// lived — rather than some other drawing path.
func TestRenderUsesCellViewModel(t *testing.T) {
	root := files.NewTestFolder(".", files.NewTestFile("big.bin", 5*1024*1024))
	state := commands.State{Folder: root, Selected: 0, MarkedFiles: emptyMarked()}
	_, grid, _ := renderScreen(t, state, 120, 12)

	// newVisualState's model is only rendered via views.NewCellView in
	// printOptions; if the row is present the CellView path was exercised.
	if findRow(grid, "big.bin") == -1 {
		t.Fatalf("expected render through views.CellView, screen rows:\n%s", screenshot(grid))
	}
	assert.Implements(t, (*views.CellModel)(nil), newVisualState(state, 10))
}

// screenshot renders the grid as a string for failure messages.
func screenshot(grid [][]rune) string {
	var b strings.Builder
	for _, row := range grid {
		b.WriteString("|")
		b.WriteString(strings.TrimRight(string(row), " "))
		b.WriteString("|\n")
	}
	return b.String()
}
