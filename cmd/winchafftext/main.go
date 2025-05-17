package main

import (
	"flag"
	"fmt"
	"ioc"
	"math/big"
	"os"
	"runer"
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
	customSequencePtr := flag.String("customSequence", "", "Custom sequence of characters to keep or throw")
	startAtOnePtr := flag.Bool("startAtOne", false, "Whether to start the sequence at 1 or 0")
	pullFromSequence := flag.Bool("pullFromSequence", false, "Whether to pull the sequence from the sequence file")

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
	var sequence *sequences.NumericSequence
	if *customSequencePtr != "" {
		sequence = &sequences.NumericSequence{}
		sequencesArray := strings.Split(*customSequencePtr, ",")
		for _, s := range sequencesArray {
			// Convert the string to a big.Int
			newNum, _ := big.NewInt(0).SetString(strings.TrimSpace(s), 10)
			// Add the number to the sequence
			sequence.Sequence = append(sequence.Sequence, newNum)
		}
		printSequence(sequence)
	} else {
		sequence = getSequence(*sequencePtr, int64(len(*textPtr)))
		printSequence(sequence)
	}

	// Process the text
	result, filteredResult := processText(*textPtr, *keepThrowPtr, *startAtOnePtr, *pullFromSequence, sequence)

	// Output result
	if *outputFilePtr == "" {
		// Print to stdout
		fmt.Println(result)
	} else {
		// Write to the file
		if len([]byte(filteredResult)) > 0 {
			err := os.WriteFile(*outputFilePtr, []byte(result), 0644)
			if err != nil {
				_, writeErr := fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
				if writeErr != nil {
					return
				}
				os.Exit(1)
			}

			err = os.WriteFile(fmt.Sprintf("%s.filtered.txt", *outputFilePtr), []byte(filteredResult), 0644)
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
}

// processText processes the text based on the filter type and sequence
// If the filter type is "keep", it keeps only the characters in the sequence
// If the filter type is "throw", it throws away the characters in the sequence.
// The sequence is a sequence of numbers, where each number represents a character in the text.
// For example, if the sequence is [1, 3, 5], then the text is processed as follows:
// 1. The first character is kept
// 2. The second character is thrown away
// 3. The third character is kept
func processText(text, filterType string, startAtOne, pullFromSequence bool, sequence *sequences.NumericSequence) (string, string) {
	var result strings.Builder
	var filteredResult strings.Builder
	stringArray := strings.Split(text, "")

	result.WriteString(fmt.Sprintf("Original Text:\n %s\n\n", text))

	result.WriteString("Sequence:\n")
	result.WriteString(fmt.Sprintf(" %v\n\n", sequence.Sequence))

	result.WriteString(fmt.Sprintf("Filtered Type: %s\n\n", filterType))

	result.WriteString(fmt.Sprintf("Start At One: %t\n\n", startAtOne))

	result.WriteString(fmt.Sprintf("Pull From Sequence: %t\n\n", pullFromSequence))

	result.WriteString("Filtered Text:\n")

	if filterType == "keep" {
		if pullFromSequence {
			for _, value := range sequence.Sequence {
				if value.Int64() > int64(len(stringArray)) {
					continue
				}

				if startAtOne {
					stringVal := stringArray[value.Int64()-1]
					result.WriteString(stringVal)
					filteredResult.WriteString(stringVal)
				} else {
					stringVal := stringArray[value.Int64()]
					result.WriteString(stringVal)
					filteredResult.WriteString(stringVal)
				}
			}
		} else {
			for counter, char := range stringArray {
				if isNumberInSequence(getCheckNumber(counter, startAtOne), sequence) {
					result.WriteString(char)
					filteredResult.WriteString(char)
				}
			}
		}
	} else {
		for counter, char := range stringArray {
			if !isNumberInSequence(getCheckNumber(counter, startAtOne), sequence) {
				result.WriteString(char)
				filteredResult.WriteString(char)
			}
		}
	}

	// Get the IOC of the filteredResult
	iocValue := ioc.CalcIOC(filteredResult.String(), ioc.Rune)
	result.WriteString(fmt.Sprintf("\n\nIOC: %f\n", iocValue))

	clearText := runer.TransposeRuneToLatin(filteredResult.String())
	result.WriteString(fmt.Sprintf("\nClear Text:\n%s\n", clearText))

	// Get the IOC of the clearText
	iocValue = ioc.CalcIOC(clearText, ioc.Runeglish)
	result.WriteString(fmt.Sprintf("\nIOC: %f\n", iocValue))

	return result.String(), filteredResult.String()
}

// getCheckNumber gets the value based off whether to start at one
func getCheckNumber(num int, startAtOne bool) int64 {
	if startAtOne {
		return int64(num - 1)
	}
	return int64(num)
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

// printSequence is a method that prints the sequence for logging
func printSequence(sequence *sequences.NumericSequence) {
	for _, n := range sequence.Sequence {
		fmt.Printf("%d ", n.Int64())
	}
	fmt.Println()
}
