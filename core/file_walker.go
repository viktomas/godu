package core

import (
	"fmt"
	"os"
	"path/filepath"
)

type File struct {
	Name  string
	Size  int64
	IsDir bool
	Files []*File
}

type ReadDir func(dirname string) ([]os.FileInfo, error)

func GetSubTree(path string, readDir ReadDir, ignoredFolders map[string]struct{}) File {
	_, name := filepath.Split(path)
	entries, err := readDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "godu: %v\n", err)
		return File{}
	}
	files := make([]*File, 0, len(entries))
	var folderSize int64
	for _, entry := range entries {
		if entry.IsDir() {
			_, ignored := ignoredFolders[entry.Name()]
			if ignored {
				continue
			}
			subDir := filepath.Join(path, entry.Name())
			subfolder := GetSubTree(subDir, readDir, ignoredFolders)
			folderSize += subfolder.Size
			files = append(files, &subfolder)
		} else {
			size := entry.Size()
			folderSize += size
			file := &File{
				entry.Name(),
				size,
				false,
				[]*File{},
			}
			files = append(files, file)
		}
	}
	return File{name, folderSize, true, files}
}
