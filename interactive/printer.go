package interactive

import (
	"fmt"
	"io"
	"strings"

	"github.com/viktomas/godu/core"
)

func PrintMarkedFiles(state *core.State, writer io.Writer) {
	markedFolders := make(map[string]*core.File)
	for file := range state.MarkedFiles {
		if file.IsDir {
			dirName := file.Path()
			// Mark folder with a slash after its name
			if _, found := markedFolders[dirName+"/"]; !found {
				markedFolders[dirName+"/"] = file
			}
		}
	}
markCheck:
	for file := range state.MarkedFiles {
		if file.Parent == nil {
			// This should never happen, just to make sure that we never nil-dereference
			continue
		}
		dirName := file.Parent.Name
		if dirName == "." {
			// Do not ignore folders in root dir
			dirName = file.Name
		} else {
			for folderName, folder := range markedFolders {
				// Check if the file is in any of the marked folders
				if strings.HasPrefix(dirName+"/", folderName) && file != folder {
					// If so, break
					continue markCheck
				}
			}
		}
		fmt.Fprintf(writer, "'%s'\n", file.Path())
	}
}
