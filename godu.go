package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/viktomas/godu/core"
	"github.com/viktomas/godu/interactive"
)

func main() {
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}
	tree := core.GetSubTree(roots[0], ioutil.ReadDir, getIgnoredFolders())
	interactive.InteractiveTree(&tree, os.Stdout, os.Stdin, 10*1024*1024)
}
