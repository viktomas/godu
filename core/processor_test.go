package core

import (
	"reflect"
	"sync"
	"testing"
)

func TestStartProcessingPrunesAndSorts(t *testing.T) {
	commands := make(chan Executer)
	states := make(chan State)
	tree := &File{"a", 130, true, []*File{
		&File{"b", 10, false, []*File{}},
		&File{"c", 50, false, []*File{}},
		&File{"d", 70, false, []*File{}},
	}}
	var wg sync.WaitGroup
	wg.Add(1)
	go StartProcessing(tree, 30, commands, states, &wg)
	close(commands)
	state := <-states
	expected := &File{"a", 130, true, []*File{
		&File{"d", 70, false, []*File{}},
		&File{"c", 50, false, []*File{}},
	}}
	wg.Wait()
	if !reflect.DeepEqual(*state.Folder, *expected) {
		t.Error("StartProcessing didn't prune and sort tree")
	}
}

func TestStartProcessing(t *testing.T) {
	commands := make(chan Executer)
	states := make(chan State, 2)
	tree := &File{"a", 60, true, []*File{
		&File{"b", 10, false, []*File{}},
		&File{"c", 50, false, []*File{}},
	}}
	var wg sync.WaitGroup
	wg.Add(1)
	go StartProcessing(tree, 0, commands, states, &wg)
	commands <- Down{}
	close(commands)
	state := <-states
	state = <-states
	wg.Wait()
	if state.Selected != 1 {
		t.Error("StartProcessing didn't process command")
	}
}

func TestDoesntProcessInvalidCommand(t *testing.T) {
	commands := make(chan Executer)
	states := make(chan State)
	var wg sync.WaitGroup
	wg.Add(1)
	tree := &File{"a", 60, true, []*File{
		&File{"b", 10, false, []*File{}},
		&File{"c", 50, false, []*File{}},
	}}
	go StartProcessing(tree, 0, commands, states, &wg)
	<-states
	commands <- Enter{}
	close(commands)
	wg.Wait() //would block if StartProcessing adds second state to the channel
}
