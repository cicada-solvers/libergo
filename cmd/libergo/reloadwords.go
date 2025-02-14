package main

import (
	"bufio"
	"config"
	"fmt"
	"gorm.io/gorm"
	"liberdatabase"
	"os"
	"path/filepath"
	"runer"
	"strings"
	"sync"
	"unicode/utf8"
)

func ReloadWords() {
	// Get the configuration folder path
	configFolderPath, err := config.GetConfigFolderPath()
	if err != nil {
		_, err2 := fmt.Fprintf(os.Stderr, "Error getting config folder path: %v\n", err)
		if err2 != nil {
			return
		}
		os.Exit(1)
	}

	// Construct the path to the words.txt file
	wordsFilePath := filepath.Join(configFolderPath, "words.txt")

	// Open the words.txt file
	file, err := os.Open(wordsFilePath)
	if err != nil {
		_, err3 := fmt.Fprintf(os.Stderr, "Error opening words.txt file: %v\n", err)
		if err3 != nil {
			return
		}
		os.Exit(1)
	}
	defer func(file *os.File) {
		fileErr := file.Close()
		if fileErr != nil {
			_, err4 := fmt.Fprintf(os.Stderr, "Error closing words.txt file: %v\n", err)
			if err4 != nil {
				return
			}
			os.Exit(1)
		}
	}(file)

	// Read all lines into memory
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		_, err5 := fmt.Fprintf(os.Stderr, "Error reading words.txt file: %v\n", err)
		if err5 != nil {
			return
		}
		os.Exit(1)
	}

	// Get the number of workers from the configuration
	configuration, _ := config.LoadConfig()
	numWorkers := configuration.NumWorkers

	// Create a wait group and channels for the words and errors
	var wg sync.WaitGroup
	wordChan := make(chan string, len(lines))
	errChan := make(chan error, len(lines))

	// Start a fixed number of worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Get the database connection for this worker
			db, err := liberdatabase.InitConnection()
			if err != nil {
				errChan <- fmt.Errorf("error initializing database connection: %v", err)
				return
			}
			defer func(db *gorm.DB) {
				closeErr := liberdatabase.CloseConnection(db)
				if closeErr != nil {
					errChan <- fmt.Errorf("error closing database connection: %v", err)
				}
			}(db)

			for latinText := range wordChan {
				runeglishText := runer.PrepLatinToRune(strings.ToUpper(latinText))
				runeText := runer.TransposeLatinToRune(runeglishText)
				gemSum := runer.CalculateGemSum(runeText, runer.Runes)

				fmt.Printf("Loading Word: %s Runeglish: %s Runes: %s Gem Sum: %d\n", strings.ToUpper(latinText), runeglishText, runeText, gemSum)

				dictionaryWord := liberdatabase.DictionaryWord{
					DictionaryWordText:   strings.ToUpper(latinText),
					RuneglishWordText:    runeglishText,
					RuneWordText:         runeText,
					GemSum:               gemSum,
					DictionaryWordLength: len(latinText),
					RuneglishWordLength:  len(runeglishText),
					RuneWordLength:       utf8.RuneCountInString(runeText),
				}

				dictionaryWord.RunePattern = dictionaryWord.GetRunePattern()

				err = liberdatabase.InsertDictionaryWord(db, dictionaryWord)
				if err != nil {
					errChan <- fmt.Errorf("error inserting dictionary word: %v", err)
				}
			}
		}()
	}

	// Send the words to the channel
	for _, line := range lines {
		wordChan <- line
	}
	close(wordChan)

	// Wait for all goroutines to finish
	wg.Wait()
	close(errChan)

	// Handle errors
	for err := range errChan {
		_, charErr := fmt.Fprintf(os.Stderr, "%v\n", err)
		if charErr != nil {
			return
		}
	}
}
