package commands

import (
	"fmt"
	"sync"

	"github.com/viktomas/godu/files"
)

// ProcessFolder removes small files and sorts folder content based on accumulated size
func ProcessFolder(folder *files.File, limit int64) error {
	files.PruneSmallFiles(folder, limit)
	if len(folder.Files) == 0 {
		return fmt.Errorf("the folder '%s' doesn't contain any files bigger than %dMB", folder.Name, limit/files.MEGABYTE)
	}
	files.SortDesc(folder)
	return nil
}

// StartProcessing reads user commands and applies them to state
func StartProcessing(
	folder *files.File,
	commands <-chan Executer,
	states chan<- State,
	lastStateChan chan<- *State,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	state := State{
		Folder:      folder,
		MarkedFiles: make(map[*files.File]struct{}),
	}
	states <- state
	for {
		command, more := <-commands
		if !more {
			close(states)
			break
		}
		if newState, err := command.Execute(state); err == nil {
			state = newState
			states <- state
		}
	}
	lastStateChan <- &state
}
