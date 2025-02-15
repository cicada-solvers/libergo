package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"titler"
)

func main() {
	titler.PrintTitle("Binary Inverter")
	// Define command-line flags
	inputFile := flag.String("inputfile", "", "Input file path")
	outputFile := flag.String("outputfile", "", "Output file path")

	// Parse command-line flags
	flag.Parse()

	// Check if input and output files are provided
	if *inputFile == "" || *outputFile == "" {
		flag.Usage()
		log.Fatalf("Both inputfile and outputfile must be specified")
	}

	// Read the input file
	inputData, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	// Convert the byte array to a binary string
	binaryString := byteArrayToBinaryString(inputData)

	// Invert the binary string
	invertedBinaryString := invertBinaryString(binaryString)

	// Write the inverted binary string to the output file
	err = os.WriteFile(*outputFile, []byte(invertedBinaryString), 0644)
	if err != nil {
		log.Fatalf("Error writing to output file: %v", err)
	}

	fmt.Println("Inversion complete. Output written to", *outputFile)
}

// byteArrayToBinaryString converts a byte array to a binary string
func byteArrayToBinaryString(data []byte) string {
	var binaryString string
	for _, b := range data {
		binaryString += fmt.Sprintf("%08b", b)
	}
	return binaryString
}

// invertBinaryString inverts the 1s and 0s in a binary string
func invertBinaryString(binaryString string) string {
	var invertedString string
	for _, bit := range binaryString {
		if bit == '0' {
			invertedString += "1"
		} else {
			invertedString += "0"
		}
	}
	return invertedString
}
