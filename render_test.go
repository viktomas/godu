package main

import (
	"strings"
	"testing"
	"time"

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
//
// NOTE for tcell upgrades: from tcell v2.12 the views layer draws with
// Screen.Put instead of SetContent, and recordingScreen must then also
// override Put, or every test below fails with "nothing was drawn" — by
// design, so an empty recording can never pass silently. The pinned tcell
// v2.7.4 has no Put method to override.
type recordingScreen struct {
	tcell.SimulationScreen
	sets []setContent
}

func (r *recordingScreen) SetContent(x, y int, mainc rune, combc []rune, style tcell.Style) {
	r.sets = append(r.sets, setContent{x, y, mainc})
	r.SimulationScreen.SetContent(x, y, mainc, combc, style)
}

// renderScreen drives the real printOptions pipeline (BoxLayout with the
// back/current/preview CellView panes and the status bar) against a
// headless tcell SimulationScreen, wrapped in a recordingScreen. It
// returns the recorded draw sequence and the resulting cell grid + styles.
//
// A rendering hang (the zero-width rune regression blanked godu's screen
// at 100% CPU) would otherwise spin until Go's package timeout, so the
// draw runs under a short deadline instead.
func renderScreen(t *testing.T, state commands.State, w, h int) ([]setContent, [][]rune, [][]tcell.Style) {
	t.Helper()
	s := &recordingScreen{SimulationScreen: tcell.NewSimulationScreen("UTF-8")}
	require.NoError(t, s.Init())
	hung := false
	t.Cleanup(func() {
		if !hung {
			s.Fini()
		}
	})
	s.SetSize(w, h)

	done := make(chan struct{})
	go func() {
		printOptions(state, s)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		hung = true
		t.Fatal("printOptions did not return within 5s — rendering hung, likely the zero-width rune regression in GetCell")
	}

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

// requireRecorded fails the test if the recording is empty. On tcell >=
// v2.12 the views layer draws via Put instead of SetContent, so the
// recording is empty there and the failure message points at the missing
// Put override instead of a silent pass.
func requireRecorded(t *testing.T, sets []setContent) {
	t.Helper()
	if len(sets) == 0 {
		t.Fatal("nothing was drawn (blank screen), or the recording is empty — on tcell >= v2.12 the views layer draws via Put, which recordingScreen does not record")
	}
}

// requireWideRunes fails the test unless wide runes were actually drawn,
// so the wide-rune assertions cannot pass vacuously.
func requireWideRunes(t *testing.T, sets []setContent) {
	t.Helper()
	wide := 0
	for _, s := range sets {
		if runewidth.RuneWidth(s.r) > 1 {
			wide++
		}
	}
	if wide == 0 {
		t.Fatalf("no wide runes were drawn in %d recorded cells; the wide-rune assertions would pass vacuously", len(sets))
	}
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

// recordedRows rebuilds one string per row from the recorded SetContent
// calls (column order). Cells that were never written — a wide rune's
// continuation column in a correct render — simply do not appear, so a
// fully drawn name is a contiguous substring of its row.
func recordedRows(sets []setContent) map[int]string {
	rows := map[int][]rune{}
	for _, s := range sets {
		rows[s.y] = append(rows[s.y], s.r)
	}
	out := map[int]string{}
	for y, r := range rows {
		out[y] = string(r)
	}
	return out
}

// assertNameDrawn fails unless every rune of name is drawn, in order and
// contiguously, on a single row — catching runes dropped by the model or
// the draw loop (wide CJK runes as well as plain ASCII).
func assertNameDrawn(t *testing.T, sets []setContent, name string) {
	t.Helper()
	rows := recordedRows(sets)
	for _, row := range rows {
		if strings.Contains(row, name) {
			return
		}
	}
	fmtRows := make([]string, 0, len(rows))
	for _, row := range rows {
		fmtRows = append(fmtRows, row)
	}
	t.Errorf("name %q was not drawn contiguously on any row; recorded rows: %q", name, fmtRows)
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
	requireRecorded(t, sets)

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

// TestRenderWritesEveryRuneOfEveryName guards the model-side content bugs:
// godu's zero-width bug let the rune after a '\n' in a filename go undrawn,
// and a model or draw-loop regression can drop runes from any name. Every
// displayed name — plain ASCII, wide CJK, and runes following a zero-width
// rune — must be drawn contiguously, in order, on one row.
//
// (Names dropped only at the terminal-emission layer — the tcell v2.8-v2.11
// real-terminal corruption — are invisible here by construction, because
// SimulationScreen does not share tcell's emission code; those are caught
// by scripts/terminal-smoke.sh.)
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
	requireRecorded(t, sets)
	requireWideRunes(t, sets)

	assertNameDrawn(t, sets, "new\nline.txt")
	assertNameDrawn(t, sets, "日本語フォルダ/")
	assertNameDrawn(t, sets, "中文文件名字很长很长很长很长.txt")

	if findRow(grid, "new") == -1 {
		t.Fatalf("file with zero-width rune in its name did not render, screen rows:\n%s", screenshot(grid))
	}
}

// TestRenderDoesNotOverwriteWideRuneContinuation guards against a draw loop
// that writes into the second column of a wide rune (tcell v2.12+ did this
// by switching views/cellarea.go to a range loop, clobbering every second
// glyph). Once a wide rune is drawn, nothing may write into its
// continuation column afterwards.
//
// In principle a wide rune sitting exactly on a pane's right edge could be
// followed by the neighbouring pane legitimately writing the next column;
// the fixture's names end in ASCII so this cannot trigger here, and if it
// ever did, the display would genuinely be broken.
func TestRenderDoesNotOverwriteWideRuneContinuation(t *testing.T) {
	root := files.NewTestFolder(".",
		files.NewTestFolder("日本語フォルダ",
			files.NewTestFile("中文文件名字很长很长很长很长.txt", 1000),
		),
	)
	// Selected line exercises the reverse-video style path too.
	state := commands.State{Folder: root, Selected: 0, MarkedFiles: emptyMarked()}
	sets, _, _ := renderScreen(t, state, 120, 12)
	requireRecorded(t, sets)
	requireWideRunes(t, sets)

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
// zero-width runes such as '\n': the rendered line must contain the runes
// after the newline (before the GetCell width fix, the draw of such a line
// hung tcell <= v2.11 forever — the watchdog in renderScreen is part of
// this test's protection).
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
// tcell's CellView draw loop relies on: after drawing a cell it advances
// the cursor by the returned width. If GetCell ever returns width 0 again,
// that loop spins forever and godu shows a blank screen. (This is the
// direct, fast hang guard; TestRenderZeroWidthRunesDoNotEatRunes covers the
// same regression through the full pipeline.)
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
	sets, grid, styles := renderScreen(t, state, 120, 12)
	requireRecorded(t, sets)

	row := findRow(grid, "precious.txt")
	if row == -1 {
		t.Fatalf("marked file not rendered, screen rows:\n%s", screenshot(grid))
	}
	if !strings.HasPrefix(strings.TrimSpace(string(grid[row])), "*") {
		t.Errorf("marked file missing '*' marker, row: %q", string(grid[row]))
	}
	// Column of the first 'p' of the name (rune scan, not byte index).
	col := -1
	for x, r := range grid[row] {
		if r == 'p' {
			col = x
			break
		}
	}
	fg, _, _ := styles[row][col].Decompose()
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
	sets, grid, _ := renderScreen(t, state, 120, 12)
	requireRecorded(t, sets)

	assert.NotEqual(t, -1, findRow(grid, "big.bin"), "back pane should show parent folder contents")
	assert.NotEqual(t, -1, findRow(grid, "sub/"), "middle pane should show current folder contents")
	assert.NotEqual(t, -1, findRow(grid, "inner.txt"), "preview pane should show selected folder contents")
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

// Compile-time check that the model keeps satisfying the CellModel contract
// the render pipeline depends on.
var _ views.CellModel = visualState{}
