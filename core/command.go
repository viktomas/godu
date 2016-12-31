package core

import (
	"errors"
	"fmt"
)

// State represents system configuration after processing user input
type State struct {
	ancestors ancestors
	folder    *File
	filepath  string
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
		ancestors: oldState.ancestors.push(oldState.folder),
		folder:    newFolder,
		filepath:  newFilepath,
	}, nil
}

func (GoBack) Execute(oldState State) (State, error) {
	parentFolder, newAncestors := oldState.ancestors.pop()
	if parentFolder == nil {
		return oldState, errors.New("Trying to go back on root")
	}
	return State{
		ancestors: newAncestors,
		folder:    parentFolder,
		filepath:  "a/",
	}, nil
}
