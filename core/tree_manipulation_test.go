package core

import (
	"reflect"
	"testing"
)

func TestSortTree(t *testing.T) {
	testTree := &File{"b", nil, 260, true, []*File{
		&File{"c", nil, 100, false, []*File{}},
		&File{"d", nil, 160, true, []*File{
			&File{"e", nil, 50, false, []*File{}},
			&File{"f", nil, 30, false, []*File{}},
			&File{"g", nil, 80, true, []*File{
				&File{"i", nil, 30, false, []*File{}},
				&File{"j", nil, 50, false, []*File{}},
			}},
		}},
	}}
	expected := &File{"b", nil, 260, true, []*File{
		&File{"d", nil, 160, true, []*File{
			&File{"g", nil, 80, true, []*File{
				&File{"j", nil, 50, false, []*File{}},
				&File{"i", nil, 30, false, []*File{}},
			}},
			&File{"e", nil, 50, false, []*File{}},
			&File{"f", nil, 30, false, []*File{}},
		}},
		&File{"c", nil, 100, false, []*File{}},
	}}
	SortDesc(testTree)
	if !reflect.DeepEqual(testTree, expected) {
		t.Errorf("tree not sorted correctly\ntree %v\nnot as expected %v", testTree, expected)
	}
}

func TestPruneTree(t *testing.T) {
	testTree := &File{"b", nil, 260, true, []*File{
		&File{"c", nil, 100, false, []*File{}},
		&File{"d", nil, 160, true, []*File{
			&File{"e", nil, 50, false, []*File{}},
			&File{"f", nil, 30, false, []*File{}},
			&File{"g", nil, 80, true, []*File{
				&File{"i", nil, 50, false, []*File{}},
				&File{"j", nil, 30, false, []*File{}},
			}},
		}},
	}}
	expected := &File{"b", nil, 260, true, []*File{
		&File{"c", nil, 100, false, []*File{}},
		&File{"d", nil, 160, true, []*File{
			&File{"g", nil, 80, true, []*File{}},
		}},
	}}
	PruneTree(testTree, 60)
	if !reflect.DeepEqual(testTree, expected) {
		t.Errorf("tree not pruned correctly\ntree %v\nnot as expected %v", testTree, expected)
	}

}
