package main

import (
	"flag"
	"fmt"
	"ioc"
	"os"
	"strings"
)

// main is the entry point of the program that calculates the Index of Coincidence (IoC) for a given text input.
// It parses command-line flags to determine the input text and alphabet type, validates the input, and outputs the IoC value.
func main() {
	// Define command line flags
	alphabetType := flag.String("alphabet", "english", "Type of alphabet to use (english, runeglish, or rune)")
	text := flag.String("text", "", "Text to analyze")

	// Parse the flags
	flag.Parse()

	// Validate that text was provided
	if *text == "" {
		fmt.Println("Error: No text provided for analysis")
		flag.Usage()
		os.Exit(1)
	}

	// Get the appropriate alphabet based on the flag
	var alphabet ioc.AlphabetType
	switch strings.ToLower(*alphabetType) {
	case "english":
		alphabet = ioc.Latin
		break
	case "runeglish":
		alphabet = ioc.Runeglish
		break
	case "rune":
		alphabet = ioc.Rune
		break
	default:
		fmt.Println("Error: Invalid alphabet type")
		flag.Usage()
		os.Exit(1)
	}

	// Calculate the Index of Coincidence
	result := ioc.CalcIOC(*text, alphabet)

	// Output the result
	fmt.Printf("IoC: %f\n", result)
}
