package main

import (
	"sync"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"github.com/viktomas/godu/commands"
	"github.com/viktomas/godu/interactive"
)

func interactiveFolder(s tcell.Screen, states chan commands.State, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		state, more := <-states
		if !more {
			break
		}
		printOptions(state, s)
	}
}

func printOptions(state commands.State, s tcell.Screen) {
	s.Clear()
	outer := views.NewBoxLayout(views.Vertical)
	inner := views.NewBoxLayout(views.Horizontal)
	var back views.Widget
	var forth views.Widget

	// Subtract a row from screen height for the status bar
	_, screenHeight := s.Size()
	innerWidgetHeight := screenHeight - 2

	middle := views.NewCellView()

	middle.SetModel(newVisualState(state, innerWidgetHeight))
	backState, err := commands.GoBack{}.Execute(state)
	if err == nil {
		backCell := views.NewCellView()
		backCell.SetModel(newVisualState(backState, innerWidgetHeight))
		back = backCell
	} else {
		back = views.NewText()
	}
	forthState, err := commands.Enter{}.Execute(state)
	if err == nil {
		forthCell := views.NewCellView()
		forthCell.SetModel(newVisualState(forthState, innerWidgetHeight))
		forth = forthCell
	} else {
		forth = views.NewText()
	}

	statusBar := views.NewSimpleStyledTextBar()
	status := interactive.ReportStatus(state.Folder, &state.MarkedFiles)

	statusBar.SetLeft(status.Total)
	statusBar.SetRight(status.Selected)

	outer.SetView(s)
	outer.AddWidget(inner, 1.0)
	outer.AddWidget(statusBar, 0.0)
	inner.AddWidget(back, 0.33)
	inner.AddWidget(middle, 0.33)
	inner.AddWidget(forth, 0.33)
	outer.Draw()
	s.Show()
}
