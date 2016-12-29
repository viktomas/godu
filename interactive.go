package godu

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

func InteractiveTree(tree File, w io.Writer, r io.Reader) {
	printOptions(tree, w)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		option, _ := strconv.Atoi(scanner.Text())
		tree := tree.Files[option]
		printOptions(tree, w)
	}
}

func printOptions(tree File, w io.Writer) {
	lines := ReportTree(tree)
	for index, line := range lines {
		fmt.Fprintf(w, "%d) %s\n", index, line)
	}
}
