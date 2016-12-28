package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type File struct {
	name  string
	size  int64
	files []File
}

type ReadDir func(dirname string) ([]os.FileInfo, error)

func main() {
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}
	file := getSubTree(roots[0], ioutil.ReadDir)
	fmt.Printf("%v", file)
}

func getSubTree(path string, readDir ReadDir) File {
	_, name := filepath.Split(path)
	entries, err := readDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "godu: %v\n", err)
		return File{}
	}
	files := make([]File, len(entries))
	for i, entry := range entries {
		if entry.IsDir() {
			subDir := filepath.Join(path, entry.Name())
			files[i] = getSubTree(subDir, readDir)
		} else {
			files[i] = File{
				entry.Name(),
				entry.Size(),
				[]File{},
			}
		}
	}
	return File{name, 0, files}
}
