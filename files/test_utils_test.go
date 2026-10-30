package files

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildFile(t *testing.T) {
	a := &File{Name: "a", Size: 100, Files: []*File{}}
	build := NewTestFile("a", 100)
	assert.Equal(t, a, build)
}

func TestBuildFolder(t *testing.T) {
	a := &File{Name: "a", IsDir: true, Files: []*File{}}
	build := NewTestFolder("a")
	assert.Equal(t, a, build)
}

func TestBuildFolderWithFile(t *testing.T) {
	e := &File{Name: "e", Size: 100, Files: []*File{}}
	d := &File{Name: "d", Size: 100, IsDir: true, Files: []*File{e}}
	e.Parent = d
	build := NewTestFolder("d", NewTestFile("e", 100))
	assert.Equal(t, d, build)
}

func TestBuildComplexFolder(t *testing.T) {
	e := &File{Name: "e", Size: 100, Files: []*File{}}
	d := &File{Name: "d", Size: 100, IsDir: true, Files: []*File{e}}
	e.Parent = d
	b := &File{Name: "b", Size: 50, Files: []*File{}}
	c := &File{Name: "c", Size: 100, Files: []*File{}}
	a := &File{Name: "a", Size: 250, IsDir: true, Files: []*File{b, c, d}}
	b.Parent = a
	c.Parent = a
	d.Parent = a
	build := NewTestFolder("a", NewTestFile("b", 50), NewTestFile("c", 100), NewTestFolder("d", NewTestFile("e", 100)))
	assert.Equal(t, a, build)
}

func TestFindTestFile(t *testing.T) {
	folder := NewTestFolder("a",
		NewTestFolder("b",
			NewTestFile("c", 10),
			NewTestFile("d", 100),
		),
	)
	expected := folder.Files[0].Files[1]
	foundFile := FindTestFile(folder, "d")
	assert.Equal(t, expected, foundFile)
}
