package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"liberdatabase"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gorm.io/gorm"
)

// dbConn is a pointer to the gorm.DB instance for interacting with the PostgreSQL database.
var dbConn *gorm.DB

// main is the entry point of the application, initializes database connection, parses command-line flags, and processes text files.
func main() {
	dir := flag.String("dir", "", "The text to decode")

	// Parse the flags
	flag.Parse()

	dbConn, _ = liberdatabase.InitConnection()

	if err := walkAndProcess(*dir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// walkAndProcess traverses the directory tree starting at root and processes only .txt files using processTextFile.
func walkAndProcess(root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// If we can't access a file/dir, log and continue
			fmt.Fprintf(os.Stderr, "Skipping %s: %v\n", path, err)
			return nil
		}
		if d.IsDir() {
			return nil
		}

		// Only process .txt files
		if err := processTextFile(path); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to process %s: %v\n", path, err)
		}
		return nil
	})
}

// processTextFile processes a text file, tracks each word, and updates the database with word counts for that file.
func processTextFile(path string) error {
	lines, _ := readAllLines(path)

	var df liberdatabase.DocumentFile
	if liberdatabase.DoesDocumentFileExist(dbConn, path) {
		df, _ = liberdatabase.GetDocumentFile(dbConn, path)
	} else {
		df = liberdatabase.AddDocumentFile(dbConn, path)
	}

	for _, line := range lines {
		separators := extractSeparators(line)
		words := getAllWords(line, separators)

		for _, word := range words {
			if liberdatabase.DoesWordExist(dbConn, word, df.FileId) {
				liberdatabase.IncrementWordCount(dbConn, word, df.FileId)
			} else {
				liberdatabase.AddDocumentWord(dbConn, word, df.FileId, 1)
			}
		}
	}

	return nil
}

// readAllLines reads all lines from the specified file path, converts them to lowercase, and returns them as a slice of strings.
func readAllLines(path string) ([]string, error) {
	var lines []string
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file) // Ensure the file is closed

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			continue
		}

		line = strings.ToLower(line)
		lines = append(lines, line)
	}

	if scanError := scanner.Err(); scanError != nil {
		fmt.Printf("Error reading file %s: %v\n", path, scanError)
	}

	return lines, nil
}

// extractSeparators takes a string and returns a string containing only the non-alphabetic characters from the input.
func extractSeparators(text string) string {
	stringArray := strings.Split(text, "")
	lettersArray := strings.Split("abcdefghijklmnopqrstuvwxyz'", "")
	var retval []string

	for _, character := range stringArray {
		if !slices.Contains(lettersArray, character) {
			if !slices.Contains(retval, character) {
				retval = append(retval, character)
			}
		}
	}

	return strings.Join(retval, "")
}

// getAllWords splits a line of text into words based on the specified separators and returns a slice of words.
func getAllWords(line, separators string) []string {
	lineArray := strings.Split(line, "")
	var words []string
	var wordBuilder strings.Builder
	for _, character := range lineArray {
		if !strings.Contains(separators, character) {
			wordBuilder.WriteString(character)
		} else {
			if wordBuilder.Len() > 0 {
				words = append(words, wordBuilder.String())
			}

			wordBuilder.Reset()
		}
	}

	return words
}
