package interactive

import (
	"fmt"

	"github.com/viktomas/godu/core"
)

// Line represents row of text in folder UI contains info about subfile
type Line struct {
	Text     []rune
	IsMarked bool
}

// Status contain info about size of all files in current godu instance
// and size of the files marked by user
type Status struct {
	Total    string
	Selected string
}

// ReportStatus reads through the folder structure and produces Status
func ReportStatus(file *core.File, markedFiles *map[*core.File]struct{}) Status {
	parent := file
	for parent.Parent != nil {
		parent = parent.Parent
	}
	var selected int64
	for f := range *markedFiles {
		if !parentMarked(f, markedFiles) {
			selected += f.Size
		}
	}
	return Status{
		Total:    fmt.Sprintf("Total size: %s", formatBytes(parent.Size)),
		Selected: fmt.Sprintf("Selected size: %s", formatBytes(selected)),
	}
}

func parentMarked(file *core.File, markedFiles *map[*core.File]struct{}) bool {
	parent := file
	for parent.Parent != nil {
		_, found := (*markedFiles)[parent.Parent]
		if found {
			return true
		}
		parent = parent.Parent
	}
	return false
}

// ReportFolder converts all subfiles into UI lines
func ReportFolder(folder *core.File, markedFiles map[*core.File]struct{}) []Line {
	report := make([]Line, len(folder.Files))
	for index, file := range folder.Files {
		name := file.Name
		if file.IsDir {
			name = name + "/"
		}
		marking := " "
		_, isMarked := markedFiles[file]
		if isMarked {
			marking = "*"
		}
		report[index] = Line{
			Text:     []rune(fmt.Sprintf("%s%s %s", marking, formatBytes(file.Size), name)),
			IsMarked: isMarked,
		}
	}
	return report
}

func formatBytes(bytesInt int64) string {
	bytes := float32(bytesInt)
	var unit string
	var amount float32
	switch {
	case core.PETABYTE <= bytes:
		unit = "P"
		amount = bytes / core.PETABYTE
	case core.TERABYTE <= bytes:
		unit = "T"
		amount = bytes / core.TERABYTE
	case core.GIGABYTE <= bytes:
		unit = "G"
		amount = bytes / core.GIGABYTE
	case core.MEGABYTE <= bytes:
		unit = "M"
		amount = bytes / core.MEGABYTE
	case core.KILOBYTE <= bytes:
		unit = "K"
		amount = bytes / core.KILOBYTE
	default:
		unit = "B"
		amount = bytes

	}
	return fmt.Sprintf("%4.0f%s", amount, unit)
}
