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

func TestGetSubTreeOnSimpleDir(t *testing.T) {
	testStructure := fakeFile{"a", 0, []fakeFile{
		fakeFile{"b", 0, []fakeFile{
			fakeFile{"c", 100, []fakeFile{}},
			fakeFile{"d", 0, []fakeFile{
				fakeFile{"e", 50, []fakeFile{}},
				fakeFile{"f", 30, []fakeFile{}},
			}},
		}},
	}}
	result := GetSubTree("b", createReadDir(testStructure))
	expected := File{"b", 180, []*File{
		&File{"c", 100, []*File{}},
		&File{"d", 80, []*File{
			&File{"e", 50, []*File{}},
			&File{"f", 30, []*File{}},
		}},
	}}
	if !reflect.DeepEqual(result, expected) {
		t.Error("file tree wasn't walked correctly")
		fmt.Printf("expected: %v", expected)
		fmt.Printf("result: %v", result)
	}

}

func TestGetSubTreeHandlesError(t *testing.T) {
	failing := func(path string) ([]os.FileInfo, error) {
		return []os.FileInfo{}, errors.New("Not found")
	}
	result := GetSubTree("xyz", failing)
	if !reflect.DeepEqual(result, File{}) {
		t.Error("GetSubTree didn't return emtpy file on ReadDir failure")
	}
}

func TestPruneTree(t *testing.T) {
	testTree := &File{"b", 260, []*File{
		&File{"c", 100, []*File{}},
		&File{"d", 160, []*File{
			&File{"e", 50, []*File{}},
			&File{"f", 30, []*File{}},
			&File{"g", 80, []*File{
				&File{"i", 50, []*File{}},
				&File{"j", 30, []*File{}},
			}},
		}},
	}}
	expected := &File{"b", 260, []*File{
		&File{"c", 100, []*File{}},
		&File{"d", 160, []*File{
			&File{"g", 80, []*File{}},
		}},
	}}
	PruneTree(testTree, 60)
	if !reflect.DeepEqual(testTree, expected) {
		t.Errorf("tree not pruned correctly\ntree %v\nnot as expected %v", testTree, expected)
	}

}
