package main

import (
	"fmt"
)

func main() {
	var start, end int
	fmt.Print("Enter the start of the range: ")
	_, err := fmt.Scan(&start)
	if err != nil {
		fmt.Printf("Invalid input: %v\n", err)
		return
	}

	fmt.Print("Enter the end of the range: ")
	_, err = fmt.Scan(&end)
	if err != nil {
		fmt.Printf("Invalid input: %v\n", err)
		return
	}

	if start > end {
		fmt.Println("Start of the range cannot be greater than the end.")
		return
	}

	for length := start; length <= end; length++ {
		fmt.Printf("Processing: %d\n", length)
		calculatePermutationRanges(length)
	}
}
