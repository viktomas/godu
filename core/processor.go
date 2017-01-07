package core

import "sync"

func StartProcessing(
	folder *File,
	limit int64,
	commands chan Executer,
	states chan State,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	PruneTree(folder, limit)
	SortDesc(folder)
	state := State{
		Folder: folder,
	}
	states <- state
	for {
		command, more := <-commands
		if !more {
			break
		}
		newState, err := command.Execute(state)
		if err == nil {
			state = newState
			states <- state
		}
	}
}
