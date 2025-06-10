package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Prompt the user for input
	fmt.Println("Enter a byte array as comma-separated values (e.g., 10,255,7,0):")
	reader := bufio.NewReader(os.Stdin)
	input, inputError := reader.ReadString('\n')
	if inputError != nil {
		fmt.Println("Error reading input:", inputError)
		return
	}

	// Trim newline and any spaces
	input = strings.TrimSpace(input)

	// Remove all spaces from the input
	input = strings.ReplaceAll(input, " ", "")

	// Split the input by comma
	byteStrings := strings.Split(input, ",")

	// Convert to bytes and then to binary strings
	var binaryStrings strings.Builder

	for _, byteStr := range byteStrings {
		if byteStr == "" {
			continue
		}

		// Parse the byte value
		byteVal, parseError := strconv.ParseUint(byteStr, 10, 8)
		if parseError != nil {
			fmt.Printf("Error parsing value '%s': %v\n", byteStr, parseError)
			continue
		}

		// Convert to 8-bit binary representation
		binaryStr := fmt.Sprintf("%08b", byteVal)
		binaryStrings.WriteString(binaryStr)
	}

	// Print to console
	fmt.Println("\nBinary representation:")
	fmt.Println(binaryStrings.String())

	// Write to file (append mode)
	file, openError := os.OpenFile("byte2binstring_output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if openError != nil {
		fmt.Println("Error opening output file:", openError)
		return
	}
	defer func(file *os.File) {
		closeError := file.Close()
		if closeError != nil {
			fmt.Println("Error closing output file:", closeError)
		}
	}(file)

	// Write output to the file with a timestamp
	timestamp := fmt.Sprintf("%s\n", time.Now().Format("2006-01-02 15:04:05"))
	outputString := fmt.Sprintf("%s%s\n\n", timestamp, binaryStrings.String())
	if _, writeError := file.WriteString(outputString); writeError != nil {
		fmt.Println("Error writing to output file:", writeError)
		return
	}

	fmt.Println("\nOutput has been appended to byte2binstring_output.txt")
}
