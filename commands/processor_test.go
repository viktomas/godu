package commands

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/viktomas/godu/core"
)

func TestProcessFolder(t *testing.T) {
	folder := core.NewTestFolder("a", core.NewTestFile("b", 10), core.NewTestFile("c", 50), core.NewTestFile("d", 70))
	ProcessFolder(folder, 30)
	c := &core.File{"c", nil, 50, false, []*core.File{}}
	d := &core.File{"d", nil, 70, false, []*core.File{}}
	a := &core.File{"a", nil, 130, true, []*core.File{d, c}}
	d.Parent = a
	c.Parent = a
	assert.Equal(t, a, folder, "ProcessFoler didn't prune and sort folder")
}

func TestProcessFolderShouldFailWithSmallFiles(t *testing.T) {
	folder := core.NewTestFolder("a", core.NewTestFile("b", 70))
	err := ProcessFolder(folder, 80)
	assert.NotNil(t, err, "ProcessFolder didn't result in error when run on folder with too small files")
}

func TestStartProcessing(t *testing.T) {
	commands := make(chan Executer)
	states := make(chan State, 2)
	lastStateChan := make(chan<- *State, 1)
	folder := core.NewTestFolder("a", core.NewTestFile("b", 10), core.NewTestFile("c", 50))
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
	folder := core.NewTestFolder("a", core.NewTestFile("b", 10), core.NewTestFile("c", 50))
	go StartProcessing(folder, commands, states, lastStateChan, &wg)
	<-states
	commands <- Enter{}
	close(commands)
	wg.Wait() //would block if StartProcessing adds second state to the channel
}
