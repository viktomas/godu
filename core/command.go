package core

import "fmt"

// State represents system configuration after processing user input
type State struct {
	parent   *File
	folder   *File
	filepath string
}

// Executer represents a user action triggered on a State
type Executer interface {
	Execute(State) (State, error)
}

type Enter struct {
	index int
}

type GoBack struct{}

func (e Enter) Execute(oldState State) (State, error) {
	if e.index < 0 || e.index >= len(oldState.folder.Files) {
		return oldState, fmt.Errorf("Trying to enter nonexistent folder on index %d", e.index)
	}
	newFolder := oldState.folder.Files[0]
	newFilepath := fmt.Sprintf("%s%s/", oldState.filepath, newFolder.Name)
	return State{
		parent:   oldState.folder,
		folder:   newFolder,
		filepath: newFilepath,
	}, nil
}

func (GoBack) Execute(oldState State) (State, error) {
	return oldState, nil
}
