package core

import (
	"reflect"
	"testing"
)

func TestStatePath(t *testing.T) {
	tree := NewTestFolder("a", NewTestFile("b", 50), NewTestFolder("c", NewTestFile("d", 50)))
	subTree := tree.Files[1]
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
	tree := NewTestFolder("a", NewTestFile("b", 50), NewTestFile("c", 50))
	initialState := State{
		Folder: tree,
	}
	newState, _ := Down{}.Execute(initialState)
	if newState.Selected != 1 {
		t.Error("Down command didn't change Selected index")
	}

}

func TestDownCommandFails(t *testing.T) {
	tree := NewTestFolder("a", NewTestFile("b", 50), NewTestFile("c", 50))
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
	tree := NewTestFolder("a", NewTestFile("b", 50), NewTestFile("c", 50))
	initialState := State{
		Folder:   tree,
		Selected: 1,
	}
	newState, _ := Up{}.Execute(initialState)
	if newState.Selected != 0 {
		t.Error("Up command didn't change Selected index")
	}

}

func TestUpCommandFails(t *testing.T) {
	tree := NewTestFolder("a", NewTestFile("b", 50), NewTestFile("c", 50))
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
	tree := NewTestFolder("a", NewTestFile("b", 50), NewTestFolder("c", NewTestFile("d", 50), NewTestFile("e", 50)))
	subTree := tree.Files[1]
	marked := make(map[*File]struct{})
	initialState := State{
		Folder:      tree,
		history:     map[*File]int{subTree: 1},
		Selected:    1,
		MarkedFiles: marked,
	}
	command := Enter{}
	newState, _ := command.Execute(initialState)
	expectedState := State{
		ancestors:   ancestors{tree},
		Folder:      subTree,
		history:     map[*File]int{tree: 1, subTree: 1},
		Selected:    1,
		MarkedFiles: marked,
	}
	if !reflect.DeepEqual(newState, expectedState) {
		t.Errorf("New state is not same as expected %v and %v", newState, expectedState)
	}
}

func TestEnterCommandFails(t *testing.T) {
	tree := NewTestFolder("a", NewTestFile("b", 50))
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
	tree := NewTestFolder("a", NewTestFile("b", 50), NewTestFolder("c", NewTestFile("d", 50), NewTestFile("e", 50)))
	subTree := tree.Files[1]
	marked := make(map[*File]struct{})
	initialState := State{
		ancestors:   ancestors{tree},
		Folder:      subTree,
		history:     map[*File]int{tree: 1},
		Selected:    1,
		MarkedFiles: marked,
	}
	command := GoBack{}
	newState, _ := command.Execute(initialState)
	expectedState := State{
		ancestors:   ancestors{},
		Folder:      tree,
		history:     map[*File]int{tree: 1, subTree: 1},
		Selected:    1,
		MarkedFiles: marked,
	}
	if !reflect.DeepEqual(newState, expectedState) {
		t.Errorf("New state is not same as expected %v and %v", newState, expectedState)
	}
}

func TestGoBackOnRoot(t *testing.T) {
	tree := NewTestFolder("a", NewTestFile("b", 50))
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

func TestMarkFile(t *testing.T) {
	tree := NewTestFolder("a", NewTestFile("b", 50))
	initialState := State{
		ancestors:   ancestors{},
		Folder:      tree,
		Selected:    0,
		MarkedFiles: make(map[*File]struct{}),
	}
	command := Mark{}
	newState, _ := command.Execute(initialState)
	if _, marked := newState.MarkedFiles[tree.Files[0]]; !marked {
		t.Error("File is not marked but should be")
	}
	newState, _ = command.Execute(newState)
	if _, marked := newState.MarkedFiles[tree.Files[0]]; marked {
		t.Error("File is marked but shouldn't be")
	}
}
