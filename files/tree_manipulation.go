package files

import (
	"sort"
)

// SortDesc sorts folder content by size from largest to smallest
func SortDesc(folder *File) {
	sort.Slice(folder.Files,
		func(i, j int) bool {
			return folder.Files[i].Size > folder.Files[j].Size
		})
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
