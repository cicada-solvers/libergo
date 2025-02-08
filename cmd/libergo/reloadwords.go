package main

import (
	"bufio"
	"config"
	"fmt"
	"liberdatabase"
	"os"
	"path/filepath"
	"runer"
	"strings"
	"unicode/utf8"
)

func ReloadWords() {
	// Get the configuration folder path
	configFolderPath, err := config.GetConfigFolderPath()
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error getting config folder path: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}

	// Construct the path to the words.txt file
	wordsFilePath := filepath.Join(configFolderPath, "words.txt")

	// Open the words.txt file
	file, err := os.Open(wordsFilePath)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error opening words.txt file: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Error closing words.txt file: %v\n", err)
			if err != nil {
				return
			}
		}
	}(file)

	// Get the database connection
	db, err := liberdatabase.InitConnection()

	// Clear the table
	err = liberdatabase.DeleteAllDictionaryWords(db)

	// Read each line from the words.txt file and write to the console
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		latinText := strings.ToUpper(scanner.Text())
		runeglishText := runer.PrepLatinToRune(latinText)
		runeText := runer.TransposeLatinToRune(runeglishText)
		gemSum := runer.CalculateGemSum(runeText, runer.Runes)

		fmt.Printf("Loading Word: %s Runeglish: %s Runes: %s Gem Sum: %d\n", latinText, runeglishText, runeText, gemSum)

		dictionaryWord := liberdatabase.DictionaryWord{
			DictionaryWordText:   latinText,
			RuneglishWordText:    runeglishText,
			RuneWordText:         runeText,
			GemSum:               gemSum,
			DictionaryWordLength: len(latinText),
			RuneglishWordLength:  len(runeglishText),
			RuneWordLength:       utf8.RuneCountInString(runeText),
		}

		dictionaryWord.RunePattern = dictionaryWord.GetRunePattern()

		err = liberdatabase.InsertDictionaryWord(db, dictionaryWord)
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error reading words.txt file: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}
