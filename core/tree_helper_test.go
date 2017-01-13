package core

import (
	"reflect"
	"testing"
)

func TstFolder(name string, files ...*File) *File {
	folder := &File{name, nil, 0, true, []*File{}}
	if files == nil {
		return folder
	}
	var size int64
	for _, file := range files {
		size += file.Size
		file.Parent = folder
	}
	folder.Size = size
	folder.Files = files
	return folder
}

func TstFile(name string, size int64) *File {
	return &File{name, nil, size, false, []*File{}}
}

func TestBuildFile(t *testing.T) {
	a := &File{"a", nil, 100, false, []*File{}}
	build := TstFile("a", 100)
	if !reflect.DeepEqual(a, build) {
		t.Error("File dsl didn't get parsed correctly")
	}
}

func TestBuildFolder(t *testing.T) {
	a := &File{"a", nil, 0, true, []*File{}}
	build := TstFolder("a")
	if !reflect.DeepEqual(a, build) {
		t.Error("Folder dsl didn't get parsed correctly")
	}
}

func TestBuildFolderWithFile(t *testing.T) {
	e := &File{"e", nil, 100, false, []*File{}}
	d := &File{"d", nil, 100, true, []*File{e}}
	e.Parent = d
	build := TstFolder("d", TstFile("e", 100))
	if !reflect.DeepEqual(d, build) {
		t.Errorf("Tree %v is not as expected: %v", d, build)
	}
}

func TestBuildTree(t *testing.T) {
	e := &File{"e", nil, 100, false, []*File{}}
	d := &File{"d", nil, 100, true, []*File{e}}
	e.Parent = d
	b := &File{"b", nil, 50, false, []*File{}}
	c := &File{"c", nil, 100, false, []*File{}}
	a := &File{"a", nil, 250, true, []*File{b, c, d}}
	b.Parent = a
	c.Parent = a
	d.Parent = a
	build := TstFolder("a", TstFile("b", 50), TstFile("c", 100), TstFolder("d", TstFile("e", 100)))
	if !reflect.DeepEqual(a, build) {
		t.Error("The dsl hasn't been parsed correctly")
	}
}
