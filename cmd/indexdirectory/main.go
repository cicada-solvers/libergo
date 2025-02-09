package main

import (
	"charindex"
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a directory path")
		os.Exit(1)
	}

	directory := flag.String("directory", "", "Directory to index characters from")
	flag.Parse()

	if *directory == "" {
		fmt.Println("Directory path is required")
		os.Exit(1)
	}

	err := charindex.IndexCharactersFromDirectory(*directory)
	if err != nil {
		fmt.Println("Error indexing characters from directory:", err)
		os.Exit(1)
	}

	fmt.Println("Character indexing completed successfully")
}
