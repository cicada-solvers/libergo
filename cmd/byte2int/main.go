package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
)

// main is the entry point of the program that reads comma-separated byte values, processes them, and outputs a big integer.
func main() {
	// Prompt the user for comma-separated values
	fmt.Print("Enter comma-separated values (0-255): ")
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
		value, conversionErr := strconv.Atoi(valueStr)
		if conversionErr != nil {
			fmt.Printf("Error: '%s' is not a valid integer\n", valueStr)
			return
		}

		if value < 0 || value > 255 {
			fmt.Printf("Error: Value %d is out of range (0-255)\n", value)
			return
		}

		bytes = append(bytes, byte(value))
	}

	// Convert bytes to big integer
	bigIntValue := new(big.Int).SetBytes(bytes)

	// Output to console
	fmt.Println("Big Integer value:", bigIntValue.String())

	// Write to output.txt file
	file, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			fmt.Println("Error closing output file:", closeErr)
		}
	}(file)

	_, err = file.WriteString(bigIntValue.String())
	if err != nil {
		fmt.Println("Error writing to output file:", err)
		return
	}

	fmt.Println("Result successfully written to output.txt")
}
