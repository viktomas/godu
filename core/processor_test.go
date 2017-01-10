package core

import (
	"reflect"
	"sync"
	"testing"
)

func TestPrepareTree(t *testing.T) {
	tree := &File{"a", nil, 130, true, []*File{
		&File{"b", nil, 10, false, []*File{}},
		&File{"c", nil, 50, false, []*File{}},
		&File{"d", nil, 70, false, []*File{}},
	}}
	PrepareTree(tree, 30)
	expected := &File{"a", nil, 130, true, []*File{
		&File{"d", nil, 70, false, []*File{}},
		&File{"c", nil, 50, false, []*File{}},
	}}
	if !reflect.DeepEqual(*tree, *expected) {
		t.Error("PrepareTree didn't prune and sort tree")
	}
}

func TestPrepareTreeShouldFailWithSmallFiles(t *testing.T) {
	tree := &File{"a", nil, 70, true, []*File{
		&File{"b", nil, 70, false, []*File{}},
	}}
	err := PrepareTree(tree, 80)
	if err == nil {
		t.Error("PrepareTree didn't result in error when run on folder with too small files")
	}
}

func TestStartProcessing(t *testing.T) {
	commands := make(chan Executer)
	states := make(chan State, 2)
	lastStateChan := make(chan<- *State, 1)
	tree := &File{"a", nil, 60, true, []*File{
		&File{"b", nil, 10, false, []*File{}},
		&File{"c", nil, 50, false, []*File{}},
	}}
	var wg sync.WaitGroup
	wg.Add(1)
	go StartProcessing(tree, commands, states, lastStateChan, &wg)
	commands <- Down{}
	close(commands)
	state := <-states
	state = <-states
	wg.Wait()
	if state.Selected != 1 {
		t.Error("StartProcessing didn't process command")
	}
	state, ok := <-states
	if ok {
		t.Error("forgot to close states channel")
	}
}

func TestDoesntProcessInvalidCommand(t *testing.T) {
	commands := make(chan Executer)
	states := make(chan State)
	lastStateChan := make(chan<- *State, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	tree := &File{"a", nil, 60, true, []*File{
		&File{"b", nil, 10, false, []*File{}},
		&File{"c", nil, 50, false, []*File{}},
	}}
	go StartProcessing(tree, commands, states, lastStateChan, &wg)
	<-states
	commands <- Enter{}
	close(commands)
	wg.Wait() //would block if StartProcessing adds second state to the channel
}
