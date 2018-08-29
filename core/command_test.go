package core

import (
	"reflect"
	"testing"
)

func TestDownCommand(t *testing.T) {
	folder := NewTestFolder("a", NewTestFile("b", 50), NewTestFile("c", 50))
	initialState := State{
		Folder: folder,
	}
	newState, _ := Down{}.Execute(initialState)
	if newState.Selected != 1 {
		t.Error("Down command didn't change Selected index")
	}

}

func TestDownCommandFails(t *testing.T) {
	folder := NewTestFolder("a", NewTestFile("b", 50), NewTestFile("c", 50))
	initialState := State{
		Folder:   folder,
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
	folder := NewTestFolder("a", NewTestFile("b", 50), NewTestFile("c", 50))
	initialState := State{
		Folder:   folder,
		Selected: 1,
	}
	newState, _ := Up{}.Execute(initialState)
	if newState.Selected != 0 {
		t.Error("Up command didn't change Selected index")
	}

}

func TestUpCommandFails(t *testing.T) {
	folder := NewTestFolder("a", NewTestFile("b", 50), NewTestFile("c", 50))
	initialState := State{
		Folder:   folder,
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
	folder := NewTestFolder("a", NewTestFile("b", 50), NewTestFolder("c", NewTestFile("d", 50), NewTestFile("e", 50)))
	subFolder := folder.Files[1]
	marked := make(map[*File]struct{})
	initialState := State{
		Folder:      folder,
		history:     map[*File]int{subFolder: 1},
		Selected:    1,
		MarkedFiles: marked,
	}
	command := Enter{}
	newState, _ := command.Execute(initialState)
	expectedState := State{
		Folder:      subFolder,
		history:     map[*File]int{folder: 1, subFolder: 1},
		Selected:    1,
		MarkedFiles: marked,
	}
	if !reflect.DeepEqual(newState, expectedState) {
		t.Errorf("New state is not same as expected %v and %v", newState, expectedState)
	}
}

func TestEnterCommandFails(t *testing.T) {
	folder := NewTestFolder("a", NewTestFile("b", 50))
	initialState := State{
		Folder: folder,
	}
	command := Enter{}
	_, err := command.Execute(initialState)
	if err == nil {
		t.Error("Command Enter entered a file")
	}

}

func TestGoBackCommand(t *testing.T) {
	folder := NewTestFolder("a", NewTestFile("b", 50), NewTestFolder("c", NewTestFile("d", 50), NewTestFile("e", 50)))
	subFolder := folder.Files[1]
	marked := make(map[*File]struct{})
	initialState := State{
		Folder:      subFolder,
		history:     map[*File]int{folder: 1},
		Selected:    1,
		MarkedFiles: marked,
	}
	command := GoBack{}
	newState, _ := command.Execute(initialState)
	expectedState := State{
		Folder:      folder,
		history:     map[*File]int{folder: 1, subFolder: 1},
		Selected:    1,
		MarkedFiles: marked,
	}
	if !reflect.DeepEqual(newState, expectedState) {
		t.Errorf("New state is not same as expected %v and %v", newState, expectedState)
	}
}

func TestGoBackOnRoot(t *testing.T) {
	folder := NewTestFolder("a", NewTestFile("b", 50))
	initialState := State{
		Folder: folder,
	}
	command := GoBack{}
	_, err := command.Execute(initialState)
	if err == nil {
		t.Error("GoBack should fail on root")
	}
}

func TestMarkFile(t *testing.T) {
	folder := NewTestFolder("a", NewTestFile("b", 50))
	initialState := State{
		Folder:      folder,
		Selected:    0,
		MarkedFiles: make(map[*File]struct{}),
	}
	command := Mark{}
	newState, _ := command.Execute(initialState)
	if _, marked := newState.MarkedFiles[folder.Files[0]]; !marked {
		t.Error("File is not marked but should be")
	}
	newState, _ = command.Execute(newState)
	if _, marked := newState.MarkedFiles[folder.Files[0]]; marked {
		t.Error("File is marked but shouldn't be")
	}
}
