package main

import (
	"flag"
	"fmt"
	"github.com/viktomas/godu"
	"io/ioutil"
)

func main() {
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}
	// file := getSubTree(roots[0], ioutil.ReadDir)
	// fmt.Printf("%v", file)
	godu.GetSubTree(roots[0], ioutil.ReadDir)
	fmt.Print("Done")
}
