package interactive

import (
	"bufio"
	"bytes"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/viktomas/godu/core"
)

func TestPrintMarkedFilesNone(t *testing.T) {
	marked := make(map[*core.File]struct{})
	state := &core.State{MarkedFiles: marked}
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	defer writer.Flush()
	PrintMarkedFiles(state, writer)
	result := buffer.String()
	if result != "" {
		t.Errorf("Expected empty output from PrintMarkedFiles, got '%s'", result)
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
	//marked[getFileByName(root, "f3")] = struct{}{}
	state := &core.State{MarkedFiles: marked}
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	PrintMarkedFiles(state, writer)
	writer.Flush()
	result := buffer.String()
	// We don't know the order, as we are using a map to store marked files :/
	expected := "'d1'\n'd1/f1'\n'd2'\n'd1/d3/f2'\n"
	if !hasSameLines(result, expected) {
		t.Errorf("Expected '%s' from PrintMarkedFiles, got '%s'", expected, result)
	}
}

func hasSameLines(value, expected string) bool {
	values := strings.Split(value, "\n")
	expecteds := strings.Split(expected, "\n")
	sort.Sort(sort.StringSlice(values))
	sort.Sort(sort.StringSlice(expecteds))
	return reflect.DeepEqual(values, expecteds)
}
