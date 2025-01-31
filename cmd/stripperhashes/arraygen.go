package main

// Program represents the program
type Program struct {
	tasks chan []byte
}

// NewProgram creates a new Program
func NewProgram() *Program {
	return &Program{
		tasks: make(chan []byte, 10000),
	}
}
