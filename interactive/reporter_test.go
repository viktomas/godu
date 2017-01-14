package interactive

import (
	"reflect"
	"testing"

	"github.com/viktomas/godu/core"
)

func TestReportTree(t *testing.T) {
	marked := make(map[*core.File]struct{})
	testTree := core.TstFolder("b",
		core.TstFile("c", 100),
		core.TstFolder("d",
			core.TstFile("e", 50),
			core.TstFile("f", 30),
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
	testTree := core.TstFolder("", core.TstFolder("a"))
	expected := []Line{
		Line{Text: []rune("    0B a/")},
	}
	testTreeAgainstOutput(testTree, marked, expected, t)
}

func TestFiveCharSize(t *testing.T) {
	marked := make(map[*core.File]struct{})
	testTree := core.TstFolder("X",
		core.TstFile("o", 1),
		core.TstFile("on", 10),
		core.TstFile("one", 100),
		core.TstFile("one1", 1000),
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
	testTree := core.TstFolder("X",
		core.TstFile("B", 1<<0),
		core.TstFile("K", 1<<10),
		core.TstFile("M", 1048576),
		core.TstFile("G", 1073741824),
		core.TstFile("T", 1099511627776),
		core.TstFile("P", 1125899906842624),
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
