package core

import (
	"errors"
	"fmt"
)

// State represents system configuration after processing user input
type State struct {
	ancestors ancestors
	Folder    *File
	filepath  string
}

// Executer represents a user action triggered on a State
type Executer interface {
	Execute(State) (State, error)
}

type Enter struct {
	Index int
}

type GoBack struct{}

func (e Enter) Execute(oldState State) (State, error) {
	if e.Index < 0 || e.Index >= len(oldState.Folder.Files) {
		return oldState, fmt.Errorf("Trying to enter nonexistent folder on index %d", e.Index)
	}
	newFolder := oldState.Folder.Files[1]
	newFilepath := fmt.Sprintf("%s%s/", oldState.filepath, newFolder.Name)
	return State{
		ancestors: oldState.ancestors.push(oldState.Folder),
		Folder:    newFolder,
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
		Folder:    parentFolder,
		filepath:  "a/",
	}, nil
}
