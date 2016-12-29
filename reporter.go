package godu

import (
	"fmt"
	"io"
)

const (
	kilobyte = 1024
	megabyte = 1024 * kilobyte
	gigabyte = 1024 * megabyte
	terabyte = 1024 * gigabyte
	petabyte = 1024 * terabyte
)

func ReportTree(folder File, w io.Writer) {
	for _, file := range folder.Files {
		name := file.Name
		if len(file.Files) > 0 {
			name = name + "/"
		}
		fmt.Fprintf(w, "%s %s\n", name, formatBytes(file.Size))
	}
}

func formatBytes(bytesInt int64) string {
	bytes := float32(bytesInt)
	var unit string
	var amount float32
	switch {
	case petabyte <= bytes:
		unit = "P"
		amount = bytes / petabyte
	case terabyte <= bytes:
		unit = "T"
		amount = bytes / terabyte
	case gigabyte <= bytes:
		unit = "G"
		amount = bytes / gigabyte
	case megabyte <= bytes:
		unit = "M"
		amount = bytes / megabyte
	case kilobyte <= bytes:
		unit = "K"
		amount = bytes / kilobyte
	default:
		unit = "B"
		amount = bytes

	}
	return fmt.Sprintf("%.0f%s", amount, unit)
}
