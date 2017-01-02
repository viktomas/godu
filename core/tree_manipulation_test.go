package core

import (
	"reflect"
	"testing"
)

func TestSortTree(t *testing.T) {
	testTree := &File{"b", 260, true, []*File{
		&File{"c", 100, false, []*File{}},
		&File{"d", 160, true, []*File{
			&File{"e", 50, false, []*File{}},
			&File{"f", 30, false, []*File{}},
			&File{"g", 80, true, []*File{
				&File{"i", 30, false, []*File{}},
				&File{"j", 50, false, []*File{}},
			}},
		}},
	}}
	expected := &File{"b", 260, true, []*File{
		&File{"d", 160, true, []*File{
			&File{"g", 80, true, []*File{
				&File{"j", 50, false, []*File{}},
				&File{"i", 30, false, []*File{}},
			}},
			&File{"e", 50, false, []*File{}},
			&File{"f", 30, false, []*File{}},
		}},
		&File{"c", 100, false, []*File{}},
	}}
	SortDesc(testTree)
	if !reflect.DeepEqual(testTree, expected) {
		t.Errorf("tree not sorted correctly\ntree %v\nnot as expected %v", testTree, expected)
	}
}

func TestPruneTree(t *testing.T) {
	testTree := &File{"b", 260, true, []*File{
		&File{"c", 100, false, []*File{}},
		&File{"d", 160, true, []*File{
			&File{"e", 50, false, []*File{}},
			&File{"f", 30, false, []*File{}},
			&File{"g", 80, true, []*File{
				&File{"i", 50, false, []*File{}},
				&File{"j", 30, false, []*File{}},
			}},
		}},
	}}
	expected := &File{"b", 260, true, []*File{
		&File{"c", 100, false, []*File{}},
		&File{"d", 160, true, []*File{
			&File{"g", 80, true, []*File{}},
		}},
	}}
	PruneTree(testTree, 60)
	if !reflect.DeepEqual(testTree, expected) {
		t.Errorf("tree not pruned correctly\ntree %v\nnot as expected %v", testTree, expected)
	}

}
