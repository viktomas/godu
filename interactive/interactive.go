package interactive

import (
	"sync"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
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
	inner := views.NewBoxLayout(views.Horizontal)
	var back views.Widget
	var forth views.Widget

	middle := views.NewCellView()

	middle.SetModel(NewVisualState(state))
	backState, err := core.GoBack{}.Execute(state)
	if err == nil {
		backCell := views.NewCellView()
		backCell.SetModel(NewVisualState(backState))
		back = backCell
	} else {
		back = views.NewText()
	}
	forthState, err := core.Enter{}.Execute(state)
	if err == nil {
		forthCell := views.NewCellView()
		forthCell.SetModel(NewVisualState(forthState))
		forth = forthCell
	} else {
		forth = views.NewText()
	}

	inner.AddWidget(back, 0.33)
	inner.AddWidget(middle, 0.33)
	inner.AddWidget(forth, 0.33)
	inner.SetView(s)
	inner.Draw()
	s.Show()
}
