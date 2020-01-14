package main

import (
	"github.com/gdamore/tcell"
	"github.com/viktomas/godu/core"
	"github.com/viktomas/godu/interactive"
)

type visualState struct {
	folders        []interactive.Line
	selected       int
	xbound, ybound int
	screenHeight   int
}

func newVisualState(state core.State, screenHeight int) visualState {
	lines := interactive.ReportFolder(state.Folder, state.MarkedFiles)
	xbound := 0
	ybound := len(lines)
	for index, line := range lines {
		if len(line.Text)-1 > xbound {
			xbound = len(line.Text) - 1
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
	// this offset enables displaying selected folders that would be otherwise hidden bellow the screen
	heightOffset := 0
	if vs.selected > vs.screenHeight {
		heightOffset = vs.selected - vs.screenHeight
	}
	shiftedIndex := y + heightOffset
	if shiftedIndex == vs.selected {
		style = style.Reverse(true)
	}
	line := vs.folders[shiftedIndex]
	if line.IsMarked {
		style = style.Foreground(tcell.ColorGreen)
	}
	if x < len(vs.folders[shiftedIndex].Text) {
		return line.Text[x], style, nil, 1
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
