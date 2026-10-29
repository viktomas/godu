package interactive

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/viktomas/godu/files"
)

func TestFilesAsSliceEmptyMap(t *testing.T) {
	marked := make(map[*files.File]struct{})
	result := FilesAsSlice(marked)
	assert.Equal(t, 0, len(result))
}

func TestFilesAsSlice(t *testing.T) {
	root := files.NewTestFolder(".",
		files.NewTestFile("'single''quotes'", 0),
		files.NewTestFolder("d1",
			files.NewTestFile("f1", 0),
			files.NewTestFolder("d3",
				files.NewTestFile("f2", 0),
			),
		),
		files.NewTestFolder("d2"),
		files.NewTestFile("f3", 0),
	)
	marked := make(map[*files.File]struct{})
	marked[files.FindTestFile(root, "d1")] = struct{}{}
	marked[files.FindTestFile(root, "d2")] = struct{}{}
	marked[files.FindTestFile(root, "f1")] = struct{}{}
	marked[files.FindTestFile(root, "f2")] = struct{}{}
	marked[files.FindTestFile(root, "'single''quotes'")] = struct{}{}

	want := []string{"'single''quotes'", filepath.Join("d1", "d3", "f2"), filepath.Join("d1", "f1"), "d1", "d2"}
	got := FilesAsSlice(marked)
	assert.Equal(t, want, got)
}

func TestFilesAsSliceSortsEqualLengthPathsAlphabetically(t *testing.T) {
	// 13 paths exceed insertion sort's threshold (12), so an unstable
	// by-length sort reorders equal-length paths; the alphabetical
	// tie-break keeps them ordered (d1/x ties with bb00..bb02 at 4 chars)
	root := files.NewTestFolder(".",
		files.NewTestFolder("d1", files.NewTestFile("x", 0)),
		files.NewTestFile("aa0", 0),
		files.NewTestFile("aa1", 0),
		files.NewTestFile("aa2", 0),
		files.NewTestFile("bb00", 0),
		files.NewTestFile("bb01", 0),
		files.NewTestFile("bb02", 0),
		files.NewTestFile("c0", 0),
		files.NewTestFile("c1", 0),
		files.NewTestFile("c2", 0),
		files.NewTestFile("c3", 0),
		files.NewTestFile("c4", 0),
		files.NewTestFile("c5", 0),
	)
	marked := make(map[*files.File]struct{})
	marked[files.FindTestFile(root, "x")] = struct{}{}
	for _, file := range root.Files {
		if !file.IsDir {
			marked[file] = struct{}{}
		}
	}
	// longer paths come first, equal-length paths alphabetically
	want := []string{"bb00", "bb01", "bb02", filepath.Join("d1", "x"), "aa0", "aa1", "aa2", "c0", "c1", "c2", "c3", "c4", "c5"}
	got := FilesAsSlice(marked)
	assert.Equal(t, want, got)
}
