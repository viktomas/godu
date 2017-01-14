package interactive

import (
	"reflect"
	"testing"

	"github.com/viktomas/godu/core"
)

func TestPrintEmptyMarkedFiles(t *testing.T) {
	marked := make(map[*core.File]struct{})
	result := QuoteMarkedFiles(marked)
	if len(result) > 0 {
		t.Errorf("Expected empty output from PrintMarkedFiles, got '%v'", result)
	}
}

func TestPrintMarkedFiles(t *testing.T) {
	root := core.NewTestFolder(".",
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
	result := QuoteMarkedFiles(marked)
	expected := []string{"'d1/d3/f2'", "'d1/f1'", "'d1'", "'d2'"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected '%v' from PrintMarkedFiles, got '%v'", expected, result)
	}
}
