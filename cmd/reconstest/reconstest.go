package main

import (
	"flag"
	"fmt"
	"runer"
	"strings"
)

func main() {
	// Define a command-line flag to take the text input
	textFlag := flag.String("text", "", "Text to be processed")

	// Parse the command-line flags
	flag.Parse()

	// Check if the text flag was provided
	if *textFlag == "" {
		fmt.Println("Please provide text using the -text flag")
		return
	}

	// Split the input text into a rune array and store its length
	inputRuneArray := strings.Split(*textFlag, "")
	inputLength := len(inputRuneArray)
	fmt.Printf("Original text: %s\n", *textFlag)
	fmt.Printf("Original text length: %d\n", inputLength)

	// Convert runes to Latin using TransposeRuneToLatin
	latinText := runer.TransposeRuneToLatin(*textFlag)
	fmt.Printf("Latin text: %s\n", latinText)

	// Convert Latin back to runes using TransposeLatinToRune
	// Note: using false as the second parameter (encodeBackwards) based on the function signature
	runicTextAgain := runer.TransposeLatinToRune(latinText, false)
	fmt.Printf("Runic text after conversion: %s\n", runicTextAgain)

	// Split the resulting runic text into an array
	outputRuneArray := strings.Split(runicTextAgain, "")
	outputLength := len(outputRuneArray)
	fmt.Printf("Output text length: %d\n", outputLength)

	// Compare lengths of the rune arrays before and after transposition
	if inputLength == outputLength {
		fmt.Println("The lengths of the rune arrays before and after transposition are the same.")
	} else {
		fmt.Println("The lengths of the rune arrays before and after transposition are different.")
		fmt.Printf("Original length: %d, Final length: %d\n", inputLength, outputLength)
	}
}
