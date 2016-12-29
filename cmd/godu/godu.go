package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/viktomas/godu"
)

func main() {
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}
	tree := godu.GetSubTree(roots[0], ioutil.ReadDir)
	godu.InteractiveTree(&tree, os.Stdout, os.Stdin, 10*1024*1024)
}
