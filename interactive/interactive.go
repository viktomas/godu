package interactive

import (
	"github.com/gdamore/tcell"
	"github.com/viktomas/godu/core"
)

func InteractiveTree(tree *core.File, s tcell.Screen, commands chan core.Executer, quit chan struct{}, limit int64) {
	core.PruneTree(tree, limit)
	core.SortDesc(tree)
	state := core.State{
		Folder: tree,
	}
	printOptions(state, s)
	for {
		select {
		case command := <-commands:
			state, _ = command.Execute(state)
			printOptions(state, s)
		case <-quit:
			break
		}
	}
}

func printOptions(state core.State, s tcell.Screen) {
	style := tcell.StyleDefault
	s.Clear()
	lines := ReportTree(state.Folder)
	for y, line := range lines {
		//fmt.Println(line)

		for x, char := range line {
			s.SetContent(x, y, char, []rune{}, style)
		}
	}
	s.Show()
}
