package main

import (
	runelib "characterrepo"
	"cipher"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// Define the flags
	text := flag.String("text", "", "The text to decode")
	alphabet := flag.String("alphabet", "rune", "The alphabet to use (rune or english)")
	outputFile := flag.String("output", "", "The output file to write the results")
	wordFile := flag.String("wordfile", "", "The CSV file of words to try for brute force decoding")
	ciphertype := flag.String("ciphertype", "caesar", "The cipher to use (vigenere, atbash, affine, autokey, caesar, trithemius)")

	// Parse the flags
	flag.Parse()

	// Validate required flags
	if *text == "" {
		log.Fatal("The -text flag is required")
	}
	if *outputFile == "" {
		log.Fatal("The -output flag is required")
	}

	// Open the output file
	file, err := os.Create(*outputFile)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer file.Close()

	// Print the parsed flags (for debugging or further processing)
	fmt.Printf("Text: %s\n", *text)
	fmt.Printf("Alphabet: %s\n", *alphabet)
	fmt.Printf("Output File: %s\n", *outputFile)
	fmt.Printf("Word File: %s\n", *wordFile)
	fmt.Printf("Cipher: %s\n", *ciphertype)

	// Add your decoding logic here
	// Determine the alphabet to use
	var alphabetSet []string
	var decodeToLatin bool
	var columnIndex int

	if strings.ToLower(*alphabet) == "rune" {
		repo := runelib.NewCharacterRepo()
		alphabetSet = repo.GetGematriaRunes()
		decodeToLatin = true
		columnIndex = 2
	} else {
		alphabetSet = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
			"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
		decodeToLatin = false
		columnIndex = 0
	}

	// Write the text to the output file
	_, err = file.WriteString(fmt.Sprintf("Text: %s\n", *text))

	// Now we are going to decode the text based on the cipher type
	var decodedText string
	var decodeErr error

	switch strings.ToLower(*ciphertype) {
	case "caesar":
		decodedText, decodeErr = cipher.BulkDecodeCaesarString(alphabetSet, strings.Split(*text, ""), decodeToLatin)
		if decodeErr != nil {
			fmt.Printf("Failed to decode using Caesar cipher: %v", decodeErr)
		}
	case "affine":
		decodedText, decodeErr = cipher.BulkDecodeAffineCipher(alphabetSet, *text, decodeToLatin)
		if decodeErr != nil {
			fmt.Printf("Failed to decode using Affine cipher: %v", decodeErr)
		}
	case "atbash":
		decodedText, decodeErr = cipher.BulkDecodeAtbashString(alphabetSet, *text, decodeToLatin)
		if decodeErr != nil {
			fmt.Printf("Failed to decode using Atbash cipher: %v", decodeErr)
		}
	case "trithemius":
		decodedText, decodeErr = cipher.BulkDecodeTrithemiusString(alphabetSet, *text, decodeToLatin)
		if decodeErr != nil {
			fmt.Printf("Failed to decode using Trithemius cipher: %v", decodeErr)
		}
	case "vigenere":
		if *wordFile == "" {
			log.Fatal("The -wordfile flag is required for Vigenere cipher")
		}

		// Read words from the CSV file
		wordlist, csvErr := ReadWordsFromCSVColumn(*wordFile, columnIndex)
		if csvErr != nil {
			return
		}

		decodedText, decodeErr = cipher.BulkDecodeVigenereCipher(alphabetSet, wordlist, *text, 1)

		// for i := 1; i <= 10; i++ {
		// decodedText, decodeErr = cipher.BulkDecodeVigenereCipher(alphabetSet, wordlist, *text, i)
		// if decodeErr != nil {
		// fmt.Printf("Failed to decode using Vigenere cipher: %v", decodeErr)
		// }
		// }
	}

	// Write the decoded text to the output file
	_, err = file.WriteString(fmt.Sprintf("Decoded Text: \n%s\n", decodedText))
}

// ReadWordsFromCSVColumn reads all the words from a specific column in a CSV file.
func ReadWordsFromCSVColumn(filePath string, columnIndex int) ([]string, error) {
	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read all rows from the CSV
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	// Extract words from the specified column
	var words []string
	for _, row := range rows {
		// Ensure the row has enough columns
		if columnIndex < len(row) {
			words = append(words, row[columnIndex])
		}
	}

	return words, nil
}
