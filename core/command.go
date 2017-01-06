package core

import "errors"

// State represents system configuration after processing user input
type State struct {
	ancestors ancestors
	Folder    *File
	Selected  int
}

// Executer represents a user action triggered on a State
type Executer interface {
	Execute(State) (State, error)
}

type Enter struct{}

type GoBack struct{}

type Down struct{}

type Up struct{}

func (s State) Path() string {
	var path string
	for _, ancestor := range s.ancestors {
		path = path + ancestor.Name + "/"
	}
	path = path + s.Folder.Name + "/"
	return path
}

func (d Down) Execute(oldState State) (State, error) {
	if oldState.Selected+2 > len(oldState.Folder.Files) {
		return oldState, errors.New("trying to go down below last file")
	}
	return State{
		ancestors: oldState.ancestors,
		Folder:    oldState.Folder,
		Selected:  oldState.Selected + 1,
	}, nil
}

func (u Up) Execute(oldState State) (State, error) {
	if oldState.Selected == 0 {
		return oldState, errors.New("trying to go above first file")
	}
	return State{
		ancestors: oldState.ancestors,
		Folder:    oldState.Folder,
		Selected:  oldState.Selected - 1,
	}, nil
}

func (e Enter) Execute(oldState State) (State, error) {
	newFolder := oldState.Folder.Files[oldState.Selected]
	if len(newFolder.Files) == 0 {
		return oldState, errors.New("Trying to enter empty file")
	}
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
