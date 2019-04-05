package interactive

import (
	"sort"

	"github.com/viktomas/godu/core"
)

type byLength []string

func (l byLength) Len() int           { return len(l) }
func (l byLength) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l byLength) Less(i, j int) bool { return len(l[i]) > len(l[j]) }

// FilesAsSlice takes files from the map and returns a sorted slice of file paths.
func FilesAsSlice(in map[*core.File]struct{}) []string {
	out := make([]string, 0, len(in))
	for file := range in {
		p := file.Path()
		out = append(out, p)
	}
	// sorting length of the path (assuming that we want to delete files in subdirs first)
	// alphabetical sorting added for determinism (map keys doesn't guarantee order)
	sort.Sort(sort.StringSlice(out))
	sort.Sort(byLength(out))
	return out
}
