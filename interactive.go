package godu

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strconv"
)

type bySizeDesc []File

func (f bySizeDesc) Len() int           { return len(f) }
func (f bySizeDesc) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f bySizeDesc) Less(i, j int) bool { return f[i].Size > f[j].Size }

func InteractiveTree(tree File, w io.Writer, r io.Reader) {
	sort.Sort(bySizeDesc(tree.Files))
	printOptions(tree, w)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		option, _ := strconv.Atoi(scanner.Text())
		tree := tree.Files[option]
		sort.Sort(bySizeDesc(tree.Files))
		printOptions(tree, w)
	}
}

func printOptions(tree File, w io.Writer) {
	lines := ReportTree(tree)
	for index, line := range lines {
		fmt.Fprintf(w, "%d) %s\n", index, line)
	}
}
