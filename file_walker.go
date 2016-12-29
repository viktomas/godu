package godu

import (
	"fmt"
	"os"
	"path/filepath"
)

type File struct {
	Name  string
	Size  int64
	Files []*File
}

type ReadDir func(dirname string) ([]os.FileInfo, error)

func GetSubTree(path string, readDir ReadDir) File {
	_, name := filepath.Split(path)
	entries, err := readDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "godu: %v\n", err)
		return File{}
	}
	files := make([]*File, len(entries))
	var folderSize int64
	for i, entry := range entries {
		if entry.IsDir() {
			subDir := filepath.Join(path, entry.Name())
			subfolder := GetSubTree(subDir, readDir)
			folderSize += subfolder.Size
			files[i] = &subfolder
		} else {
			size := entry.Size()
			files[i] = &File{
				entry.Name(),
				size,
				[]*File{},
			}
			folderSize += size
		}
	}
	return File{name, folderSize, files}
}

func PruneTree(tree *File, limit int64) {
	prunedFiles := []*File{}
	for _, file := range tree.Files {
		if file.Size >= limit {
			PruneTree(file, limit)
			prunedFiles = append(prunedFiles, file)
		}
	}
	tree.Files = prunedFiles

}
