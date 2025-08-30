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
	"sort"
	"strings"
	"sync"

	"gorm.io/gorm"
)

type FileStruct struct {
	FileName string
	Size     int64
}

var fileList []FileStruct
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

	// Get FileList
	fmt.Printf("Getting file list from %s\n", *dir)
	err := getFileList(*dir)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	fmt.Printf("Found %d files\n", len(fileList))
	sortFileListAscending()

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
		for _, file := range fileList {
			fileChannel <- file.FileName
		}
		close(fileChannel)
	}()

	wg.Wait()

	// Get All the Words
	liberdatabase.DeleteAllWordStatistics(connections[0])
	currentFileId := uint(0)
	wordBatch := make([]liberdatabase.WordStatistics, 0, 500)
	distinctWords := liberdatabase.GetAllDistinctWords(connections[0], 0)

	for len(distinctWords) > 0 {
		for _, dw := range distinctWords {
			fmt.Printf("Processing word %s\n", dw.Word)
			if dw.ID > currentFileId {
				currentFileId = dw.ID
			}

			average := liberdatabase.GetAveraegePercentageOfTextByWord(connections[0], dw.Word)

			wordBatch = append(wordBatch, liberdatabase.WordStatistics{
				Word:                    dw.Word,
				AveragePercentageOfText: average,
			})

			if len(wordBatch) >= 500 {
				liberdatabase.AddWordStatistics(connections[0], wordBatch)
				wordBatch = []liberdatabase.WordStatistics{}
			}
		}

		distinctWords = liberdatabase.GetAllDistinctWords(connections[0], currentFileId)
	}

	// Close the DB connections
	for i := 0; i < numWorkers; i++ {
		_ = liberdatabase.CloseConnection(connections[i])
	}
}

// walkAndProcess traverses the directory tree starting at root and processes only .txt files using processTextFile.
func getFileList(root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// If we can't access a file/dir, log and continue
			_, _ = fmt.Fprintf(os.Stderr, "Skipping %s: %v\n", path, err)
			return nil
		}
		if d.IsDir() {
			return nil
		}

		fileInfo, _ := os.Stat(path)

		file := FileStruct{
			FileName: path,
			Size:     fileInfo.Size(),
		}

		fileList = append(fileList, file)

		return nil
	})
}

func sortFileListAscending() {
	sort.Slice(fileList, func(i, j int) bool {
		return fileList[i].Size < fileList[j].Size
	})
}

func processTextFileChannel(workerId int, wg *sync.WaitGroup) {
	for document := range fileChannel {
		err := processTextFile(document, workerId)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error processing file %s: %v\n", document, err)
			continue
		}

		// Remove file after processing
		_ = os.Remove(document)
	}

	wg.Done()
}

func processTextFile(path string, workerId int) error {
	mapOfWords := make(map[string]int64)
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
		liberdatabase.DeleteStatisticsByFileId(dbConn, df.FileId)
		liberdatabase.DeleteWordsByFileId(dbConn, df.FileId)
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
			_, keyExists := mapOfWords[word]
			if keyExists {
				mapOfWords[word]++
			} else {
				mapOfWords[word] = 1
			}
		}
	}

	if scanError := scanner.Err(); scanError != nil {
		return fmt.Errorf("error reading file %s: %w", path, scanError)
	}

	tempWords := make([]liberdatabase.DocumentWord, 0, len(mapOfWords))
	for word, count := range mapOfWords {
		if len(tempWords) >= 500 {
			liberdatabase.AddDocumentWord(dbConn, tempWords)
			tempWords = []liberdatabase.DocumentWord{}
		} else {
			tempWord := liberdatabase.DocumentWord{
				FileId:    df.FileId,
				Word:      word,
				WordCount: count,
			}

			tempWords = append(tempWords, tempWord)
		}
	}

	liberdatabase.AddDocumentWord(dbConn, tempWords)
	tempWords = []liberdatabase.DocumentWord{}

	calculateWordPercentages(dbConn, df.FileId)

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

func calculateWordPercentages(dbConn *gorm.DB, fileId string) {
	words := liberdatabase.GetDistinctWords(dbConn, fileId)
	totalCount := int64(0)
	percentages := make([]liberdatabase.DocumentWordStatistics, 0, len(words))

	for _, word := range words {
		totalCount += word.WordCount
	}

	for _, word := range words {
		wordPercent := (float64(word.WordCount) / float64(totalCount)) * 100

		documentPercent := liberdatabase.DocumentWordStatistics{
			Word:             word.Word,
			PercentageOfText: wordPercent,
			FileId:           fileId,
		}

		percentages = append(percentages, documentPercent)

		if len(percentages) >= 500 {
			liberdatabase.AddDocumentWordStatistics(dbConn, percentages)
			percentages = []liberdatabase.DocumentWordStatistics{}
		}
	}

	liberdatabase.AddDocumentWordStatistics(dbConn, percentages)
	percentages = []liberdatabase.DocumentWordStatistics{}
}
