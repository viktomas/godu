package core

import (
	"errors"
	"fmt"
)

// State represents system configuration after processing user input
type State struct {
	ancestors ancestors
	Folder    *File
}

// Executer represents a user action triggered on a State
type Executer interface {
	Execute(State) (State, error)
}

type Enter struct {
	Index int
}

type GoBack struct{}

func (s State) Path() string {
	var path string
	for _, ancestor := range s.ancestors {
		path = path + ancestor.Name + "/"
	}
	path = path + s.Folder.Name + "/"
	return path
}

func (e Enter) Execute(oldState State) (State, error) {
	if e.Index < 0 || e.Index >= len(oldState.Folder.Files) {
		return oldState, fmt.Errorf("Trying to enter nonexistent folder on index %d", e.Index)
	}
	newFolder := oldState.Folder.Files[e.Index]
	return State{
		ancestors: oldState.ancestors.push(oldState.Folder),
		Folder:    newFolder,
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
	}, nil
}
