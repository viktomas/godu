package interactive

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/viktomas/godu/files"
)

func TestReportFolder(t *testing.T) {
	marked := make(map[*files.File]struct{})
	folder := files.NewTestFolder("b",
		files.NewTestFile("c", 100),
		files.NewTestFolder("d",
			files.NewTestFile("e", 50),
			files.NewTestFile("f", 30),
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
	marked := make(map[*files.File]struct{})
	folder := files.NewTestFolder("", files.NewTestFolder("a"))
	expected := []Line{
		Line{Text: []rune("    0B a/")},
	}
	testFolderAgainstOutput(folder, marked, expected, t)
}

func TestFiveCharSize(t *testing.T) {
	marked := make(map[*files.File]struct{})
	folder := files.NewTestFolder("X",
		files.NewTestFile("o", 1),
		files.NewTestFile("on", 10),
		files.NewTestFile("one", 100),
		files.NewTestFile("one1", 1000),
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
	marked := make(map[*files.File]struct{})
	folder := files.NewTestFolder("X",
		files.NewTestFile("B", 1<<0),
		files.NewTestFile("K", 1<<10),
		files.NewTestFile("M", 1048576),
		files.NewTestFile("G", 1073741824),
		files.NewTestFile("T", 1099511627776),
		files.NewTestFile("P", 1125899906842624),
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

func testFolderAgainstOutput(folder *files.File, marked map[*files.File]struct{}, expected []Line, t *testing.T) {
	result := ReportFolder(folder, marked)
	assert.Equal(t, expected, result)
}

func TestReportStatusTotalSize(t *testing.T) {
	folder := files.NewTestFolder("root",
		files.NewTestFile("a", 10),
		files.NewTestFile("b", 20),
		files.NewTestFile("c", 30),
	)
	marked := make(map[*files.File]struct{})
	status := ReportStatus(files.FindTestFile(folder, "b"), &marked)
	assert.Equal(t, "Total size:   60B", status.Total)
}

func TestReportStatusSelectedSize(t *testing.T) {
	folder := files.NewTestFolder("root",
		files.NewTestFile("a", 30),
		files.NewTestFile("b", 40),
	)
	marked := make(map[*files.File]struct{})
	marked[files.FindTestFile(folder, "b")] = struct{}{}
	status := ReportStatus(folder, &marked)
	assert.Equal(t, "Selected size:   40B", status.Selected)
}

func TestReportStatusSelectedSizeWithParent(t *testing.T) {
	folder := files.NewTestFolder("root",
		files.NewTestFolder("f1",
			files.NewTestFile("a", 10),
			files.NewTestFile("b", 20),
		),
		files.NewTestFile("c", 30),
		files.NewTestFile("d", 40),
	)
	marked := make(map[*files.File]struct{})
	marked[files.FindTestFile(folder, "f1")] = struct{}{}
	marked[files.FindTestFile(folder, "a")] = struct{}{}
	status := ReportStatus(folder, &marked)
	assert.Equal(t, "Selected size:   30B", status.Selected)
}
