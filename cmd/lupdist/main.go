package main

import (
	runelib "characterrepo"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"lgstructs"
	"log"
	"os"
	"runer"
	"strconv"
	"strings"
)

type LiberWordDistance struct {
	LiberWordGuid     string
	LiberWord         string
	LiberWordLength   int
	LiberWordPosition int
	LiberWordSection  string
	Patterns          []WordDistancePattern
}

type WordDistancePattern struct {
	WordDistancePatternGuid       string
	DictionaryWord                string
	DictionaryWordDistancePattern string
	WordListOrigin                string
	LiberWordGuid                 string
	TranslatedLatin               string
}

func main() {
	charRepo := runelib.NewCharacterRepo()
	runes := charRepo.GetGematriaRunes()

	section := flag.String("section", "", "The section to measure")
	text := flag.String("text", "", "The text to measure")
	wordFile := flag.String("wordfile", "", "The CSV file of words to try for measuring")
	outputFile := flag.String("output", "", "The output file to write the results")

	// Parse the flags
	flag.Parse()

	// Validate required flags
	if *text == "" {
		log.Fatal("The -text flag is required")
	}
	if *outputFile == "" {
		log.Fatal("The -output flag is required")
	}
	if *wordFile == "" {
		log.Fatal("The -wordfile flag is required")
	}
	if *section == "" {
		log.Fatal("The -wordfile flag is required")
	}

	fmt.Printf("Section: %s\n", *section)
	fmt.Printf("Text: %s\n", *text)
	fmt.Printf("Output File: %s\n", *outputFile)
	fmt.Printf("Word File: %s\n", *wordFile)

	// Gets the base file name without the extension
	baseFileName := strings.TrimSuffix(*outputFile, ".csv")
	headerFileName := fmt.Sprintf("%s_header.csv", baseFileName)
	detailsFileName := fmt.Sprintf("%s_detail.csv", baseFileName)

	words := breakTextApart(strings.Split(*text, ""))

	for wordPos, word := range words {
		wordDistance := LiberWordDistance{
			LiberWordGuid:     uuid.NewString(),
			LiberWord:         word,
			LiberWordLength:   len(strings.Split(word, "")),
			LiberWordPosition: wordPos + 1,
			LiberWordSection:  *section,
			Patterns:          []WordDistancePattern{},
		}

		wordLength := len(strings.Split(word, ""))
		listWords, _ := ReadWordsFromCSVColumn(*wordFile, 2, wordLength)
		for _, listWord := range listWords {
			distancePattern := lgstructs.CalculateWordDistances(strings.Split(word, ""), strings.Split(listWord, ""), runes)
			pattern := WordDistancePattern{
				DictionaryWord:                listWord,
				DictionaryWordDistancePattern: distancePattern,
				WordDistancePatternGuid:       uuid.NewString(),
				WordListOrigin:                *wordFile,
				LiberWordGuid:                 wordDistance.LiberWordGuid,
				TranslatedLatin:               runer.TransposeRuneToLatin(listWord),
			}

			wordDistance.Patterns = append(wordDistance.Patterns, pattern)
		}

		writeCsvHeaderFile(headerFileName, wordDistance)
		writeCsvDetailFile(detailsFileName, wordDistance)

		fmt.Printf("Word: %s, Word Length: %d, Word Position: %d, Section: %s\n", wordDistance.LiberWord, wordDistance.LiberWordLength, wordDistance.LiberWordPosition, wordDistance.LiberWordSection)
	}

	// Now we are going to write the CSV files

	fmt.Printf("Done\n\n\n")
}

func writeCsvHeaderFile(headerFile string, word LiberWordDistance) {
	// Create the CSV file
	file, err := os.OpenFile(headerFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Failed to create output file: %v\n", err)
		return
	}
	defer func(file *os.File) {
		fileError := file.Close()
		if fileError != nil {
			fmt.Printf("Failed to close output file: %v\n", fileError)
		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	fileInfo, _ := file.Stat()
	if fileInfo.Size() == 0 {
		// Write the header
		header := []string{
			"LiberWordGuid", "LiberWord", "LiberWordLength", "LiberWordPosition", "LiberWordSection",
		}
		if err := writer.Write(header); err != nil {
			fmt.Printf("Failed to write header to output file: %v\n", err)
			return
		}
	}

	// Write the data
	record := []string{
		word.LiberWordGuid,
		word.LiberWord,
		strconv.Itoa(word.LiberWordLength),
		strconv.Itoa(word.LiberWordPosition),
		word.LiberWordSection,
	}
	if err := writer.Write(record); err != nil {
		fmt.Printf("Failed to write record to output file: %v\n", err)
		return
	}
}

func writeCsvDetailFile(detailFile string, word LiberWordDistance) {
	// Create the CSV file
	file, err := os.OpenFile(detailFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Failed to create output file: %v\n", err)
		return
	}
	defer func(file *os.File) {
		fileError := file.Close()
		if fileError != nil {
			fmt.Printf("Failed to close output file: %v\n", fileError)
		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	fileInfo, _ := file.Stat()
	if fileInfo.Size() == 0 {
		// Write the header
		header := []string{
			"DictionaryWord", "DictionaryWordDistancePattern", "WordDistancePatternGuid", "WordListOrigin", "LiberWordGuid", "TranslatedLatin",
		}
		if err := writer.Write(header); err != nil {
			fmt.Printf("Failed to write header to output file: %v\n", err)
			return
		}
	}

	// Write the data
	for _, pattern := range word.Patterns {
		record := []string{
			pattern.DictionaryWord,
			pattern.DictionaryWordDistancePattern,
			pattern.WordDistancePatternGuid,
			pattern.WordListOrigin,
			pattern.LiberWordGuid,
			pattern.TranslatedLatin,
		}
		if err := writer.Write(record); err != nil {
			fmt.Printf("Failed to write record to output file: %v\n", err)
			return
		}
	}
}

func breakTextApart(text []string) []string {
	var words []string
	var combinedText strings.Builder
	charrepo := runelib.NewCharacterRepo()

	for _, character := range text {
		if charrepo.IsRune(character, false) {
			combinedText.WriteString(character)
		} else {
			combinedText.WriteString(" ")
		}
	}

	tempString := combinedText.String()
	tempString = strings.ReplaceAll(tempString, "  ", " ")
	tempString = strings.TrimSpace(tempString)

	// Split the combined text into words
	words = strings.Fields(tempString)

	return words
}

// ReadWordsFromCSVColumn reads all the words from a specific column in a CSV file.
func ReadWordsFromCSVColumn(filePath string, columnIndex int, length int) ([]string, error) {
	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			fmt.Printf("Failed to close file: %v", err)
		}
	}(file)

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
			wordLength := len(strings.Split(row[columnIndex], ""))
			if wordLength == length {
				words = append(words, row[columnIndex])
			}
		}
	}

	return words, nil
}
