package interactive

import (
	"fmt"
	"sort"
	"strings"

	"github.com/viktomas/godu/core"
)

type byLength []string

func (l byLength) Len() int           { return len(l) }
func (l byLength) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l byLength) Less(i, j int) bool { return len(l[i]) > len(l[j]) }

// FilesAsSlice takes files from the map and returns a sorted slice of file paths.
// Each file path is quoted if 'quote' is set to true.
func FilesAsSlice(in map[*core.File]struct{}, quote bool) []string {
	out := make([]string, 0, len(in))
	for file := range in {
		p := file.Path()
		if quote {
			// Escape single quotes
			p = strings.Replace(p, "'", "\\'", -1)
			// Quote full file path
			p = fmt.Sprintf("'%s'", p)
		}
		out = append(out, p)
	}
	// sorting lenght of the path (assuming that we want to deleate files in subdirs first)
	// alphabetical sorting added for determinism (map keys doesn't guarantee order)
	sort.Sort(sort.StringSlice(out))
	sort.Sort(byLength(out))
	return out
}
