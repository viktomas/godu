package godu

// Command represents a user action triggered on a State
type Command func(State) State
