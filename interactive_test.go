package godu

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestPrependsNumberTree(t *testing.T) {
	testTree := File{"b", 180, []File{
		File{"c", 100, []File{}},
	}}
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	InteractiveTree(testTree, writer, os.Stdin)
	writer.Flush()
	result := buffer.String()
	expected := "0) c 100B\n"
	if result != expected {
		t.Errorf("expected:\n%sbut got:\n%s", expected, result)
	}
}

func TestTakesANumberAndGoesDeeper(t *testing.T) {
	testTree := File{"b", 180, []File{
		File{"c", 100, []File{
			File{"d", 10, []File{}},
		}},
	}}
	reader := strings.NewReader("0\n")
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	InteractiveTree(testTree, writer, reader)
	writer.Flush()
	result := buffer.String()
	expected := "0) c/ 100B\n0) d 10B\n"
	if result != expected {
		t.Errorf("expected:\n%sbut got:\n%s", expected, result)
	}
}
