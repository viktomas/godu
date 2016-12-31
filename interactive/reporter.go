package interactive

import (
	"fmt"
	"github.com/viktomas/godu/core"
)

func ReportTree(folder *core.File) []string {
	report := make([]string, len(folder.Files))
	for index, file := range folder.Files {
		name := file.Name
		if len(file.Files) > 0 {
			name = name + "/"
		}
		report[index] = fmt.Sprintf("%s %s", name, formatBytes(file.Size))
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
	return fmt.Sprintf("%.0f%s", amount, unit)
}
