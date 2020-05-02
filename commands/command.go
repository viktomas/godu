package commands

import (
	"errors"

	"github.com/viktomas/godu/files"
)

// State represents system configuration after processing user input
type State struct {
	Folder      *files.File
	Selected    int
	history     map[*files.File]int // last cursor location in each folder
	MarkedFiles map[*files.File]struct{}
}

// Executer represents a user action triggered on a State
type Executer interface {
	Execute(State) (State, error)
}

// Enter is an action opening selected directory
type Enter struct{}

// GoBack is an action returning to parent directory
type GoBack struct{}

// Down is an action selecting next file in the list
type Down struct{}

// Up is an action selecting previous file in the list
type Up struct{}

// Mark is an action that saves current directory for later use
type Mark struct{}

func copyState(state State) State {
	return State{
		Folder:      state.Folder,
		history:     state.history,
		Selected:    state.Selected,
		MarkedFiles: state.MarkedFiles,
	}
}

func (d Down) Execute(oldState State) (State, error) {
	if oldState.Selected+2 > len(oldState.Folder.Files) {
		return oldState, errors.New("trying to go down below last file")
	}
	newState := copyState(oldState)
	newState.Selected = oldState.Selected + 1
	return newState, nil
}

func (u Up) Execute(oldState State) (State, error) {
	if oldState.Selected == 0 {
		return oldState, errors.New("trying to go above first file")
	}
	newState := copyState(oldState)
	newState.Selected = oldState.Selected - 1
	return newState, nil
}

func (e Enter) Execute(oldState State) (State, error) {
	newFolder := oldState.Folder.Files[oldState.Selected]
	if len(newFolder.Files) == 0 {
		return oldState, errors.New("Trying to enter empty file")
	}
	newHistory := map[*files.File]int{}
	for fp, selected := range oldState.history {
		newHistory[fp] = selected
	}
	newHistory[oldState.Folder] = oldState.Selected
	return State{
		Folder:      newFolder,
		history:     newHistory,
		Selected:    newHistory[newFolder],
		MarkedFiles: oldState.MarkedFiles,
	}, nil
}

func (GoBack) Execute(oldState State) (State, error) {
	parentFolder := oldState.Folder.Parent
	if parentFolder == nil {
		return oldState, errors.New("Trying to go back on root")
	}
	newHistory := map[*files.File]int{}
	for fp, selected := range oldState.history {
		newHistory[fp] = selected
	}
	newHistory[oldState.Folder] = oldState.Selected
	return State{
		Folder:      parentFolder,
		history:     newHistory,
		Selected:    newHistory[parentFolder],
		MarkedFiles: oldState.MarkedFiles,
	}, nil
}

func (m Mark) Execute(oldState State) (State, error) {
	newState := copyState(oldState)
	selectedFile := newState.Folder.Files[newState.Selected]
	if _, exists := newState.MarkedFiles[selectedFile]; exists {
		delete(newState.MarkedFiles, selectedFile)
	} else {
		newState.MarkedFiles[selectedFile] = struct{}{}
	}
	return newState, nil
}
