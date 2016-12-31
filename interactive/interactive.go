package interactive

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"github.com/viktomas/godu/core"
)

func InteractiveTree(tree *core.File, w io.Writer, r io.Reader, limit int64) {
	core.PruneTree(tree, limit)
	core.SortDesc(tree)
	state := core.State{
		Folder: tree,
	}
	printOptions(state, w)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		option, err := strconv.Atoi(scanner.Text())
		if err != nil || option > len(tree.Files) {
			back := core.GoBack{}
			state, _ = back.Execute(state)
		} else {
			enter := core.Enter{option}
			state, _ = enter.Execute(state)
		}
		printOptions(state, w)
	}
}

func printOptions(state core.State, w io.Writer) {
	fmt.Fprintf(w, "%s\n---\n", state.Path())
	lines := ReportTree(state.Folder)
	for index, line := range lines {
		fmt.Fprintf(w, "%d)\t%s\n", index, line)
	}
}
