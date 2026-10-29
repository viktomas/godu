package interactive

import (
	"sort"

	"github.com/viktomas/godu/files"
)

// FilesAsSlice takes files from the map and returns a sorted slice of file paths.
func FilesAsSlice(in map[*files.File]struct{}) []string {
	out := make([]string, 0, len(in))
	for file := range in {
		p := file.Path()
		out = append(out, p)
	}
	// sorting by length of the path (assuming that we want to delete files in subdirs first)
	// with alphabetical order as a tie-breaker — without it, equal-length paths
	// could be reordered, since sort.Slice is not stable
	sort.Slice(out, func(i, j int) bool {
		if len(out[i]) != len(out[j]) {
			return len(out[i]) > len(out[j])
		}
		return out[i] < out[j]
	})
	return out
}
