package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownCommand(t *testing.T) {
	folder := NewTestFolder("a", NewTestFile("b", 50), NewTestFile("c", 50))
	initialState := State{
		Folder: folder,
	}
	newState, _ := Down{}.Execute(initialState)
	assert.Equal(t, 1, newState.Selected, "Down command didn't change Selected index")

}

func TestDownCommandFails(t *testing.T) {
	folder := NewTestFolder("a", NewTestFile("b", 50), NewTestFile("c", 50))
	initialState := State{
		Folder:   folder,
		Selected: 1,
	}
	newState, err := Down{}.Execute(initialState)
	assert.NotNil(t, err, "Down command didn't fail")
	assert.Equal(t, initialState, newState, "State mutated when performing Down")
}

func TestUpCommand(t *testing.T) {
	folder := NewTestFolder("a", NewTestFile("b", 50), NewTestFile("c", 50))
	initialState := State{
		Folder:   folder,
		Selected: 1,
	}
	newState, _ := Up{}.Execute(initialState)
	assert.Equal(t, 0, newState.Selected, "Up command didn't change Selected index")
}

func TestUpCommandFails(t *testing.T) {
	folder := NewTestFolder("a", NewTestFile("b", 50), NewTestFile("c", 50))
	initialState := State{
		Folder:   folder,
		Selected: 0,
	}
	newState, err := Up{}.Execute(initialState)
	assert.NotNil(t, err, "Up command didn't fail")
	assert.Equal(t, initialState, newState, "State mutated when performing Up")
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
	assert.Equal(t, expectedState, newState)
}

func TestEnterCommandFails(t *testing.T) {
	folder := NewTestFolder("a", NewTestFile("b", 50))
	initialState := State{
		Folder: folder,
	}
	command := Enter{}
	_, err := command.Execute(initialState)
	assert.NotNil(t, err, "Command Enter entered a file")
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
	assert.Equal(t, expectedState, newState)
}

func TestGoBackOnRoot(t *testing.T) {
	folder := NewTestFolder("a", NewTestFile("b", 50))
	initialState := State{
		Folder: folder,
	}
	command := GoBack{}
	_, err := command.Execute(initialState)
	assert.NotNil(t, err, "GoBack should fail on root")
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
	_, marked := newState.MarkedFiles[folder.Files[0]]
	assert.True(t, marked)
	newState, _ = command.Execute(newState)
	_, marked = newState.MarkedFiles[folder.Files[0]]
	assert.False(t, marked)
}
