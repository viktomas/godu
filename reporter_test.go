package godu

import (
	"bufio"
	"bytes"
	"testing"
)

func TestReportTree(t *testing.T) {
	testTree := File{"b", 180, []File{
		File{"c", 100, []File{}},
		File{"d", 80, []File{
			File{"e", 50, []File{}},
			File{"f", 30, []File{}},
		}},
	}}
	expected := "c 100B\nd/ 80B\n"
	testTreeAgainstOutput(testTree, expected, t)
}

func TestReportingUnits(t *testing.T) {
	testTree := File{"X", 0, []File{
		File{"B", 1 << 0, []File{}},
		File{"K", 1 << 10, []File{}},
		File{"M", 1048576, []File{}},
		File{"G", 1073741824, []File{}},
		File{"T", 1099511627776, []File{}},
		File{"P", 1125899906842624, []File{}},
	}}
	ex := "B 1B\n"
	ex += "K 1K\n"
	ex += "M 1M\n"
	ex += "G 1G\n"
	ex += "T 1T\n"
	ex += "P 1P\n"
	testTreeAgainstOutput(testTree, ex, t)
}

func testTreeAgainstOutput(testTree File, expected string, t *testing.T) {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	ReportTree(testTree, writer)
	writer.Flush()
	result := buffer.String()
	if result != expected {
		t.Errorf("Result was not equal\n%sbut\n%s", expected, result)
	}
}
