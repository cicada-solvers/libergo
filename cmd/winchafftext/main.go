package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"
	"sequences"
	"strings"
)

// main is the entry point for the program
func main() {
	// Define command-line flags
	textPtr := flag.String("text", "", "The text to process")
	keepThrowPtr := flag.String("filter", "keep", "Whether to 'keep' or 'throw' characters specified in the sequence")
	sequencePtr := flag.String("sequence", "", "The sequence of characters to keep or throw")
	outputFilePtr := flag.String("output", "", "Output file path (if not specified, prints to stdout)")

	// Parse the flags
	flag.Parse()

	// Validate input
	if *textPtr == "" {
		fmt.Println("Error: Text is required")
		flag.Usage()
		os.Exit(1)
	}

	if *sequencePtr == "" {
		fmt.Println("Error: Sequence is required")
		flag.Usage()
		os.Exit(1)
	}

	if *keepThrowPtr != "keep" && *keepThrowPtr != "throw" {
		fmt.Println("Error: Filter must be either 'keep' or 'throw'")
		flag.Usage()
		os.Exit(1)
	}

	// Getting the sequence
	sequence := getSequence(*sequencePtr, int64(len(*textPtr)))

	// Process the text
	result := processText(*textPtr, *keepThrowPtr, sequence)

	// Output result
	if *outputFilePtr == "" {
		// Print to stdout
		fmt.Println(result)
	} else {
		// Write to file
		err := os.WriteFile(*outputFilePtr, []byte(result), 0644)
		if err != nil {
			_, writeErr := fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
			if writeErr != nil {
				return
			}
			os.Exit(1)
		}
		fmt.Printf("Output written to %s\n", *outputFilePtr)
	}
}

// processText processes the text based on the filter type and sequence
// If the filter type is "keep", it keeps only the characters in the sequence
// If the filter type is "throw", it throws away the characters in the sequence
// The sequence is a sequence of numbers, where each number represents a character in the text
// For example, if the sequence is [1, 3, 5], then the text is processed as follows:
// 1. The first character is kept
// 2. The second character is thrown away
// 3. The third character is kept
func processText(text, filterType string, sequence *sequences.NumericSequence) string {
	var result strings.Builder

	if filterType == "keep" {
		for counter, char := range text {
			if isNumberInSequence(int64(counter), sequence) {
				result.WriteString(string(char))
			}
		}
	} else {
		for counter, char := range text {
			if !isNumberInSequence(int64(counter), sequence) {
				result.WriteString(string(char))
			}
		}
	}

	return result.String()
}

// getSequence returns a sequence based on the sequence type
func getSequence(sequenceStr string, maxNumber int64) *sequences.NumericSequence {
	var sequence *sequences.NumericSequence

	switch sequenceStr {
	case "central_polygonal":
		sequence, _ = sequences.GetCentralPolygonalNumbersSequence(big.NewInt(maxNumber), false)
	case "cubes":
		sequence, _ = sequences.GetCubesSequence(big.NewInt(maxNumber), false)
	case "natural":
		sequence, _ = sequences.GetNaturalSequence(big.NewInt(maxNumber), false)
	case "prime":
		sequence, _ = sequences.GetPrimeSequence(big.NewInt(maxNumber), false)
	case "fibonacci_prime":
		sequence, _ = sequences.GetFibonacciPrimeSequence(big.NewInt(maxNumber), false)
	case "cake":
		sequence, _ = sequences.GetCakeSequence(big.NewInt(maxNumber), false)
	case "catalan":
		sequence, _ = sequences.GetCatalanSequence(big.NewInt(maxNumber), false)
	case "totient":
		sequence, _ = sequences.GetTotientSequence(big.NewInt(maxNumber))
	case "totient_prime":
		sequence, _ = sequences.GetTotientPrimeSequence(big.NewInt(maxNumber))
	case "fibonacci":
		sequence, _ = sequences.GetFibonacciSequence(big.NewInt(maxNumber))
	case "zekendorf":
		sequence, _ = sequences.GetZekendorfRepresentationSequence(big.NewInt(maxNumber), false)
	case "lucas":
		sequence, _ = sequences.GenerateLucas(big.NewInt(maxNumber), false)
	case "collatz":
		sequence, _ = sequences.GetCollatzSequence(maxNumber, false)
	default:
		fmt.Printf("Unknown sequence type: %s\n", sequenceStr)
		os.Exit(1)
	}

	return sequence
}

// isNumberInSequence returns true if the number is in the sequence, false otherwise
func isNumberInSequence(num int64, sequence *sequences.NumericSequence) bool {
	for _, n := range sequence.Sequence {
		if n.Int64() == num {
			return true
		}
	}
	return false
}
