package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"strings"
)

func main() {
	// Prompt the user for a hex value
	fmt.Print("Enter a hexadecimal value: ")
	reader := bufio.NewReader(os.Stdin)
	hexInput, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// Trim whitespace and newlines from the input
	hexInput = strings.TrimSpace(hexInput)

	// Convert hex string to big integer
	bigIntValue := new(big.Int)
	_, success := bigIntValue.SetString(hexInput, 16)
	if !success {
		fmt.Println("Error: Invalid hexadecimal value")
		return
	}

	// Output to console
	fmt.Println("Decimal value:", bigIntValue.String())

	// Write to output.txt file
	file, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer func(file *os.File) {
		closeError := file.Close()
		if closeError != nil {
			fmt.Println("Error closing output file:", closeError)
		}
	}(file)

	_, err = file.WriteString(bigIntValue.String())
	if err != nil {
		fmt.Println("Error writing to output file:", err)
		return
	}

	fmt.Println("Result successfully written to output.txt")
}
