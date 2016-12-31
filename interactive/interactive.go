package interactive

import (
	"bufio"
	"fmt"
	"github.com/viktomas/godu/core"
	"io"
	"sort"
	"strconv"
)

type bySizeDesc []*core.File

func (f bySizeDesc) Len() int           { return len(f) }
func (f bySizeDesc) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f bySizeDesc) Less(i, j int) bool { return f[i].Size > f[j].Size }

func InteractiveTree(tree *core.File, w io.Writer, r io.Reader, limit int64) {
	core.PruneTree(tree, limit)
	sort.Sort(bySizeDesc(tree.Files))
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
		fmt.Fprintf(w, "%d) %s\n", index, line)
	}
}
