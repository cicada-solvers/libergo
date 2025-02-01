package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "init" {
		fmt.Println("Initializing database.")
		_, err := initDatabase()
		if err != nil {
			fmt.Printf("Error initializing database: %v\n", err)
			return
		}

		fmt.Println("Database initialized successfully.")
		return
	}

	var start, end int
	fmt.Print("Enter the start length: ")
	_, err := fmt.Scan(&start)
	if err != nil {
		fmt.Printf("Invalid input: %v\n", err)
		return
	}

	fmt.Print("Enter the end length: ")
	_, err = fmt.Scan(&end)
	if err != nil {
		fmt.Printf("Invalid input: %v\n", err)
		return
	}

	if start > end {
		fmt.Println("Start length cannot be greater than the end.")
		return
	}

	for length := start; length <= end; length++ {
		fmt.Printf("Processing: %d\n", length)
		calculatePermutationRanges(length)
	}
}
