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
				fakeFile{"g", 70, []fakeFile{ //thisfolder should get ignored
					fakeFile{"h", 10, []fakeFile{}},
					fakeFile{"i", 20, []fakeFile{}},
				}},
			}},
		}},
	}}
	ignoredFolders := map[string]struct{}{"g": struct{}{}}
	result := GetSubTree("b", createReadDir(testStructure), ignoredFolders)
	expected := File{"b", 180, true, []*File{
		&File{"c", 100, false, []*File{}},
		&File{"d", 80, true, []*File{
			&File{"e", 50, false, []*File{}},
			&File{"f", 30, false, []*File{}},
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
	result := GetSubTree("xyz", failing, map[string]struct{}{})
	if !reflect.DeepEqual(result, File{}) {
		t.Error("GetSubTree didn't return emtpy file on ReadDir failure")
	}
}
