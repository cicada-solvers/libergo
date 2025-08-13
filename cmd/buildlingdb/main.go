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
	"runtime"
	"slices"
	"strings"
	"sync"

	"gorm.io/gorm"
)

var fileChannel chan string
var connections map[int]*gorm.DB
var lettersArray []string

// main is the entry point of the application, initializes database connection, parses command-line flags, and processes text files.
func main() {
	lettersArray = strings.Split("abcdefghijklmnopqrstuvwxyz'", "")
	fileChannel = make(chan string, 16384) // Increased buffer size

	dir := flag.String("dir", "", "The text to decode")

	// Parse the flags
	flag.Parse()

	// Threading for sentence processing.
	var wg sync.WaitGroup

	numWorkers := runtime.NumCPU() * 2 // Adjusted number of workers
	connections = make(map[int]*gorm.DB, numWorkers)
	for i := 0; i < numWorkers; i++ {
		connections[i], _ = liberdatabase.InitConnection()
		wg.Add(1)
		go processTextFileChannel(i, &wg)
	}

	go func() {
		err := walkAndProcess(*dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}
		close(fileChannel)
	}()

	wg.Wait()

	// Close the DB connections
	for i := 0; i < numWorkers; i++ {
		_ = liberdatabase.CloseConnection(connections[i])
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

		fileChannel <- path

		return nil
	})
}

func processTextFileChannel(workerId int, wg *sync.WaitGroup) {
	for document := range fileChannel {
		err := processTextFile(document, workerId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing file %s: %v\n", document, err)
			continue
		}
	}

	wg.Done()
}

// processTextFile processes a text file, tracks each word, and updates the database with word counts for that file.
func processTextFile(path string, workerId int) error {
	dbConn := connections[workerId]
	lines, _ := readAllLines(path)

	var df liberdatabase.DocumentFile
	if liberdatabase.DoesDocumentFileExist(dbConn, path) {
		df, _ = liberdatabase.GetDocumentFile(dbConn, path)
	} else {
		df = liberdatabase.AddDocumentFile(dbConn, path)
	}

	for _, line := range lines {
		words := getAllWords(line)

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
	buf := make([]byte, 64*1024)
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

// getAllWords splits a line of text into words based on the specified separators and returns a slice of words.
func getAllWords(line string) []string {
	lineArray := strings.Split(line, "")

	var words []string
	var wordBuilder strings.Builder
	for _, character := range lineArray {
		if slices.Contains(lettersArray, character) {
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
