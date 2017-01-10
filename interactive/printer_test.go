package interactive

import (
	"bufio"
	"bytes"
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
	marked := make(map[*core.File]struct{})
	root := &core.File{Name: "."}
	f1 := &core.File{Name: "f1"}
	f2 := &core.File{Name: "f2"}
	f3 := &core.File{Name: "f3", Parent: nil}
	d1 := &core.File{Name: "d1", IsDir: true, Parent: root, Files: []*core.File{f1}}
	d2 := &core.File{Name: "d2", IsDir: true, Parent: root}
	d3 := &core.File{Name: "d3", IsDir: true, Parent: d1, Files: []*core.File{f2}}
	f1.Parent = d1
	f2.Parent = d3
	marked[d1] = struct{}{}
	marked[d2] = struct{}{}
	marked[f1] = struct{}{}
	marked[f2] = struct{}{}
	marked[f3] = struct{}{}
	state := &core.State{MarkedFiles: marked}
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	PrintMarkedFiles(state, writer)
	writer.Flush()
	result := buffer.String()
	// We don't know the order, as we are using a map to store marked files :/
	expected := `'d1'
'd2'
'd1/d3/f2'`
	if hasSameLines(result, expected) {
		t.Errorf("Expected '%s' from PrintMarkedFiles, got '%s'", expected, result)
	}
}

func hasSameLines(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}
	ss1 := strings.Split(s1, "\n")
	for _, l1 := range ss1 {
		if !strings.Contains(s2, l1+"\n") {
			return false
		}
	}
	return true
}
