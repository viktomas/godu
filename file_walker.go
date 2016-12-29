package godu

import (
	"fmt"
	"os"
	"path/filepath"
)

type File struct {
	Name  string
	Size  int64
	Files []File
}

type ReadDir func(dirname string) ([]os.FileInfo, error)

func GetSubTree(path string, readDir ReadDir) File {
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
			files[i] = GetSubTree(subDir, readDir)
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
