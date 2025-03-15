package main

import (
	"flag"
	"fmt"
	"lgstructs"
)

func main() {
	// Define flags
	textFlag := flag.String("text", "", "Text to get rune pattern")

	// Parse flags
	flag.Parse()

	// Check if the text flag is empty
	if *textFlag == "" {
		flag.Usage()
		return
	}

	// Create a DictionaryWord instance
	dw := lgstructs.DictionaryWord{
		RuneWordText: *textFlag,
	}

	// Get the rune pattern
	pattern := lgstructs.GetRunePattern(dw)

	// Output the result
	fmt.Println("Pattern:", pattern)
}
