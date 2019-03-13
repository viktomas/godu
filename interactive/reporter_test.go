package interactive

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/viktomas/godu/core"
)

func TestReportFolder(t *testing.T) {
	marked := make(map[*core.File]struct{})
	folder := core.NewTestFolder("b",
		core.NewTestFile("c", 100),
		core.NewTestFolder("d",
			core.NewTestFile("e", 50),
			core.NewTestFile("f", 30),
		),
	)
	marked[folder.Files[0]] = struct{}{}
	expected := []Line{
		Line{Text: []rune("* 100B c"), IsMarked: true},
		Line{Text: []rune("   80B d/")},
	}
	testFolderAgainstOutput(folder, marked, expected, t)
}

func TestPrintsEmptyDir(t *testing.T) {
	marked := make(map[*core.File]struct{})
	folder := core.NewTestFolder("", core.NewTestFolder("a"))
	expected := []Line{
		Line{Text: []rune("    0B a/")},
	}
	testFolderAgainstOutput(folder, marked, expected, t)
}

func TestFiveCharSize(t *testing.T) {
	marked := make(map[*core.File]struct{})
	folder := core.NewTestFolder("X",
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
	testFolderAgainstOutput(folder, marked, ex, t)
}

func TestReportingUnits(t *testing.T) {
	marked := make(map[*core.File]struct{})
	folder := core.NewTestFolder("X",
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
	testFolderAgainstOutput(folder, marked, ex, t)
}

func testFolderAgainstOutput(folder *core.File, marked map[*core.File]struct{}, expected []Line, t *testing.T) {
	result := ReportFolder(folder, marked)
	assert.Equal(t, expected, result)
}

func TestReportStatusTotalSize(t *testing.T) {
	folder := core.NewTestFolder("root",
		core.NewTestFile("a", 10),
		core.NewTestFile("b", 20),
		core.NewTestFile("c", 30),
	)
	marked := make(map[*core.File]struct{})
	status := ReportStatus(core.FindTestFile(folder, "b"), &marked)
	assert.Equal(t, "Total size:   60B", status.Total)
}

func TestReportStatusSelectedSize(t *testing.T) {
	folder := core.NewTestFolder("root",
		core.NewTestFile("a", 30),
		core.NewTestFile("b", 40),
	)
	marked := make(map[*core.File]struct{})
	marked[core.FindTestFile(folder, "b")] = struct{}{}
	status := ReportStatus(folder, &marked)
	assert.Equal(t, "Selected size:   40B", status.Selected)
}

func TestReportStatusSelectedSizeWithParent(t *testing.T) {
	folder := core.NewTestFolder("root",
		core.NewTestFolder("f1",
			core.NewTestFile("a", 10),
			core.NewTestFile("b", 20),
		),
		core.NewTestFile("c", 30),
		core.NewTestFile("d", 40),
	)
	marked := make(map[*core.File]struct{})
	marked[core.FindTestFile(folder, "f1")] = struct{}{}
	marked[core.FindTestFile(folder, "a")] = struct{}{}
	status := ReportStatus(folder, &marked)
	assert.Equal(t, "Selected size:   30B", status.Selected)
}
