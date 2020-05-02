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
	folder := &File{"b", nil, 260, true, []*File{
		&File{"c", nil, 100, false, []*File{}},
		&File{"d", nil, 160, true, []*File{
			&File{"e", nil, 50, false, []*File{}},
			&File{"f", nil, 30, false, []*File{}},
			&File{"g", nil, 80, true, []*File{
				&File{"i", nil, 50, false, []*File{}},
				&File{"j", nil, 30, false, []*File{}},
			}},
		}},
	}}
	expected := &File{"b", nil, 260, true, []*File{
		&File{"c", nil, 100, false, []*File{}},
		&File{"d", nil, 160, true, []*File{
			&File{"g", nil, 80, true, []*File{}},
		}},
	}}
	PruneSmallFiles(folder, 60)
	assert.Equal(t, expected, folder)
}
