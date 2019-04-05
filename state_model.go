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
}

func newVisualState(state core.State) visualState {
	lines := interactive.ReportFolder(state.Folder, state.MarkedFiles)
	xbound := 0
	ybound := len(lines)
	for index, line := range lines {
		if len(line.Text)-1 > xbound {
			xbound = len(line.Text) - 1
		}
		lines[index] = line
	}
	return visualState{lines, state.Selected, xbound, ybound}
}

func (vs visualState) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	style := tcell.StyleDefault
	if y == vs.selected {
		style = style.Reverse(true)
	}
	if y < len(vs.folders) {
		line := vs.folders[y]
		if line.IsMarked {
			style = style.Foreground(tcell.ColorGreen)
		}
		if x < len(vs.folders[y].Text) {
			return line.Text[x], style, nil, 1
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
