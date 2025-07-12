package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

// main is the entry point of the application, prompting the user to convert a hex string into comma-separated decimals.
func main() {
	// Prompt the user for the hex string
	fmt.Print("Enter a hex string: ")

	// Read the user input
	reader := bufio.NewReader(os.Stdin)
	hexString, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// Clean the input (remove whitespace, newlines)
	hexString = strings.TrimSpace(hexString)

	// Convert hex string to bytes
	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		fmt.Println("Error decoding hex string:", err)
		return
	}

	// Convert bytes to base 10 values and create comma separated string
	var decimalValues []string
	for _, b := range bytes {
		decimalValues = append(decimalValues, fmt.Sprintf("%d", b))
	}

	// Output the comma-separated decimal values
	fmt.Println(strings.Join(decimalValues, ","))
}
