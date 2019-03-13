package interactive

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/viktomas/godu/core"
)

func TestFilesAsSliceEmptyMap(t *testing.T) {
	marked := make(map[*core.File]struct{})
	result := FilesAsSlice(marked)
	assert.Equal(t, 0, len(result))
}

func TestFilesAsSlice(t *testing.T) {
	root := core.NewTestFolder(".",
		core.NewTestFile("'single''quotes'", 0),
		core.NewTestFolder("d1",
			core.NewTestFile("f1", 0),
			core.NewTestFolder("d3",
				core.NewTestFile("f2", 0),
			),
		),
		core.NewTestFolder("d2"),
		core.NewTestFile("f3", 0),
	)
	marked := make(map[*core.File]struct{})
	marked[core.FindTestFile(root, "d1")] = struct{}{}
	marked[core.FindTestFile(root, "d2")] = struct{}{}
	marked[core.FindTestFile(root, "f1")] = struct{}{}
	marked[core.FindTestFile(root, "f2")] = struct{}{}
	marked[core.FindTestFile(root, "'single''quotes'")] = struct{}{}

	want := []string{"'single''quotes'", "d1/d3/f2", "d1/f1", "d1", "d2"}
	got := FilesAsSlice(marked)
	assert.Equal(t, want, got)
}
