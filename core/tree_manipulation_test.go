package core

import (
	"reflect"
	"testing"
)

func TestSortTree(t *testing.T) {
	testTree := &File{"b", 260, []*File{
		&File{"c", 100, []*File{}},
		&File{"d", 160, []*File{
			&File{"e", 50, []*File{}},
			&File{"f", 30, []*File{}},
			&File{"g", 80, []*File{
				&File{"i", 30, []*File{}},
				&File{"j", 50, []*File{}},
			}},
		}},
	}}
	expected := &File{"b", 260, []*File{
		&File{"d", 160, []*File{
			&File{"g", 80, []*File{
				&File{"j", 50, []*File{}},
				&File{"i", 30, []*File{}},
			}},
			&File{"e", 50, []*File{}},
			&File{"f", 30, []*File{}},
		}},
		&File{"c", 100, []*File{}},
	}}
	SortDesc(testTree)
	if !reflect.DeepEqual(testTree, expected) {
		t.Errorf("tree not sorted correctly\ntree %v\nnot as expected %v", testTree, expected)
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
