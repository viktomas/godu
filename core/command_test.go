package core

import (
	"reflect"
	"testing"
)

func TestEnterCommand(t *testing.T) {
	subTree := &File{"b", 50, []*File{
		&File{"c", 50, []*File{}},
	}}
	tree := &File{"a", 50, []*File{
		subTree,
	}}
	initialState := State{
		folder:   tree,
		filepath: "a/",
	}
	command := Enter{0}
	newState, _ := command.Execute(initialState)
	expectedState := State{
		ancestors: ancestors{tree},
		folder:    subTree,
		filepath:  "a/b/",
	}
	if !reflect.DeepEqual(newState, expectedState) {
		t.Errorf("New state is not same as expected %v and %v", newState, expectedState)
	}
}

func TestEnterCommandFails(t *testing.T) {
	tree := &File{"a", 50, []*File{
		&File{"b", 50, []*File{
			&File{"c", 50, []*File{}},
		}},
	}}
	initialState := State{
		folder: tree,
	}
	command := Enter{1}
	_, err := command.Execute(initialState)
	if err == nil {
		t.Error("Command Enter didn't fail when trying to enter nonexistent folder")
	}

}

func TestGoBackCommand(t *testing.T) {
	subTree := &File{"b", 50, []*File{
		&File{"c", 50, []*File{}},
	}}
	tree := &File{"a", 50, []*File{
		subTree,
	}}
	initialState := State{
		ancestors: ancestors{tree},
		folder:    subTree,
		filepath:  "a/b/",
	}
	command := GoBack{}
	newState, _ := command.Execute(initialState)
	expectedState := State{
		folder:   tree,
		filepath: "a/",
	}
	if !reflect.DeepEqual(newState, expectedState) {
		t.Errorf("New state is not same as expected %v and %v", newState, expectedState)
	}
}
