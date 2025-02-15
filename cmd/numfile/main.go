package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Define flags
	inputFlag := flag.String("input", "", "Comma-separated string of numbers (0-255)")
	outputFlag := flag.String("output", "output.bin", "Output file name")

	// Parse flags
	flag.Parse()

	// Check if input flag is provided
	if *inputFlag == "" {
		fmt.Println("Input string is required")
		os.Exit(1)
	}

	// Split the input string by commas
	numStrings := strings.Split(*inputFlag, ",")

	// Create a byte slice to hold the converted numbers
	var byteSlice []byte

	// Convert each number string to a byte and append to the byte slice
	for _, numStr := range numStrings {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			fmt.Println("Error converting string to number:", err)
			os.Exit(1)
		}
		if num < 0 || num > 255 {
			fmt.Println("Number out of byte range:", num)
			os.Exit(1)
		}
		byteSlice = append(byteSlice, byte(num))
	}

	// Write the byte slice to a file
	err := os.WriteFile(*outputFlag, byteSlice, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		os.Exit(1)
	}

	fmt.Println("Bytes written to", *outputFlag)
}
