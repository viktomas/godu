package core

import "errors"

// State represents system configuration after processing user input
type State struct {
	ancestors ancestors
	Folder    *File
	Selected  int
	history   map[*File]int // history of all selected postions
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

func copyState(state State) State {
	return State{
		ancestors: state.ancestors,
		Folder:    state.Folder,
		Selected:  state.Selected,
		history:   state.history,
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
	newHistory := map[*File]int{}
	for fp, selected := range oldState.history {
		newHistory[fp] = selected
	}
	newHistory[oldState.Folder] = oldState.Selected
	return State{
		ancestors: oldState.ancestors.push(oldState.Folder),
		Folder:    newFolder,
		history:   newHistory,
		Selected:  newHistory[newFolder],
	}, nil
}

func (GoBack) Execute(oldState State) (State, error) {
	parentFolder, newAncestors := oldState.ancestors.pop()
	if parentFolder == nil {
		return oldState, errors.New("Trying to go back on root")
	}
	newHistory := map[*File]int{}
	for fp, selected := range oldState.history {
		newHistory[fp] = selected
	}
	newHistory[oldState.Folder] = oldState.Selected
	return State{
		ancestors: newAncestors,
		Folder:    parentFolder,
		Selected:  newHistory[parentFolder],
		history:   newHistory,
	}, nil
}
