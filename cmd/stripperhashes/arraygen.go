package main

import "liberdatabase"

// Program represents the program
type Program struct {
	tasks chan liberdatabase.ReadPermutation
}

// NewProgram creates a new Program
func NewProgram() *Program {
	return &Program{
		tasks: make(chan liberdatabase.ReadPermutation, 10000),
	}
}
