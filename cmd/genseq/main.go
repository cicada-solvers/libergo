package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"
	"sequences"
	"strings"
	"titler"
)

// main is the entry point for the program
func main() {
	titler.PrintTitle("Sequence Generator")

	// Define the named parameters
	maxNumberString := flag.String("max", "100", "The maximum number")
	sequenceType := flag.String("type", "default", "The type of sequence")
	positional := flag.Bool("positional", false, "Whether the sequence is positional")
	help := flag.Bool("help", false, "List supported sequences and describe flags")
	output := flag.String("output", "", "Output file to write the sequence as a comma-separated string")

	// Parse the command-line flags
	flag.Parse()

	// If help flag is set, print the supported sequences, flag descriptions, and examples, then exit
	if *help {
		fmt.Println("Supported sequences:")
		fmt.Println(" - central_polygonal")
		fmt.Println(" - cubes")
		fmt.Println(" - natural")
		fmt.Println(" - prime")
		fmt.Println(" - fibonacci_prime")
		fmt.Println(" - cake")
		fmt.Println(" - catalan")
		fmt.Println(" - totient")
		fmt.Println(" - totient_prime")
		fmt.Println(" - fibonacci")
		fmt.Println(" - zekendorf")
		fmt.Println(" - lucas")
		fmt.Println(" - collatz")
		fmt.Println("\nFlags:")
		fmt.Println(" -max: The maximum number (default: 100)")
		fmt.Println(" -type: The type of sequence (default: default)")
		fmt.Println(" -positional: Whether the sequence is positional (default: false)")
		fmt.Println(" -help: List supported sequences and describe flags")
		fmt.Println(" -output: Output file to write the sequence as a comma-separated string")
		fmt.Println("\nExamples:")
		fmt.Println(" ./genseq")
		fmt.Println(" ./genseq -max=200")
		fmt.Println(" ./genseq -type=fibonacci")
		fmt.Println(" ./genseq -positional=true")
		fmt.Println(" ./genseq -max=200 -type=fibonacci -positional=true")
		fmt.Println(" ./genseq -output=sequence.txt")
		fmt.Println(" ./genseq -help")
		os.Exit(0)
	}

	// Print the parameters to the console
	fmt.Printf("Max Number: %s\n", *maxNumberString)
	fmt.Printf("Sequence Type: %s\n", *sequenceType)
	fmt.Printf("Positional: %t\n", *positional)

	maxNumber := new(big.Int)
	maxNumber, ok := maxNumber.SetString(*maxNumberString, 10)
	if !ok {
		fmt.Printf("Invalid max number: %s\n", *maxNumberString)
		os.Exit(1)
	}

	// Generate and print the sequence based on the sequence type
	var sequence *sequences.NumericSequence
	var err error

	switch *sequenceType {
	case "central_polygonal":
		sequence, err = sequences.GetCentralPolygonalNumbersSequence(maxNumber, *positional)
	case "cubes":
		sequence, err = sequences.GetCubesSequence(maxNumber, *positional)
	case "natural":
		sequence, err = sequences.GetNaturalSequence(maxNumber, *positional)
	case "prime":
		sequence, err = sequences.GetPrimeSequence(maxNumber, *positional)
	case "fibonacci_prime":
		sequence, err = sequences.GetFibonacciPrimeSequence(maxNumber, *positional)
	case "cake":
		sequence, err = sequences.GetCakeSequence(maxNumber, *positional)
	case "catalan":
		sequence, err = sequences.GetCatalanSequence(maxNumber, *positional)
	case "totient":
		sequence, err = sequences.GetTotientSequence(maxNumber)
	case "totient_prime":
		sequence, err = sequences.GetTotientPrimeSequence(maxNumber)
	case "fibonacci":
		sequence, err = sequences.GetFibonacciSequence(maxNumber)
	case "zekendorf":
		sequence, err = sequences.GetZekendorfRepresentationSequence(maxNumber, *positional)
	case "lucas":
		sequence, err = sequences.GenerateLucas(maxNumber, *positional)
	case "collatz":
		sequence, err = sequences.GetCollatzSequence(maxNumber.Int64(), *positional)
	default:
		fmt.Printf("Unknown sequence type: %s\n", *sequenceType)
		os.Exit(1)
	}

	if err != nil {
		fmt.Println("Error generating sequence:", err)
		os.Exit(1)
	}

	var sequenceStrings []string
	for _, num := range sequence.Sequence {
		sequenceStrings = append(sequenceStrings, num.String())
	}

	// Print the sequence to the console
	fmt.Printf("Sequence Count: %d - IsPrime: %v\n", len(sequence.Sequence), sequences.IsPrime(new(big.Int).SetInt64(int64(len(sequence.Sequence)))))
	fmt.Printf("Sequence: %s\n", strings.Join(sequenceStrings, ","))

	// If output flag is set, write the sequence to the specified file
	if *output != "" {
		file, err := os.Create(*output)
		if err != nil {
			fmt.Println("Error creating file:", err)
			os.Exit(1)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				fmt.Println("Error closing file:", err)
				os.Exit(1)
			}
		}(file)

		_, err = file.WriteString(strings.Join(sequenceStrings, ","))
		if err != nil {
			fmt.Println("Error writing to file:", err)
			os.Exit(1)
		}
		fmt.Printf("Sequence written to file: %s\n", *output)
	}
}
