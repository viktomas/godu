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

func TestEnterCommand(t *testing.T) {
	subTree := &File{"c", 50, true, []*File{
		&File{"d", 50, false, []*File{}},
	}}
	tree := &File{"a", 50, true, []*File{
		&File{"b", 50, false, []*File{}},
		subTree,
	}}
	initialState := State{
		Folder: tree,
	}
	command := Enter{1}
	newState, _ := command.Execute(initialState)
	expectedState := State{
		ancestors: ancestors{tree},
		Folder:    subTree,
	}
	if !reflect.DeepEqual(newState, expectedState) {
		t.Errorf("New state is not same as expected %v and %v", newState, expectedState)
	}
}

func TestEnterCommandFails(t *testing.T) {
	tree := &File{"a", 50, true, []*File{
		&File{"b", 50, true, []*File{
			&File{"c", 50, false, []*File{}},
		}},
	}}
	initialState := State{
		Folder: tree,
	}
	command := Enter{1}
	_, err := command.Execute(initialState)
	if err == nil {
		t.Error("Command Enter didn't fail when trying to enter nonexistent folder")
	}

}

func TestGoBackCommand(t *testing.T) {
	subTree := &File{"b", 50, true, []*File{
		&File{"c", 50, false, []*File{}},
	}}
	tree := &File{"a", 50, true, []*File{
		subTree,
	}}
	initialState := State{
		ancestors: ancestors{tree},
		Folder:    subTree,
	}
	command := GoBack{}
	newState, _ := command.Execute(initialState)
	expectedState := State{
		ancestors: ancestors{},
		Folder:    tree,
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
