package core

import (
	"reflect"
	"testing"
)

func TestStatePath(t *testing.T) {
	subTree := &File{"c", 50, true, []*File{
		&File{"d", 50, false, []*File{}},
	}}
	tree := &File{"a", 50, true, []*File{
		&File{"b", 50, false, []*File{}},
		subTree,
	}}
	initialState := State{
		ancestors: ancestors{tree},
		Folder:    subTree,
	}
	expected := "a/c/"
	if initialState.Path() != expected {
		t.Errorf("expected path %s, but got %s", expected, initialState.Path())
	}
}

func TestDownCommand(t *testing.T) {
	tree := &File{"a", 100, true, []*File{
		&File{"b", 50, false, []*File{}},
		&File{"c", 50, false, []*File{}},
	}}
	initialState := State{
		Folder: tree,
	}
	newState, _ := Down{}.Execute(initialState)
	if newState.Selected != 1 {
		t.Error("Down command didn't change Selected index")
	}

}

func TestDownCommandFails(t *testing.T) {
	tree := &File{"a", 100, true, []*File{
		&File{"b", 50, false, []*File{}},
		&File{"c", 50, false, []*File{}},
	}}
	initialState := State{
		Folder:   tree,
		Selected: 1,
	}
	newState, err := Down{}.Execute(initialState)
	if err == nil {
		t.Error("Down command didn't fail")
	}
	if !reflect.DeepEqual(newState, initialState) {
		t.Error("State mutated when performing Down")
	}

}

func TestUpCommand(t *testing.T) {
	tree := &File{"a", 100, true, []*File{
		&File{"b", 50, true, []*File{}},
		&File{"c", 50, true, []*File{}},
	}}
	initialState := State{
		Selected: 1,
		Folder:   tree,
	}
	newState, _ := Up{}.Execute(initialState)
	if newState.Selected != 0 {
		t.Error("Up command didn't change Selected index")
	}

}

func TestUpCommandFails(t *testing.T) {
	tree := &File{"a", 100, true, []*File{
		&File{"b", 50, true, []*File{}},
		&File{"c", 50, true, []*File{}},
	}}
	initialState := State{
		Folder:   tree,
		Selected: 0,
	}
	newState, err := Up{}.Execute(initialState)
	if err == nil {
		t.Error("Up command didn't fail")
	}
	if !reflect.DeepEqual(newState, initialState) {
		t.Error("State mutated when performing Up")
	}

}

func TestEnterCommand(t *testing.T) {
	subTree := &File{"c", 50, true, []*File{
		&File{"d", 50, false, []*File{}},
		&File{"f", 50, false, []*File{}},
	}}
	tree := &File{"a", 50, true, []*File{
		&File{"b", 50, false, []*File{}},
		subTree,
	}}
	initialState := State{
		Selected: 1,
		Folder:   tree,
		history:  map[*File]int{subTree: 1},
	}
	command := Enter{}
	newState, _ := command.Execute(initialState)
	expectedState := State{
		ancestors: ancestors{tree},
		Folder:    subTree,
		history:   map[*File]int{tree: 1, subTree: 1},
		Selected:  1,
	}
	if !reflect.DeepEqual(newState, expectedState) {
		t.Errorf("New state is not same as expected %v and %v", newState, expectedState)
	}
}

func TestEnterCommandFails(t *testing.T) {
	tree := &File{"a", 50, true, []*File{
		&File{"b", 50, false, []*File{}},
	}}
	initialState := State{
		Folder: tree,
	}
	command := Enter{}
	_, err := command.Execute(initialState)
	if err == nil {
		t.Error("Command Enter entered a file")
	}

}

func TestGoBackCommand(t *testing.T) {
	subTree := &File{"c", 50, true, []*File{
		&File{"d", 50, false, []*File{}},
		&File{"e", 50, false, []*File{}},
	}}
	tree := &File{"a", 50, true, []*File{
		&File{"b", 50, false, []*File{}},
		subTree,
	}}
	initialState := State{
		ancestors: ancestors{tree},
		Folder:    subTree,
		Selected:  1,
		history:   map[*File]int{tree: 1},
	}
	command := GoBack{}
	newState, _ := command.Execute(initialState)
	expectedState := State{
		ancestors: ancestors{},
		Folder:    tree,
		history:   map[*File]int{tree: 1, subTree: 1},
		Selected:  1,
	}
	if !reflect.DeepEqual(newState, expectedState) {
		t.Errorf("New state is not same as expected %v and %v", newState, expectedState)
	}
}

func TestGoBackOnRoot(t *testing.T) {
	tree := &File{"a", 50, true, []*File{
		&File{"b", 50, false, []*File{}},
	}}
	initialState := State{
		ancestors: ancestors{},
		Folder:    tree,
	}
	command := GoBack{}
	_, err := command.Execute(initialState)
	if err == nil {
		t.Error("GoBack should fail on root")
	}
}
