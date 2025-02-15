package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/big"
	"os"
	"runer"
	"strings"
	"titler"
)

// ParseTextType converts a string to the corresponding TextType
func ParseTextType(textType string) (runer.TextType, error) {
	switch strings.ToLower(textType) {
	case "latin":
		return runer.Latin, nil
	case "runeglish":
		return runer.Runeglish, nil
	case "runes":
		return runer.Runes, nil
	default:
		return runer.Latin, fmt.Errorf("invalid text type: %s", textType)
	}
}

// CalculateGemProductForText calculates the gem product for the entire text and each word
func CalculateGemProductForText(input string, outputToFile bool, outputFileName string, textType runer.TextType) error {
	var text string
	if fileInfo, err := os.Stat(input); err == nil && !fileInfo.IsDir() {
		file, err := os.Open(input)
		if err != nil {
			return err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}(file)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			text += scanner.Text() + " "
		}
		if err := scanner.Err(); err != nil {
			return err
		}
	} else {
		text = input
	}

	// Split text into words using space, •, and ⊹ as delimiters
	words := strings.FieldsFunc(text, func(r rune) bool {
		return r == ' ' || r == '•' || r == '⊹' || r == '.' || r == ',' || r == '!' || r == '?' || r == ':' || r == ';' || r == '(' || r == ')'
	})

	totalGemProduct := runer.CalculateGemProduct(text, textType)
	wordGemProducts := make(map[string]big.Int)
	for _, word := range words {
		wordGemProducts[word] = runer.CalculateGemProduct(word, textType)
	}

	if outputToFile {
		file, err := os.Create(outputFileName)
		if err != nil {
			return err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}(file)
		writer := bufio.NewWriter(file)
		_, fileError := fmt.Fprintf(writer, "Total Gem Product: %d\n", totalGemProduct)
		if fileError != nil {
			fmt.Printf("Error: %v\n", fileError)
		}

		for word, gemProduct := range wordGemProducts {
			_, fileError = fmt.Fprintf(writer, "Word: %s, Gem Product: %d\n", word, gemProduct)
			if fileError != nil {
				fmt.Printf("Error: %v\n", fileError)
			}
		}
		flushError := writer.Flush()
		if flushError != nil {
			fmt.Printf("Error: %v\n", flushError)
		}
	} else {
		fmt.Printf("Total Gem Product: %d\n", totalGemProduct)
		for word, gemProduct := range wordGemProducts {
			fmt.Printf("Word: %s, Gem Product: %d\n", word, gemProduct)
		}
	}

	return nil
}

func main() {
	titler.PrintTitle("Gematria Product")
	var outputToFile = false
	input := flag.String("input", "", "Input string or file path")
	outputFileName := flag.String("outputFileName", "", "Output file name (optional)")
	textTypeStr := flag.String("textType", "latin", "Type of text (e.g., latin, runeglish, runes)")

	flag.Parse()

	if *input == "" {
		fmt.Println("Input is required")
		flag.Usage()
		return
	}

	if *outputFileName == "" {
		outputToFile = false
	}

	textType, err := ParseTextType(*textTypeStr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	err = CalculateGemProductForText(*input, outputToFile, *outputFileName, textType)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
