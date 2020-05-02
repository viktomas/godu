package files

import (
	"sort"
)

type bySizeDesc []*File

func (f bySizeDesc) Len() int           { return len(f) }
func (f bySizeDesc) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f bySizeDesc) Less(i, j int) bool { return f[i].Size > f[j].Size }

// SortDesc sorts folder content by size from largest to smallest
func SortDesc(folder *File) {
	sort.Sort(bySizeDesc(folder.Files))
	for _, file := range folder.Files {
		SortDesc(file)
	}
}

func PruneSmallFiles(folder *File, limit int64) {
	prunedFiles := []*File{}
	for _, file := range folder.Files {
		if file.Size >= limit {
			PruneSmallFiles(file, limit)
			prunedFiles = append(prunedFiles, file)
		}
	}
	folder.Files = prunedFiles

}
