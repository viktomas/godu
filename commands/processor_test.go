package commands

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/viktomas/godu/files"
)

func TestProcessFolder(t *testing.T) {
	folder := files.NewTestFolder("a", files.NewTestFile("b", 10), files.NewTestFile("c", 50), files.NewTestFile("d", 70))
	ProcessFolder(folder, 30)
	c := &files.File{"c", nil, 50, false, []*files.File{}}
	d := &files.File{"d", nil, 70, false, []*files.File{}}
	a := &files.File{"a", nil, 130, true, []*files.File{d, c}}
	d.Parent = a
	c.Parent = a
	assert.Equal(t, a, folder, "ProcessFoler didn't prune and sort folder")
}

func TestProcessFolderShouldFailWithSmallFiles(t *testing.T) {
	folder := files.NewTestFolder("a", files.NewTestFile("b", 70))
	err := ProcessFolder(folder, 80)
	assert.NotNil(t, err, "ProcessFolder didn't result in error when run on folder with too small files")
}

func TestStartProcessing(t *testing.T) {
	commands := make(chan Executer)
	states := make(chan State, 2)
	lastStateChan := make(chan<- *State, 1)
	folder := files.NewTestFolder("a", files.NewTestFile("b", 10), files.NewTestFile("c", 50))
	var wg sync.WaitGroup
	wg.Add(1)
	go StartProcessing(folder, commands, states, lastStateChan, &wg)
	commands <- Down{}
	close(commands)
	state := <-states
	state = <-states
	wg.Wait()
	assert.Equal(t, 1, state.Selected, "StartProcessing didn't process command")
	state, ok := <-states
	assert.False(t, ok, "forgot to close states channel")
}

func TestDoesntProcessInvalidCommand(t *testing.T) {
	commands := make(chan Executer)
	states := make(chan State)
	lastStateChan := make(chan<- *State, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	folder := files.NewTestFolder("a", files.NewTestFile("b", 10), files.NewTestFile("c", 50))
	go StartProcessing(folder, commands, states, lastStateChan, &wg)
	<-states
	commands <- Enter{}
	close(commands)
	wg.Wait() //would block if StartProcessing adds second state to the channel
}
