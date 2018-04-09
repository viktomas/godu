package interactive

import (
	"reflect"
	"testing"

	"github.com/viktomas/godu/core"
)

func TestFilesAsSliceEmptyMap(t *testing.T) {
	marked := make(map[*core.File]struct{})
	result := FilesAsSlice(marked, false)
	if len(result) > 0 {
		t.Errorf("Expected empty output, got '%v'", result)
	}
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

	evalFunc := func(want, got []string, t *testing.T) {
		// t.Helper() // TODO: Requires Go 1.9
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Expected '%v', got '%v'", want, got)
		}
	}

	t.Run("unquoted", func(t *testing.T) {
		want := []string{"'single''quotes'", "d1/d3/f2", "d1/f1", "d1", "d2"}
		got := FilesAsSlice(marked, false)
		evalFunc(want, got, t)
	})

	t.Run("quoted", func(t *testing.T) {
		want := []string{"'\\'single\\'\\'quotes\\''", "'d1/d3/f2'", "'d1/f1'", "'d1'", "'d2'"}
		got := FilesAsSlice(marked, true)
		evalFunc(want, got, t)
	})
}
