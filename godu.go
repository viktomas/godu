package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/viktomas/godu/core"
	"github.com/viktomas/godu/interactive"
)

func main() {
	limit := flag.Int64("l", 10, "show only files larger than limit (in MB)")
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}
	tree := core.GetSubTree(roots[0], ioutil.ReadDir, getIgnoredFolders())
	interactive.InteractiveTree(&tree, os.Stdout, os.Stdin, *limit*core.MEGABYTE)
}
