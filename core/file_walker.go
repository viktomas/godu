package core

import (
	"log"
	"os"
	"path/filepath"
)

type File struct {
	Name   string
	Parent *File
	Size   int64
	IsDir  bool
	Files  []*File
}

func (f *File) Path() string {
	var path string
	cur := f
	for {
		if cur = cur.Parent; cur == nil {
			break
		}
		path = filepath.Join(cur.Name, path)
	}
	return filepath.Join(path, f.Name)
}

type ReadDir func(dirname string) ([]os.FileInfo, error)

func GetSubTree(path string, parent *File, readDir ReadDir, ignoredFolders map[string]struct{}) *File {
	_, name := filepath.Split(path)
	ret := &File{}
	entries, err := readDir(path)
	if err != nil {
		log.Println(err)
		return ret
	}
	files := make([]*File, 0, len(entries))
	var folderSize int64
	for _, entry := range entries {
		if entry.IsDir() {
			if _, ignored := ignoredFolders[entry.Name()]; ignored {
				continue
			}
			subDir := filepath.Join(path, entry.Name())
			subfolder := GetSubTree(subDir, ret, readDir, ignoredFolders)
			folderSize += subfolder.Size
			files = append(files, subfolder)
		} else {
			size := entry.Size()
			folderSize += size
			file := &File{
				entry.Name(),
				ret,
				size,
				false,
				[]*File{},
			}
			files = append(files, file)
		}
	}
	ret.Name = name
	if parent != nil {
		ret.Parent = parent
	}
	ret.Size = folderSize
	ret.IsDir = true
	ret.Files = files
	return ret
}
