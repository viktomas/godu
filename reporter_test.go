package godu

import (
	"reflect"
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
	expected := []string{"c 100B", "d/ 80B"}
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
	ex := []string{
		"B 1B",
		"K 1K",
		"M 1M",
		"G 1G",
		"T 1T",
		"P 1P",
	}
	testTreeAgainstOutput(testTree, ex, t)
}

func testTreeAgainstOutput(testTree File, expected []string, t *testing.T) {
	result := ReportTree(testTree)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected:\n%sbut got:\n%s", expected, result)
	}
}
