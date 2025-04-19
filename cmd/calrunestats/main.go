package main

import (
	"bufio"
	runelib "characterrepo"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/jdkato/prose/v2"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"runer"
	"runtime"
	"strings"
	"sync"
	"time"
)

// LineStatistics is a struct that holds the statistics of a line.
type LineStatistics struct {
	FileName            string
	KeyWord             string
	Text                string
	LatinText           string
	DoubletCount        int
	SentenceProbability float64
	Statistics          []RuneStatistic
}

// RuneStatistic is a struct that holds the rune and its percentage.
type RuneStatistic struct {
	Rune       string
	Percentage float64
}

// LineTextWithFile is a struct that holds the line text and the file name.
type LineTextWithFile struct {
	LineText string
	FileName string
}

// processedCounter is used to track the number of processed items.
var processedCounter = big.NewInt(0)

// rateCounter is used to track the rate of processed items.
var rateCounter = big.NewInt(0)

// main function initializes the program, processes files, and writes results to CSV.
func main() {
	// Initialize the repositories
	charRepo := runelib.NewCharacterRepo()
	gemRunes := charRepo.GetGematriaRunes()
	processedTicker := time.NewTicker(time.Minute)
	defer processedTicker.Stop()

	go func() {
		for range processedTicker.C {
			fmt.Printf("Rate: %s/min - Processed %s items\n", rateCounter.String(), processedCounter.String())
			rateCounter.SetInt64(int64(0))
		}
	}()

	// Define the flag for the directory path
	dirPath := flag.String("dir", "./your-directory-path", "Path to the directory containing text files")
	flag.Parse()

	// Create channels and a WaitGroup
	lineChan := make(chan LineTextWithFile, 100)
	resultsChan := make(chan LineStatistics, 100)
	var wg sync.WaitGroup

	// Start worker goroutines
	numWorkers := runtime.NumCPU() // Adjust the number of workers as needed
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range lineChan {
				lineStatistics := processLine(line, gemRunes, charRepo)
				resultsChan <- lineStatistics // Send results to the results channel
			}
		}()
	}

	// Start a single writer goroutine
	go func() {
		for result := range resultsChan {
			processedCounter.Add(processedCounter, big.NewInt(1))
			rateCounter.Add(rateCounter, big.NewInt(1))

			if result.SentenceProbability >= 50.0 {
				fmt.Printf("Writing results for file: %s\n", result.FileName)
				writeErr := writeResultsToFile(len(gemRunes), result, result.FileName+".csv")
				if writeErr != nil {
					fmt.Printf("Error writing results: %v\n", writeErr)
				}
			}
		}
	}()

	// Walk through the directory and send lines to the channel
	err := filepath.WalkDir(*dirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Open the file
		f, openErr := os.Open(path)
		if openErr != nil {
			fmt.Printf("Failed to open file %s: %v\n", path, openErr)
			return nil
		} else {
			fmt.Printf("Opened file %s successfully.\n", path)
		}

		defer func(f *os.File) {
			closeError := f.Close()
			if closeError != nil {
				fmt.Printf("Error closing file %s: %v\n", path, closeError)
			} else {
				fmt.Printf("File %s closed successfully.\n", path)
				removeErr := os.Remove(path)
				if removeErr != nil {
					fmt.Printf("Error removing file %s: %v\n", path, removeErr)
				}
			}
		}(f)

		// Read the file line by line
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := LineTextWithFile{
				LineText: scanner.Text(),
				FileName: path,
			}

			lineChan <- line // Send line to the channel
		}

		// Check for errors during scanning
		if scanErr := scanner.Err(); scanErr != nil {
			fmt.Printf("Error reading file %s: %v\n", path, scanErr)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the directory: %v\n", err)
	}

	// Close the line channel and wait for workers to finish
	close(lineChan)
	wg.Wait()

	// Close the results channel after all workers are done
	close(resultsChan)
}

// processLine processes a line of text and returns the statistics.
func processLine(line LineTextWithFile, gemRunes []string, charRepo *runelib.CharacterRepo) LineStatistics {
	var lineStatistics []RuneStatistic
	lineStatisticsAll := LineStatistics{}
	parts := strings.Split(line.LineText, " : ")
	if len(parts) > 1 {
		latinText := runer.TransposeRuneToLatin(parts[1])
		posCounts, totalWords := analyzeText(parts[1])
		sentenceProbability := calculateSentenceProbability(posCounts, totalWords)

		partTwoArray := strings.Split(parts[1], "")
		totalCount := getTotalRunes(partTwoArray)
		if totalCount > 0 {
			for _, runeChar := range gemRunes {
				count := getCountOfParticularRune(partTwoArray, runeChar)
				percentage := (float64(count) / float64(totalCount)) * 100
				statistic := RuneStatistic{
					Rune:       runeChar,
					Percentage: percentage,
				}

				lineStatistics = append(lineStatistics, statistic)
			}

			lineStatisticsAll = LineStatistics{
				FileName:            line.FileName,
				KeyWord:             parts[0],
				Text:                parts[1],
				LatinText:           latinText,
				SentenceProbability: sentenceProbability,
				DoubletCount:        charRepo.GetDoubletCount(parts[1], gemRunes),
				Statistics:          lineStatistics,
			}
		}
	}

	return lineStatisticsAll
}

// getTotalRunes counts the total number of runes in the line.
func getTotalRunes(line []string) int {
	charRepo := runelib.NewCharacterRepo()
	total := 0
	for _, str := range line {
		if charRepo.IsRune(str, false) {
			total += len(str)
		}
	}
	return total
}

// getCountOfParticularRune counts the occurrences of a specific rune in the line.
func getCountOfParticularRune(line []string, rune string) int {
	charRepo := runelib.NewCharacterRepo()
	count := 0
	for _, str := range line {
		if charRepo.IsRune(str, false) && str == rune {
			count++
		}
	}
	return count
}

// writeResultsToFile writes the results to a CSV file.
func writeResultsToFile(alphabetCount int, results LineStatistics, outputPath string) error {
	// Open or create the result CSV file
	resultFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Error opening/creating result.csv: %v\n", err)
		return err
	}
	defer func(resultFile *os.File) {
		closeErr := resultFile.Close()
		if closeErr != nil {
			fmt.Printf("Error closing result.csv: %v\n", closeErr)
		}
	}(resultFile)

	// Create a CSV writer
	csvWriter := csv.NewWriter(resultFile)
	defer csvWriter.Flush()

	// Write the header if the file is empty
	fileInfo, _ := resultFile.Stat()
	if fileInfo.Size() == 0 {
		header := []string{"KeyWord", "DoubletCount", "SentenceProbability", "Text", "LatinText"}

		for i := 0; i < alphabetCount; i++ {
			header = append(header, fmt.Sprintf("Rune%d", i+1))
		}

		if headerErr := csvWriter.Write(header); headerErr != nil {
			fmt.Printf("Error writing header to CSV: %v\n", headerErr)
			return headerErr
		}
	}

	var record []string
	record = append(record, results.KeyWord)
	record = append(record, fmt.Sprintf("%d", results.DoubletCount))
	record = append(record, fmt.Sprintf("%.2f", results.SentenceProbability))
	record = append(record, results.Text)
	record = append(record, results.LatinText)

	// Write the data
	for _, stat := range results.Statistics {
		rText := stat.Rune
		pText := fmt.Sprintf("%.2f", stat.Percentage)
		aText := fmt.Sprintf("%s|%s", rText, pText)
		record = append(record, aText)
	}

	if csvErr := csvWriter.Write(record); csvErr != nil {
		fmt.Printf("Error writing record to CSV: %v\n", csvErr)
		return err
	}

	return nil
}

// analyzeText analyzes the given text and returns the part-of-speech counts and total word count.
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

// calculateSentenceProbability calculates the probability of a sentence being a valid English sentence.
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
