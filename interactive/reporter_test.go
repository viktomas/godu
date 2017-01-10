package interactive

import (
	"reflect"
	"testing"

	"github.com/viktomas/godu/core"
)

func TestReportTree(t *testing.T) {
	marked := make(map[*core.File]struct{})
	testTree := &core.File{"b", nil, 180, true, []*core.File{
		&core.File{"c", nil, 100, false, []*core.File{}},
		&core.File{"d", nil, 80, true, []*core.File{
			&core.File{"e", nil, 50, false, []*core.File{}},
			&core.File{"f", nil, 30, false, []*core.File{}},
		}},
	}}
	marked[testTree.Files[0]] = struct{}{}
	expected := []Line{
		Line{Text: []rune("* 100B c"), IsMarked: true},
		Line{Text: []rune("  80B d/")},
	}
	testTreeAgainstOutput(testTree, marked, expected, t)
}

func TestPrintsEmptyDir(t *testing.T) {
	marked := make(map[*core.File]struct{})
	testTree := &core.File{"", nil, 50, true, []*core.File{
		&core.File{"a", nil, 50, true, []*core.File{}},
	}}
	expected := []Line{
		Line{Text: []rune("  50B a/")},
	}
	testTreeAgainstOutput(testTree, marked, expected, t)
}

func TestFiveCharSize(t *testing.T) {
	marked := make(map[*core.File]struct{})
	testTree := &core.File{"X", nil, 0, true, []*core.File{
		&core.File{"o", nil, 1, false, []*core.File{}},
		&core.File{"on", nil, 10, false, []*core.File{}},
		&core.File{"one", nil, 100, false, []*core.File{}},
		&core.File{"one1", nil, 1000, false, []*core.File{}},
	}}
	ex := []Line{
		Line{Text: []rune("   1B o")},
		Line{Text: []rune("  10B on")},
		Line{Text: []rune(" 100B one")},
		Line{Text: []rune("1000B one1")},
	}
	testTreeAgainstOutput(testTree, marked, ex, t)
}

func TestReportingUnits(t *testing.T) {
	marked := make(map[*core.File]struct{})
	testTree := &core.File{"X", nil, 0, true, []*core.File{
		&core.File{"B", nil, 1 << 0, false, []*core.File{}},
		&core.File{"K", nil, 1 << 10, false, []*core.File{}},
		&core.File{"M", nil, 1048576, false, []*core.File{}},
		&core.File{"G", nil, 1073741824, false, []*core.File{}},
		&core.File{"T", nil, 1099511627776, false, []*core.File{}},
		&core.File{"P", nil, 1125899906842624, false, []*core.File{}},
	}}
	ex := []Line{
		Line{Text: []rune("   1B B")},
		Line{Text: []rune("   1K K")},
		Line{Text: []rune("   1M M")},
		Line{Text: []rune("   1G G")},
		Line{Text: []rune("   1T T")},
		Line{Text: []rune("   1P P")},
	}
	testTreeAgainstOutput(testTree, marked, ex, t)
}

func testTreeAgainstOutput(testTree *core.File, marked map[*core.File]struct{}, expected []Line, t *testing.T) {
	result := ReportTree(testTree, marked)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected:\n%sbut got:\n%s", expected, result)
	}
}
