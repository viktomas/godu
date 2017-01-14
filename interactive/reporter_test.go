package interactive

import (
	"reflect"
	"testing"

	"github.com/viktomas/godu/core"
)

func TestReportTree(t *testing.T) {
	marked := make(map[*core.File]struct{})
	testTree := core.NewTestFolder("b",
		core.NewTestFile("c", 100),
		core.NewTestFolder("d",
			core.NewTestFile("e", 50),
			core.NewTestFile("f", 30),
		),
	)
	marked[testTree.Files[0]] = struct{}{}
	expected := []Line{
		Line{Text: []rune("* 100B c"), IsMarked: true},
		Line{Text: []rune("   80B d/")},
	}
	testTreeAgainstOutput(testTree, marked, expected, t)
}

func TestPrintsEmptyDir(t *testing.T) {
	marked := make(map[*core.File]struct{})
	testTree := core.NewTestFolder("", core.NewTestFolder("a"))
	expected := []Line{
		Line{Text: []rune("    0B a/")},
	}
	testTreeAgainstOutput(testTree, marked, expected, t)
}

func TestFiveCharSize(t *testing.T) {
	marked := make(map[*core.File]struct{})
	testTree := core.NewTestFolder("X",
		core.NewTestFile("o", 1),
		core.NewTestFile("on", 10),
		core.NewTestFile("one", 100),
		core.NewTestFile("one1", 1000),
	)
	ex := []Line{
		Line{Text: []rune("    1B o")},
		Line{Text: []rune("   10B on")},
		Line{Text: []rune("  100B one")},
		Line{Text: []rune(" 1000B one1")},
	}
	testTreeAgainstOutput(testTree, marked, ex, t)
}

func TestReportingUnits(t *testing.T) {
	marked := make(map[*core.File]struct{})
	testTree := core.NewTestFolder("X",
		core.NewTestFile("B", 1<<0),
		core.NewTestFile("K", 1<<10),
		core.NewTestFile("M", 1048576),
		core.NewTestFile("G", 1073741824),
		core.NewTestFile("T", 1099511627776),
		core.NewTestFile("P", 1125899906842624),
	)
	ex := []Line{
		Line{Text: []rune("    1B B")},
		Line{Text: []rune("    1K K")},
		Line{Text: []rune("    1M M")},
		Line{Text: []rune("    1G G")},
		Line{Text: []rune("    1T T")},
		Line{Text: []rune("    1P P")},
	}
	testTreeAgainstOutput(testTree, marked, ex, t)
}

func testTreeAgainstOutput(testTree *core.File, marked map[*core.File]struct{}, expected []Line, t *testing.T) {
	result := ReportTree(testTree, marked)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected:\n%sbut got:\n%s", expected, result)
	}
}
