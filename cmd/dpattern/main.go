package main

import (
	"flag"
	"fmt"
	"liberdatabase"
)

func main() {
	// Define flags
	textFlag := flag.String("text", "", "Text to get rune pattern")

	// Parse flags
	flag.Parse()

	// Check if no flags are provided
	if len(flag.Args()) == 0 {
		flag.Usage()
		return
	}

	// Create a DictionaryWord instance
	dw := liberdatabase.DictionaryWord{
		RuneWordText: *textFlag,
	}

	// Get the rune pattern
	pattern := dw.GetRunePattern()

	// Output the result
	fmt.Println("Rune Pattern:", pattern)
}
