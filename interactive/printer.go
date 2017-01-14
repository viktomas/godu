package interactive

import (
	"fmt"
	"io"

	"github.com/viktomas/godu/core"
)

func PrintMarkedFiles(state *core.State, writer io.Writer) {
	for file := range state.MarkedFiles {
		fmt.Fprintf(writer, "'%s'\n", file.Path())
	}
}
