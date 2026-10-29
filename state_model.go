package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/viktomas/godu/commands"
	"github.com/viktomas/godu/interactive"
)

type visualState struct {
	folders        []interactive.Line
	selected       int
	xbound, ybound int
	screenHeight   int
}

func newVisualState(state commands.State, screenHeight int) visualState {
	lines := interactive.ReportFolder(state.Folder, state.MarkedFiles)
	xbound := 0
	ybound := len(lines)
	for index, line := range lines {
		// GetCell addresses lines by display width (runewidth),
		// so the horizontal bound has to be measured in display columns
		// rather than in number of runes
		if width := runewidth.StringWidth(string(line.Text)); width-1 > xbound {
			xbound = width - 1
		}
		lines[index] = line
	}
	return visualState{lines, state.Selected, xbound, ybound, screenHeight}
}

func (vs visualState) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	style := tcell.StyleDefault
	// return empty cell if we are asking for a line that doesn't exist
	if y >= len(vs.folders) {
		return ' ', style, nil, 1
	}
	// For some reason tcell is asking for cells below the viewport, we will return empty cell
	if y > vs.screenHeight {
		return ' ', style, nil, 1
	}
	shiftedIndex := y
	if vs.selected > vs.screenHeight {
		// shifting the index enables displaying selected folders that would be otherwise hidden bellow the screen
		shiftedIndex += vs.selected - vs.screenHeight
	}
	if shiftedIndex == vs.selected {
		style = style.Reverse(true)
	}
	line := vs.folders[shiftedIndex]
	if line.IsMarked {
		style = style.Foreground(tcell.ColorGreen)
	}

	visualIndex := 0
	for _, r := range line.Text {
		width := runewidth.RuneWidth(r)
		if visualIndex == x {
			return r, style, nil, width
		}
		visualIndex += width
		if visualIndex > x {
			break
		}
	}
	return ' ', style, nil, 1
}
func (vs visualState) GetBounds() (int, int) {
	return vs.xbound, vs.ybound
}
func (visualState) SetCursor(int, int) {
}

func (visualState) GetCursor() (int, int, bool, bool) {
	return 0, 0, false, false
}
func (visualState) MoveCursor(offx, offy int) {

}
