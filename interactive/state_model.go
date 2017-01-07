package interactive

import (
	"github.com/gdamore/tcell"
	"github.com/viktomas/godu/core"
)

type VisualState struct {
	folders        [][]rune
	selected       int
	xbound, ybound int
}

func NewVisualState(state core.State) VisualState {
	lines := ReportTree(state.Folder)
	folders := make([][]rune, len(lines))
	xbound := 0
	ybound := len(folders)
	for index, line := range lines {
		runeLine := []rune(line)
		if len(runeLine)-1 > xbound {
			xbound = len(runeLine) - 1
		}
		folders[index] = runeLine
	}
	return VisualState{folders, state.Selected, xbound, ybound}
}

func (vs VisualState) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	style := tcell.StyleDefault
	if y == vs.selected {
		style = style.Reverse(true)
	}
	if y < len(vs.folders) && x < len(vs.folders[y]) {
		return vs.folders[y][x], style, nil, 1
	}
	return ' ', style, nil, 1
}
func (vs VisualState) GetBounds() (int, int) {
	return vs.xbound, vs.ybound
}
func (VisualState) SetCursor(int, int) {
}

func (VisualState) GetCursor() (int, int, bool, bool) {
	return 0, 0, false, false
}
func (VisualState) MoveCursor(offx, offy int) {

}
