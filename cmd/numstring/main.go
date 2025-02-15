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

	// Parse flags
	flag.Parse()

	// Check if input flag is provided
	if *inputFlag == "" {
		fmt.Println("Input string is required")
		os.Exit(1)
	}

	// Split the input string by commas
	numStrings := strings.Split(*inputFlag, ",")

	// Create a string slice to hold the converted numbers
	var strValues []string

	// Convert each number string to its string representation and append to the string slice
	for _, numStr := range numStrings {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			fmt.Println("Error converting string to number:", err)
			os.Exit(1)
		}
		if num < 0 || num > 255 {
			fmt.Println("Number out of range:", num)
			os.Exit(1)
		}
		strValues = append(strValues, string(byte(num)))
	}

	// Join the string slice into a single string
	outputString := strings.Join(strValues, "")
	outputString = outputString + "\n"
	fmt.Println(outputString)
}
