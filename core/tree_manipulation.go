package core

import (
	"sort"
)

type bySizeDesc []*File

func (f bySizeDesc) Len() int           { return len(f) }
func (f bySizeDesc) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f bySizeDesc) Less(i, j int) bool { return f[i].Size > f[j].Size }

func SortDesc(tree *File) {
	sort.Sort(bySizeDesc(tree.Files))
	for _, file := range tree.Files {
		SortDesc(file)
	}
}

func PruneTree(tree *File, limit int64) {
	prunedFiles := []*File{}
	for _, file := range tree.Files {
		if file.Size >= limit {
			PruneTree(file, limit)
			prunedFiles = append(prunedFiles, file)
		}
	}
	tree.Files = prunedFiles

}
