package interactive

import (
	"path"
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

	want := []string{"'single''quotes'", path.Join("d1", "d3", "f2"), path.Join("d1", "f1"), "d1", "d2"}
	got := FilesAsSlice(marked)
	assert.Equal(t, want, got)
}
