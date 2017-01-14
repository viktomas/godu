package core

import (
	"reflect"
	"sync"
	"testing"
)

func TestPrepareTree(t *testing.T) {
	tree := TstFolder("a", TstFile("b", 10), TstFile("c", 50), TstFile("d", 70))
	PrepareTree(tree, 30)
	c := &File{"c", nil, 50, false, []*File{}}
	d := &File{"d", nil, 70, false, []*File{}}
	a := &File{"a", nil, 130, true, []*File{d, c}}
	d.Parent = a
	c.Parent = a
	if !reflect.DeepEqual(*tree, *a) {
		t.Error("PrepareTree didn't prune and sort tree")
	}
}

func TestPrepareTreeShouldFailWithSmallFiles(t *testing.T) {
	tree := TstFolder("a", TstFile("b", 70))
	err := PrepareTree(tree, 80)
	if err == nil {
		t.Error("PrepareTree didn't result in error when run on folder with too small files")
	}
}

func TestStartProcessing(t *testing.T) {
	commands := make(chan Executer)
	states := make(chan State, 2)
	lastStateChan := make(chan<- *State, 1)
	tree := TstFolder("a", TstFile("b", 10), TstFile("c", 50))
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
	tree := TstFolder("a", TstFile("b", 10), TstFile("c", 50))
	go StartProcessing(tree, commands, states, lastStateChan, &wg)
	<-states
	commands <- Enter{}
	close(commands)
	wg.Wait() //would block if StartProcessing adds second state to the channel
}
