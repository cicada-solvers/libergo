package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Prompt the user for comma-separated byte values
	fmt.Print("Enter comma-separated byte values (0-255): ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// Trim whitespace and newlines from the input
	input = strings.TrimSpace(input)

	// Split the input by commas
	valueStrings := strings.Split(input, ",")

	// Create a byte slice to hold the values
	bytes := make([]byte, 0, len(valueStrings))

	// Parse each value and add it to the byte slice
	for _, valueStr := range valueStrings {
		valueStr = strings.TrimSpace(valueStr)
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			fmt.Printf("Error: '%s' is not a valid integer\n", valueStr)
			return
		}

		if value < 0 || value > 255 {
			fmt.Printf("Error: Value %d is out of range (0-255)\n", value)
			return
		}

		bytes = append(bytes, byte(value))
	}

	// Ask for output filename
	fmt.Print("Enter output filename (default: output.bin): ")
	filename, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading filename:", err)
		return
	}

	filename = strings.TrimSpace(filename)
	if filename == "" {
		filename = "output.bin"
	}

	// Write bytes to the output file
	err = os.WriteFile(filename, bytes, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Printf("Successfully wrote %d bytes to %s\n", len(bytes), filename)
}
