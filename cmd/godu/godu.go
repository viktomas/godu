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
	// file := getSubTree(roots[0], ioutil.ReadDir)
	// fmt.Printf("%v", file)
	tree := godu.GetSubTree(roots[0], ioutil.ReadDir)
	godu.ReportTree(tree, os.Stdout)
}
