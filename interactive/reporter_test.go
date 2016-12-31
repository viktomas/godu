package interactive

import (
	"reflect"
	"testing"

	"github.com/viktomas/godu/core"
)

func TestReportTree(t *testing.T) {
	testTree := &core.File{"b", 180, []*core.File{
		&core.File{"c", 100, []*core.File{}},
		&core.File{"d", 80, []*core.File{
			&core.File{"e", 50, []*core.File{}},
			&core.File{"f", 30, []*core.File{}},
		}},
	}}
	expected := []string{"c 100B", "d/ 80B"}
	testTreeAgainstOutput(testTree, expected, t)
}

func TestReportingUnits(t *testing.T) {
	testTree := &core.File{"X", 0, []*core.File{
		&core.File{"B", 1 << 0, []*core.File{}},
		&core.File{"K", 1 << 10, []*core.File{}},
		&core.File{"M", 1048576, []*core.File{}},
		&core.File{"G", 1073741824, []*core.File{}},
		&core.File{"T", 1099511627776, []*core.File{}},
		&core.File{"P", 1125899906842624, []*core.File{}},
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

func testTreeAgainstOutput(testTree *core.File, expected []string, t *testing.T) {
	result := ReportTree(testTree)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected:\n%sbut got:\n%s", expected, result)
	}
}
