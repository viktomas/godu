package core

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
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
	want := "root/folder1/file1"
	file1 := FindTestFile(root, "file1")
	if p := file1.Path(); p != want {
		t.Errorf("unexpected file path, got '%s', want '%s'", p, want)
	}
}

func TestGetSubTreeOnSimpleDir(t *testing.T) {
	testStructure := fakeFile{"a", 0, []fakeFile{
		fakeFile{"b", 0, []fakeFile{
			fakeFile{"c", 100, []fakeFile{}},
			fakeFile{"d", 0, []fakeFile{
				fakeFile{"e", 50, []fakeFile{}},
				fakeFile{"f", 30, []fakeFile{}},
				fakeFile{"g", 70, []fakeFile{ //thisfolder should get ignored
					fakeFile{"h", 10, []fakeFile{}},
					fakeFile{"i", 20, []fakeFile{}},
				}},
			}},
		}},
	}}
	ignoredFolders := map[string]struct{}{"g": struct{}{}}
	progress := make(chan int, 0)
	go dummyProgressConsumer(progress)
	result := GetSubTree("b", nil, createReadDir(testStructure), ignoredFolders, progress)
	buildExpected := func() *File {
		b := &File{"b", nil, 180, true, []*File{}}
		c := &File{"c", b, 100, false, []*File{}}
		d := &File{"d", b, 80, true, []*File{}}
		b.Files = []*File{c, d}

		e := &File{"e", nil, 50, false, []*File{}}
		e.Parent = d
		f := &File{"f", nil, 30, false, []*File{}}
		f.Parent = d
		d.Files = []*File{e, f}

		return b
	}
	expected := buildExpected()
	if !reflect.DeepEqual(*result, *expected) {
		t.Error("file tree wasn't walked correctly")
		fmt.Printf("expected: %v", *expected)
		fmt.Printf("result: %v", *result)
	}
}

func TestGetSubTreeHandlesError(t *testing.T) {
	failing := func(path string) ([]os.FileInfo, error) {
		return []os.FileInfo{}, errors.New("Not found")
	}
	progress := make(chan int, 0)
	go dummyProgressConsumer(progress)
	result := GetSubTree("xyz", nil, failing, map[string]struct{}{}, progress)
	if !reflect.DeepEqual(*result, File{}) {
		t.Error("GetSubTree didn't return emtpy file on ReadDir failure")
	}
}

func dummyProgressConsumer(progress <-chan int) {
	for range progress {
	}
}
