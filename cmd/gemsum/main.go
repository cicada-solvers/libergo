package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runer"
	"strings"
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

// CalculateGemSumForText calculates the gem sum for the entire text and each word
func CalculateGemSumForText(input string, outputToFile bool, outputFileName string, textType runer.TextType) error {
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

	totalGemSum := runer.CalculateGemSum(text, textType)
	wordGemSums := make(map[string]int64)
	for _, word := range words {
		wordGemSums[word] = runer.CalculateGemSum(word, textType)
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
		_, fileError := fmt.Fprintf(writer, "Total Gem Sum: %d\n", totalGemSum)
		if fileError != nil {
			fmt.Printf("Error: %v\n", fileError)
		}

		for word, gemSum := range wordGemSums {
			_, fileError = fmt.Fprintf(writer, "Word: %s, Gem Sum: %d\n", word, gemSum)
			if fileError != nil {
				fmt.Printf("Error: %v\n", fileError)
			}
		}
		flushError := writer.Flush()
		if flushError != nil {
			fmt.Printf("Error: %v\n", flushError)
		}
	} else {
		fmt.Printf("Total Gem Sum: %d\n", totalGemSum)
		for word, gemSum := range wordGemSums {
			fmt.Printf("Word: %s, Gem Sum: %d\n", word, gemSum)
		}
	}

	return nil
}

func main() {
	var outputToFile = false
	input := flag.String("input", "", "Input string or file path")
	outputFileName := flag.String("outputFileName", "", "Output file name (optional)")
	textTypeStr := flag.String("textType", "latin", "Type of text (e.g., latin, runeglish, runes)")

	flag.Parse()

	if *input == "" {
		fmt.Println("Input is required")
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

	err = CalculateGemSumForText(*input, outputToFile, *outputFileName, textType)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
