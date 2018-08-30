package core

import (
	"sort"
)

type bySizeDesc []*File

func (f bySizeDesc) Len() int           { return len(f) }
func (f bySizeDesc) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f bySizeDesc) Less(i, j int) bool { return f[i].Size > f[j].Size }

func SortDesc(folder *File) {
	sort.Sort(bySizeDesc(folder.Files))
	for _, file := range folder.Files {
		SortDesc(file)
	}
}

func pruneFolder(folder *File, limit int64) {
	prunedFiles := []*File{}
	for _, file := range folder.Files {
		if file.Size >= limit {
			pruneFolder(file, limit)
			prunedFiles = append(prunedFiles, file)
		}
	}
	folder.Files = prunedFiles

}
