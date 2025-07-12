package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// main is the entry point of the program that collects user input, generates permutations, and writes them to a file.
func main() {
	// remove output files
	removeOutputFileIfExists()

	// Create a scanner to read user input
	scanner := bufio.NewScanner(os.Stdin)

	// Prompt for the number of strings
	fmt.Print("Enter the number of strings you want to permute: ")
	scanner.Scan()
	numStr := scanner.Text()

	// Convert the input to an integer
	n, err := strconv.Atoi(numStr)
	if err != nil || n <= 0 {
		fmt.Println("Please enter a valid positive number")
		return
	}

	// Collect the strings from the user
	stringsArray := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Printf("Enter string #%d: ", i+1)
		scanner.Scan()
		stringsArray[i] = scanner.Text()
	}

	// Generate and print all permutations
	fmt.Println("\nAll Combinations:")
	permutations := generatePermutations(stringsArray)
	for i, perm := range permutations {
		output := fmt.Sprintf("Combination %d: %s\n\n\n", i+1, strings.Join(perm, ","))
		fmt.Print(output)
		_ = writePermutationsToFile(output, "output.txt")
	}
}

// generatePermutations returns all possible permutations of the input strings
func generatePermutations(input []string) [][]string {
	var result [][]string

	// Helper function to generate permutations recursively
	var generate func([]string, int)
	generate = func(arr []string, n int) {
		if n == 1 {
			// Create a copy of the current permutation
			tmp := make([]string, len(arr))
			copy(tmp, arr)
			result = append(result, tmp)
			return
		}

		for i := 0; i < n; i++ {
			// Generate permutations with current n-1 elements
			generate(arr, n-1)

			// Swap elements based on whether n is odd or even
			if n%2 == 1 {
				arr[0], arr[n-1] = arr[n-1], arr[0]
			} else {
				arr[i], arr[n-1] = arr[n-1], arr[i]
			}
		}
	}

	// Start the recursive permutation generation
	generate(input, len(input))
	return result
}

// removeOutputFileIfExists checks if "output.txt" exists and removes it if it does.
func removeOutputFileIfExists() {
	if _, err := os.Stat("output.txt"); !os.IsNotExist(err) {
		_ = os.Remove("output.txt")
	}
}

// writePermutationsToFile writes all permutations to the specified output file
func writePermutationsToFile(permutation string, filename string) error {
	// Create or open the file for writing
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// Create a buffered writer for better performance
	writer := bufio.NewWriter(file)
	defer func(writer *bufio.Writer) {
		_ = writer.Flush()
	}(writer)

	// Write each permutation to the file
	_, writeError := writer.WriteString(permutation)
	if writeError != nil {
		return writeError
	}

	return nil
}
