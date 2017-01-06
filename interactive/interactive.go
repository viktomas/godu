package interactive

import (
	"sync"

	"github.com/gdamore/tcell"
	"github.com/viktomas/godu/core"
)

func InteractiveTree(tree *core.File, s tcell.Screen, commands chan core.Executer, wg *sync.WaitGroup, limit int64) {
	defer wg.Done()
	core.PruneTree(tree, limit)
	core.SortDesc(tree)
	state := core.State{
		Folder: tree,
	}
	printOptions(state, s)
	for {
		command, more := <-commands
		if !more {
			break
		}
		state, _ = command.Execute(state)
		printOptions(state, s)
	}
}

func printOptions(state core.State, s tcell.Screen) {
	s.Clear()
	lines := ReportTree(state.Folder)
	for y, line := range lines {
		style := tcell.StyleDefault
		if y == state.Selected {
			style = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
		}

		for x, char := range line {
			s.SetContent(x, y, char, []rune{}, style)
		}
	}
	s.Show()
}
