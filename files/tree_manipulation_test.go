package files

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortFolder(t *testing.T) {
	folder := NewTestFolder("b",
		NewTestFile("c", 100),
		NewTestFolder("d",
			NewTestFile("e", 50),
			NewTestFile("f", 30),
			NewTestFolder("g",
				NewTestFile("i", 30),
				NewTestFile("j", 50),
			),
		),
	)
	expected := NewTestFolder("b",
		NewTestFolder("d",
			NewTestFolder("g",
				NewTestFile("j", 50),
				NewTestFile("i", 30),
			),
			NewTestFile("e", 50),
			NewTestFile("f", 30),
		),
		NewTestFile("c", 100),
	)

	SortDesc(folder)
	assert.Equal(t, expected, folder)
}

func TestPruneFolder(t *testing.T) {
	folder := &File{Name: "b", Size: 260, IsDir: true, Files: []*File{
		{Name: "c", Size: 100, Files: []*File{}},
		{Name: "d", Size: 160, IsDir: true, Files: []*File{
			{Name: "e", Size: 50, Files: []*File{}},
			{Name: "f", Size: 30, Files: []*File{}},
			{Name: "g", Size: 80, IsDir: true, Files: []*File{
				{Name: "i", Size: 50, Files: []*File{}},
				{Name: "j", Size: 30, Files: []*File{}},
			}},
		}},
	}}
	expected := &File{Name: "b", Size: 260, IsDir: true, Files: []*File{
		{Name: "c", Size: 100, Files: []*File{}},
		{Name: "d", Size: 160, IsDir: true, Files: []*File{
			{Name: "g", Size: 80, IsDir: true, Files: []*File{}},
		}},
	}}
	PruneSmallFiles(folder, 60)
	assert.Equal(t, expected, folder)
}
