package files

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// File structure representing files and folders with their accumulated sizes
// ReadError is set on folders whose contents could not be read
// (the node then represents an empty folder)
type File struct {
	Name      string
	Parent    *File
	Size      int64
	IsDir     bool
	Files     []*File
	ReadError string
}

// Path builds a file system location for given file
func (f *File) Path() string {
	if f.Parent == nil {
		return f.Name
	}
	return filepath.Join(f.Parent.Path(), f.Name)
}

// UpdateSize goes through subfiles and subfolders and accumulates their size
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

// ReadDir function can return list of files for given folder path
type ReadDir func(dirname string) ([]os.FileInfo, error)

// ShouldIgnoreFolder function decides whether a folder should be ignored
type ShouldIgnoreFolder func(absolutePath string) bool

func ignoringReadDir(shouldIgnore ShouldIgnoreFolder, originalReadDir ReadDir) ReadDir {
	return func(path string) ([]os.FileInfo, error) {
		if shouldIgnore(path) {
			return []os.FileInfo{}, nil
		}
		return originalReadDir(path)
	}
}

// WalkFolder will go through a given folder and subfolders and produces file structure
// with aggregated file sizes
func WalkFolder(
	path string,
	readDir ReadDir,
	ignoreFunction ShouldIgnoreFolder,
	progress chan<- int,
) *File {
	var wg sync.WaitGroup
	c := make(chan bool, 2*runtime.NumCPU())
	root := walkSubFolderConcurrently(path, nil, ignoringReadDir(ignoreFunction, readDir), c, &wg, progress)
	wg.Wait()
	close(progress)
	root.UpdateSize()
	return root
}

func walkSubFolderConcurrently(
	path string,
	parent *File,
	readDir ReadDir,
	c chan bool,
	wg *sync.WaitGroup,
	progress chan<- int,
) *File {
	result := &File{}
	entries, err := readDir(path)
	if err != nil {
		// an unreadable folder is represented as an empty folder so that
		// it stays a valid node in the tree (with a name and a parent);
		// ReadError lets the UI tell unreadable folders from empty ones
		log.Println(err)
		result.ReadError = err.Error()
		entries = nil
	}
	dirName, name := filepath.Split(path)
	result.Files = make([]*File, 0, len(entries))
	numSubFolders := 0
	defer updateProgress(progress, &numSubFolders)
	var mutex sync.Mutex
	for _, entry := range entries {
		if entry.IsDir() {
			numSubFolders++
			subFolderPath := filepath.Join(path, entry.Name())
			wg.Add(1)
			go func() {
				c <- true
				subFolder := walkSubFolderConcurrently(subFolderPath, result, readDir, c, wg, progress)
				mutex.Lock()
				result.Files = append(result.Files, subFolder)
				mutex.Unlock()
				<-c
				wg.Done()
			}()
		} else {
			size := entry.Size()
			file := &File{
				Name:   entry.Name(),
				Parent: result,
				Size:   size,
				Files:  []*File{},
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

func updateProgress(progress chan<- int, count *int) {
	if *count > 0 {
		progress <- *count
	}
}
