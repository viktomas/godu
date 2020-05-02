package files

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildFile(t *testing.T) {
	a := &File{"a", nil, 100, false, []*File{}}
	build := NewTestFile("a", 100)
	assert.Equal(t, a, build)
}

func TestBuildFolder(t *testing.T) {
	a := &File{"a", nil, 0, true, []*File{}}
	build := NewTestFolder("a")
	assert.Equal(t, a, build)
}

func TestBuildFolderWithFile(t *testing.T) {
	e := &File{"e", nil, 100, false, []*File{}}
	d := &File{"d", nil, 100, true, []*File{e}}
	e.Parent = d
	build := NewTestFolder("d", NewTestFile("e", 100))
	assert.Equal(t, d, build)
}

func TestBuildComplexFolder(t *testing.T) {
	e := &File{"e", nil, 100, false, []*File{}}
	d := &File{"d", nil, 100, true, []*File{e}}
	e.Parent = d
	b := &File{"b", nil, 50, false, []*File{}}
	c := &File{"c", nil, 100, false, []*File{}}
	a := &File{"a", nil, 250, true, []*File{b, c, d}}
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
