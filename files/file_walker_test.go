package files

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type fakeFile struct {
	fileName  string
	fileSize  int64
	fakeFiles []fakeFile
}

func (f fakeFile) Name() string       { return f.fileName }
func (f fakeFile) Size() int64        { return f.fileSize }
func (f fakeFile) Mode() os.FileMode  { return 0 }
func (f fakeFile) ModTime() time.Time { return time.Now() }
func (f fakeFile) IsDir() bool        { return len(f.fakeFiles) > 0 }
func (f fakeFile) Sys() interface{}   { return nil }

func createReadDir(ff fakeFile) ReadDir {
	return func(path string) ([]os.FileInfo, error) {
		names := strings.Split(path, "/")
		fakeFolder := ff
		var found bool
		for _, name := range names {
			found = false
			for _, testFile := range fakeFolder.fakeFiles {
				if testFile.fileName == name {
					fakeFolder = testFile
					found = true
					break
				}
			}
			if !found {
				return []os.FileInfo{}, fmt.Errorf("file not found")
			}

		}
		result := make([]os.FileInfo, len(fakeFolder.fakeFiles))
		for i, resultFile := range fakeFolder.fakeFiles {
			result[i] = resultFile
		}
		return result, nil
	}
}

func TestFilePath(t *testing.T) {
	root := NewTestFolder("root",
		NewTestFolder("folder1",
			NewTestFile("file1", 0),
		),
	)
	want := filepath.Join("root", "folder1", "file1")
	file1 := FindTestFile(root, "file1")
	assert.Equal(t, want, file1.Path())
}

func TestWalkFolderOnSimpleDir(t *testing.T) {
	testStructure := fakeFile{"a", 0, []fakeFile{
		{"b", 0, []fakeFile{
			{"c", 100, []fakeFile{}},
			{"d", 0, []fakeFile{
				{"e", 50, []fakeFile{}},
				{"f", 30, []fakeFile{}},
				{"g", 70, []fakeFile{ //thisfolder should get ignored
					{"h", 10, []fakeFile{}},
					{"i", 20, []fakeFile{}},
				}},
			}},
		}},
	}}
	dummyIgnoreFunction := func(p string) bool { return p == filepath.Join("b", "d", "g") }
	progress := make(chan int, 3)
	result := WalkFolder("b", createReadDir(testStructure), dummyIgnoreFunction, progress)
	buildExpected := func() *File {
		b := &File{"b", nil, 180, true, []*File{}}
		c := &File{"c", b, 100, false, []*File{}}
		d := &File{"d", b, 80, true, []*File{}}
		b.Files = []*File{c, d}

		e := &File{"e", nil, 50, false, []*File{}}
		e.Parent = d
		f := &File{"f", nil, 30, false, []*File{}}
		g := &File{"g", nil, 0, true, []*File{}}
		f.Parent = d
		g.Parent = d
		d.Files = []*File{e, f, g}

		return b
	}
	expected := buildExpected()
	assert.Equal(t, expected, result)
	resultProgress := 0
	resultProgress += <-progress
	resultProgress += <-progress
	_, more := <-progress
	assert.Equal(t, 2, resultProgress)
	assert.False(t, more, "the progress channel should be closed")
}

func TestWalkFolderHandlesError(t *testing.T) {
	failing := func(path string) ([]os.FileInfo, error) {
		return []os.FileInfo{}, errors.New("Not found")
	}
	progress := make(chan int, 2)
	result := WalkFolder("xyz", failing, func(string) bool { return false }, progress)
	assert.Equal(t, File{}, *result, "WalkFolder didn't return empty file on ReadDir failure")
}
