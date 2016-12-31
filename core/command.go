package core

// State represents system configuration after processing user input
type State struct {
	parent   *File
	folder   *File
	history  map[*File]int
	selected int
	filepath string
}

// Command represents a user action triggered on a State
type Command func(State) State
