package main

import (
	"flag"
	"fmt"
	"github.com/jdkato/prose/v2"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"
)

var fileMutex sync.Mutex

type ColInformation struct {
	ColName   string
	RowCounts int
}

type Sentence struct {
	Content     string
	Output      string
	ColumnIndex int
}

func main() {
	// Define command-line flags
	inputFile := flag.String("input", "", "Path to the input Excel file")
	outputFile := flag.String("output", "", "Path to the output file")
	sheetName := "Worksheet"

	// Parse the flags
	flag.Parse()

	// Check if the input file is provided
	if *inputFile == "" {
		log.Fatalf("Input file is required")
	}

	// Open the Excel file
	f, err := excelize.OpenFile(*inputFile)
	if err != nil {
		log.Fatalf("Failed to open the Excel file: %v", err)
	}
	defer func(f *excelize.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("Failed to close the Excel file: %v", err)
		}
	}(f)

	colInfo, err := getColInformation(f, sheetName)
	if err != nil {
		fmt.Printf("Failed to get column info: %v", err)
		return
	}

	// Print the column information
	fmt.Printf("%v\n", colInfo)

	// Initialize a strings.Builder
	var builder strings.Builder

	// Call permuteCols with the provided output file name
	err = permuteCols(f, *outputFile, sheetName, colInfo, builder, 0)
	if err != nil {
		fmt.Printf("Failed to permute cols: %v", err)
	}
}

func permuteCols(f *excelize.File, outputName, sheetName string, cols []ColInformation, builder strings.Builder, currentColIdx int) (err error) {
	localBuilder := cloneStringBuilder(&builder)

	if currentColIdx < (len(cols) - 1) {
		for i := 3; i < cols[currentColIdx].RowCounts+3; i++ {
			localBuilder = cloneStringBuilder(&builder)
			cellValue, cellError := f.GetCellValue(sheetName, fmt.Sprintf("%s%d", cols[currentColIdx].ColName, i))
			if cellError != nil {
				return fmt.Errorf("Failed to get cell value: %v\n", cellError)
			}

			spacer := " "
			if currentColIdx == 0 {
				spacer = ""
			}

			localBuilder.WriteString(spacer + cellValue)

			// Write the builder content to the console
			fmt.Printf("%d %s\n", currentColIdx, localBuilder.String())

			permuteErr := permuteCols(f, outputName, sheetName, cols, *localBuilder, currentColIdx+1)
			if permuteErr != nil {
				return fmt.Errorf("Failed to permute columns: %v\n", permuteErr)
			}
		}
	} else {
		var wg sync.WaitGroup
		sentenceChan := make(chan Sentence, 16384) // Increased buffer size

		go func() {
			for i := 3; i < cols[currentColIdx].RowCounts+3; i++ {
				cellValue, cellError := f.GetCellValue(sheetName, fmt.Sprintf("%s%d", cols[currentColIdx].ColName, i))
				if cellError != nil {
					fmt.Printf("Failed to get cell value: %v\n", cellError)
				}

				localBuilder = cloneStringBuilder(&builder)

				localBuilder.WriteString(" " + cellValue)

				sentence := Sentence{
					Content:     localBuilder.String(),
					Output:      outputName,
					ColumnIndex: currentColIdx,
				}

				sentenceChan <- sentence
			}
			close(sentenceChan)
		}()

		numWorkers := runtime.NumCPU() * 2 // Adjusted number of workers
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go calculateProbabilityAndWriteToFile(sentenceChan, &wg)
		}

		wg.Wait()
	}

	return nil
}

func calculateProbabilityAndWriteToFile(sentChan chan Sentence, wg *sync.WaitGroup) {
	defer wg.Done()

	for sentence := range sentChan {
		// Write the builder content to the console
		fmt.Printf("%d %s\n", sentence.ColumnIndex, sentence.Content)

		posCounts, totalWords := analyzeText(sentence.Content)
		probability := calculateSentenceProbability(posCounts, totalWords)

		fmt.Printf("POS Counts: %+v\n", posCounts)
		fmt.Printf("Total Words: %d\n", totalWords)

		if probability > 0 {
			fmt.Printf("Sentence Probability: %.2f%%\n", probability)

			// Write the content to the output file
			outputBytes := []byte(sentence.Content + "\n\n")

			for {
				fileMutex.Lock()
				file, openError := os.OpenFile(sentence.Output, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
				if openError != nil {
					fmt.Printf("Failed to open file: %v\n", openError)
					fileMutex.Unlock()
					time.Sleep(100 * time.Millisecond) // Wait before retrying
					continue
				}

				_, writeErr := file.Write(outputBytes)
				if writeErr != nil {
					fmt.Printf("Failed to write to file: %v\n", writeErr)
					err := file.Close()
					if err != nil {
						fmt.Printf("Failed to close file: %v\n", err)
					}
					fileMutex.Unlock()
					time.Sleep(100 * time.Millisecond) // Wait before retrying
					continue
				}

				closeError := file.Close()
				if closeError != nil {
					fmt.Printf("Failed to close file: %v\n", closeError)
				}
				fileMutex.Unlock()
				break
			}
		}
	}
}

func analyzeText(text string) (map[string]int, int) {
	doc, err := prose.NewDocument(text)
	if err != nil {
		log.Fatalf("Failed to create document: %v", err)
	}

	posCounts := map[string]int{
		"Noun":        0,
		"Verb":        0,
		"Adjective":   0,
		"Adverb":      0,
		"Determiner":  0,
		"Conjunction": 0,
		"Preposition": 0,
		"Pronoun":     0,
		"Punctuation": 0,
		"NamedEntity": 0,
	}
	totalWords := 0

	for _, tok := range doc.Tokens() {
		switch tok.Tag {
		case "NN", "NNS", "NNP", "NNPS":
			posCounts["Noun"]++
		case "VB", "VBD", "VBG", "VBN", "VBP", "VBZ":
			posCounts["Verb"]++
		case "JJ", "JJR", "JJS":
			posCounts["Adjective"]++
		case "RB", "RBR", "RBS":
			posCounts["Adverb"]++
		case "DT":
			posCounts["Determiner"]++
		case "CC":
			posCounts["Conjunction"]++
		case "IN":
			posCounts["Preposition"]++
		case "PRP", "PRP$", "WP", "WP$":
			posCounts["Pronoun"]++
		case ".", ",", ":", ";", "!", "?":
			posCounts["Punctuation"]++
		}
		totalWords++
	}

	posCounts["NamedEntity"] = len(doc.Entities())

	return posCounts, totalWords
}

func calculateSentenceProbability(posCounts map[string]int, totalWords int) float64 {
	if totalWords == 0 {
		return 0.0
	}

	probability := 0.0
	if posCounts["Noun"] > 0 && posCounts["Verb"] > 0 {
		probability = 50.0
		if posCounts["Adjective"] > 0 {
			probability += 10.0
		}
		if posCounts["Adverb"] > 0 {
			probability += 10.0
		}
		if posCounts["Determiner"] > 0 {
			probability += 5.0
		}
		if posCounts["Conjunction"] > 0 {
			probability += 5.0
		}
		if posCounts["Preposition"] > 0 {
			probability += 5.0
		}
		if posCounts["Pronoun"] > 0 {
			probability += 5.0
		}
		if posCounts["Punctuation"] > 0 {
			probability += 10.0
		}
		if posCounts["NamedEntity"] > 0 {
			probability += 5.0
		}
	}

	return probability
}

// cloneStringBuilder clones the given strings.Builder and returns a new instance with the same content.
func cloneStringBuilder(sb *strings.Builder) *strings.Builder {
	// Create a new strings.Builder
	newSb := &strings.Builder{}
	// Write the content of the original builder to the new builder
	newSb.WriteString(sb.String())
	return newSb
}

// getColInformation gets the column information from the given Excel file and sheet name.
func getColInformation(f *excelize.File, sheetName string) (cols []ColInformation, err error) {
	// Get the column information
	cols = make([]ColInformation, 0)

	// Get all the columns in the sheet
	columns, err := f.GetCols(sheetName)
	if err != nil {
		return nil, fmt.Errorf("Failed to get columns: %v\n", err)
	}

	for colId, colValues := range columns {
		// Get the row count
		rowCounts := len(colValues)
		colName := intToExcelColumn(colId + 1)
		colInfo := ColInformation{
			ColName:   colName,
			RowCounts: 0,
		}

		for i := 3; i < rowCounts+5; i++ {
			cellValue, cellError := f.GetCellValue(sheetName, fmt.Sprintf("%s%d", colName, i))
			if cellError != nil {
				return nil, fmt.Errorf("Failed to get cell value: %v\n", cellError)
			}

			if cellValue != "" {
				colInfo.RowCounts++
			} else {
				break
			}
		}

		cols = append(cols, colInfo)
	}

	return cols, nil
}

// intToExcelColumn converts the given integer to an Excel column name.
func intToExcelColumn(n int) string {
	column := ""
	for n > 0 {
		n-- // Adjust for 1-based indexing
		column = string('A'+(n%26)) + column
		n /= 26
	}
	return column
}
