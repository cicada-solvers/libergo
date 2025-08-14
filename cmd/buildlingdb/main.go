package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"liberdatabase"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"gorm.io/gorm"
)

var fileChannel chan string
var connections map[int]*gorm.DB

// Create letterMap once at initialization instead of for each call
var letterMap map[rune]bool

func init() {
	lettersArray := strings.Split("abcdefghijklmnopqrstuvwxyz'", "")
	letterMap = make(map[rune]bool, len(lettersArray))
	for _, letter := range lettersArray {
		letterMap[rune(letter[0])] = true
	}
}

// main is the entry point of the application, initializes database connection, parses command-line flags, and processes text files.
func main() {
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
			_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
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
			_, _ = fmt.Fprintf(os.Stderr, "Skipping %s: %v\n", path, err)
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
			_, _ = fmt.Fprintf(os.Stderr, "Error processing file %s: %v\n", document, err)
			continue
		}
	}

	wg.Done()
}

func processTextFile(path string, workerId int) error {
	fmt.Printf("Processing file %s\n", path)
	dbConn := connections[workerId]

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	var df liberdatabase.DocumentFile
	if liberdatabase.DoesDocumentFileExist(dbConn, path) {
		df, _ = liberdatabase.GetDocumentFile(dbConn, path)
	} else {
		df = liberdatabase.AddDocumentFile(dbConn, path)
	}

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			continue
		}

		line = strings.ToLower(line)
		words := getAllWords(line)

		for _, word := range words {
			if liberdatabase.DoesWordExist(dbConn, word, df.FileId) {
				liberdatabase.IncrementWordCount(dbConn, word, df.FileId)
			} else {
				liberdatabase.AddDocumentWord(dbConn, word, df.FileId, 1)
			}
		}
	}

	if scanError := scanner.Err(); scanError != nil {
		return fmt.Errorf("error reading file %s: %w", path, scanError)
	}

	return nil
}

// getAllWords splits a line of text into words based on the specified separators and returns a slice of words.
func getAllWords(line string) []string {
	var words []string
	var wordBuilder strings.Builder

	// Pre-allocate space for words to reduce reallocations
	words = make([]string, 0, 16) // Assuming average of ~16 words per line

	// Iterate through runes directly
	for _, r := range line {
		if letterMap[r] {
			wordBuilder.WriteRune(r)
		} else if wordBuilder.Len() > 0 {
			words = append(words, wordBuilder.String())
			wordBuilder.Reset()
		}
	}

	// Add the last word if the line ends with a letter
	if wordBuilder.Len() > 0 {
		words = append(words, wordBuilder.String())
	}

	return words
}
