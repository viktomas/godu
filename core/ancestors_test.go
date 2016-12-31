package core

import (
	"reflect"
	"testing"
)

func TestPush(t *testing.T) {
	file := &File{}
	a := ancestors{file}
	newA := a.push(file)
	if len(newA) != 2 {
		t.Errorf("ancestors size is not 2 but %d", len(newA))
	}
	if !reflect.DeepEqual(a, ancestors{file}) {
		t.Error("ancestors got mutated :(")
	}
}

func TestPop(t *testing.T) {
	fileA := &File{Name: "a"}
	fileB := &File{Name: "b"}
	a := ancestors{fileA, fileB}
	file, newA := a.pop()
	if !reflect.DeepEqual(file, fileB) {
		t.Error("We popped wrong file")
	}
	if !reflect.DeepEqual(newA, ancestors{fileA}) {
		t.Error("We don't have correct new stack")
	}
	if !reflect.DeepEqual(a, ancestors{fileA, fileB}) {
		t.Error("Ancestors got mutated")
	}
}

func TestPopEmpty(t *testing.T) {
	a := ancestors{}
	file, _ := a.pop()
	if file != nil {
		t.Error("pop managed to pop something from empty stack")
	}
}
