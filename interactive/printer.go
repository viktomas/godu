package interactive

import (
	"fmt"
	"sort"

	"github.com/viktomas/godu/core"
)

type byLength []string

func (l byLength) Len() int           { return len(l) }
func (l byLength) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l byLength) Less(i, j int) bool { return len(l[i]) > len(l[j]) }

// QuoteMarkedFiles takes files from the map and returns slice of qoted file paths
func QuoteMarkedFiles(markedFiles map[*core.File]struct{}) []string {
	quotedFiles := make([]string, len(markedFiles))
	i := 0
	for file := range markedFiles {
		quotedFiles[i] = fmt.Sprintf("'%s'", file.Path())
		i++
	}
	// sorting lenght of the path (assuming that we want to deleate files in subdirs first)
	// alphabetical sorting added for determinism (map keys doesn't guarantee order)
	sort.Sort(sort.StringSlice(quotedFiles))
	sort.Sort(byLength(quotedFiles))
	return quotedFiles
}
