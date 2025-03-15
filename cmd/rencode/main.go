package main

import (
	"flag"
	"fmt"
	"os"
	"runer"
	"titler"
)

// main reads input text, encodes it, and writes the result to an output file or stdout.
func main() {
	titler.PrintTitle("Gematria Encoder / Decoder")

	// Define flags
	textFlag := flag.String("text", "", "Text to be encoded")
	fileFlag := flag.String("file", "", "File containing text to be encoded (overrides text flag)")
	typeFlag := flag.String("type", "l2r", "Type of encoding: 'l2r' or 'r2l'")
	outputFile := flag.String("output", "", "Output file to write the encoded text")
	helpFlag := flag.Bool("help", false, "Display help")

	// Parse flags
	flag.Parse()

	// Display help if requested
	if *helpFlag {
		flag.Usage()
		return
	}

	// Read input text
	var inputText string
	if *fileFlag != "" {
		data, err := os.ReadFile(*fileFlag)
		if err != nil {
			fmt.Println("Error reading file:", err)
			os.Exit(1)
		}
		inputText = string(data)
	} else {
		inputText = *textFlag
	}

	// Perform encoding
	var encodedText string
	switch *typeFlag {
	case "l2r":
		prepped := runer.PrepLatinToRune(inputText)
		encodedText = runer.TransposeLatinToRune(prepped)
	case "r2l":
		encodedText = runer.TransposeRuneToLatin(inputText)
	default:
		fmt.Println("Invalid encoding type:", *typeFlag)
		os.Exit(1)
	}

	// Output result
	if *outputFile != "" {
		err := os.WriteFile(*outputFile, []byte(encodedText), 0644)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			os.Exit(1)
		}
	} else {
		fmt.Println(encodedText)
	}
}
