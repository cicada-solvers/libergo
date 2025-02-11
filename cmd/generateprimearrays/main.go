package main

import (
	"fmt"
	"titler"
)

func main() {
	titler.PrintTitle("Generate Prime Arrays")
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

	fmt.Printf("Done generating prime arrays from %d to %d\n", start, end)
}
