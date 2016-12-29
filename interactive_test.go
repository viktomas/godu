package godu

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func TestPrependsNumberTree(t *testing.T) {
	testTree := &File{"b", 180, []*File{
		&File{"c", 100, []*File{}},
	}}
	expected := "0) c 100B\n"
	testInteractive(testTree, "", expected, t)
}

func TestTakesANumberAndGoesDeeper(t *testing.T) {
	testTree := &File{"b", 180, []*File{
		&File{"c", 100, []*File{
			&File{"d", 10, []*File{}},
		}},
	}}
	input := "0\n"
	expected := "0) c/ 100B\n0) d 10B\n"
	testInteractive(testTree, input, expected, t)
}

func TestOrdersOutputDesc(t *testing.T) {
	testTree := &File{"b", 180, []*File{
		&File{"c", 10, []*File{}},
		&File{"d", 100, []*File{
			&File{"e", 10, []*File{}},
		}},
	}}
	input := "0\n"
	expected := "0) d/ 100B\n1) c 10B\n0) e 10B\n"
	testInteractive(testTree, input, expected, t)
}

func TestGoesBackWhenNotNumber(t *testing.T) {
	testTree := &File{"b", 180, []*File{
		&File{"c", 10, []*File{}},
		&File{"d", 100, []*File{
			&File{"e", 10, []*File{}},
		}},
	}}
	input := "0\nb\n"
	expected := "0) d/ 100B\n1) c 10B\n0) e 10B\n0) d/ 100B\n1) c 10B\n"
	testInteractive(testTree, input, expected, t)

}

func testInteractive(tree *File, input string, expected string, t *testing.T) {
	reader := strings.NewReader(input)
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	InteractiveTree(tree, writer, reader)
	writer.Flush()
	result := buffer.String()
	if result != expected {
		t.Errorf("expected:\n%sbut got:\n%s", expected, result)
	}
}
