package core

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type File struct {
	Name   string
	Parent *File
	Size   int64
	IsDir  bool
	Files  []*File
}

func (f *File) Path() string {
	if f.Parent == nil {
		return f.Name
	}
	return filepath.Join(f.Parent.Path(), f.Name)
}

func (f *File) UpdateSize() {
	if !f.IsDir {
		return
	}
	var size int64
	for _, child := range f.Files {
		child.UpdateSize()
		size += child.Size
	}
	f.Size = size
}

type ReadDir func(dirname string) ([]os.FileInfo, error)

func WalkFolder(path string, readDir ReadDir, ignoredFolders map[string]struct{}) *File {
	var wg sync.WaitGroup
	c := make(chan bool, 2*runtime.NumCPU())
	root := walkSubFolderConcurrently(path, nil, readDir, ignoredFolders, c, &wg)
	wg.Wait()
	root.UpdateSize()
	return root
}

func walkSubFolderConcurrently(path string, parent *File, readDir ReadDir, ignoredFolders map[string]struct{}, c chan bool, wg *sync.WaitGroup) *File {
	result := &File{}
	entries, err := readDir(path)
	if err != nil {
		log.Println(err)
		return result
	}
	dirName, name := filepath.Split(path)
	result.Files = make([]*File, 0, len(entries))
	var mutex sync.Mutex
	for _, entry := range entries {
		if entry.IsDir() {
			if _, ignored := ignoredFolders[entry.Name()]; ignored {
				continue
			}
			subFolderPath := filepath.Join(path, entry.Name())
			wg.Add(1)
			go func() {
				c <- true
				subFolder := walkSubFolderConcurrently(subFolderPath, result, readDir, ignoredFolders, c, wg)
				mutex.Lock()
				result.Files = append(result.Files, subFolder)
				mutex.Unlock()
				<-c
				wg.Done()
			}()
		} else {
			size := entry.Size()
			file := &File{
				entry.Name(),
				result,
				size,
				false,
				[]*File{},
			}
			mutex.Lock()
			result.Files = append(result.Files, file)
			mutex.Unlock()
		}
	}
	if parent != nil {
		result.Name = name
		result.Parent = parent
	} else {
		// Root dir
		// TODO unit test this Join
		result.Name = filepath.Join(dirName, name)
	}
	result.IsDir = true
	return result
}
