package main

import (
	"bufio"
	runelib "characterrepo"
	"flag"
	"fmt"
	"io/fs"
	"liberdatabase"
	"os"
	"path/filepath"
	"runer"
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

var charRepo *runelib.CharacterRepo
var isCharExtract *bool
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
	charRepo = runelib.NewCharacterRepo()
	fileChannel = make(chan string, 16384) // Increased buffer size

	dir := flag.String("dir", "", "The text to decode")
	isCharExtract = flag.Bool("charextract", false, "Extract characters from text")

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

		if *isCharExtract {
			err := processCharacters(document, workerId)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error processing file %s: %v\n", document, err)
				continue
			}
		} else {
			err := processTextFile(document, workerId)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error processing file %s: %v\n", document, err)
				continue
			}
		}

		// Remove file after processing
		_ = os.Remove(document)
	}

	wg.Done()
}

// processCharacters processes a text file and extracts characters from it.
func processCharacters(path string, workerId int) error {
	mapOfCharacters := make(map[string]int64)
	mapOfRunes := make(map[string]int64)
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
		liberdatabase.DeleteCharactersByFileId(dbConn, df.FileId)
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

		line = strings.ToUpper(line)
		line = runer.PrepLatinToRune(line)
		runes := runer.TransposeLatinToRune(line, false)
		runeArray := strings.Split(runes, "")
		lineArray := strings.Split(line, "")
		for _, character := range lineArray {
			if charRepo.IsLetterInAlphabet(character) {
				_, keyExists := mapOfCharacters[character]
				if keyExists {
					mapOfCharacters[character]++
				} else {
					mapOfCharacters[character] = 1
				}
			}
		}

		for _, character := range runeArray {
			if charRepo.IsRune(character, false) {
				_, keyExists := mapOfRunes[character]
				if keyExists {
					mapOfRunes[character]++
				} else {
					mapOfRunes[character] = 1
				}
			}
		}
	}

	if scanError := scanner.Err(); scanError != nil {
		return fmt.Errorf("error reading file %s: %w", path, scanError)
	}

	totalCharacterCount := int64(0)
	characterArray := make([]liberdatabase.DocumentCharacter, 0, len(mapOfCharacters))
	for character, count := range mapOfCharacters {
		documentCharacter := liberdatabase.DocumentCharacter{
			FileId:         df.FileId,
			Character:      character,
			CharacterCount: count,
			CharacterType:  "RUNEGLISH",
		}

		totalCharacterCount += count

		characterArray = append(characterArray, documentCharacter)
	}

	if len(characterArray) > 0 {
		liberdatabase.AddDocumentCharacters(dbConn, characterArray)
		characterArray = []liberdatabase.DocumentCharacter{}
	}

	totalRuneCount := int64(0)
	runeArray := make([]liberdatabase.DocumentCharacter, 0, len(mapOfRunes))
	for character, count := range mapOfRunes {
		documentCharacter := liberdatabase.DocumentCharacter{
			FileId:         df.FileId,
			Character:      character,
			CharacterCount: count,
			CharacterType:  "RUNE",
		}

		totalRuneCount += count

		runeArray = append(runeArray, documentCharacter)
	}

	if len(characterArray) > 0 {
		liberdatabase.AddDocumentCharacters(dbConn, characterArray)
		characterArray = []liberdatabase.DocumentCharacter{}
	}

	if len(runeArray) > 0 {
		liberdatabase.AddDocumentCharacters(dbConn, runeArray)
		characterArray = []liberdatabase.DocumentCharacter{}
	}

	liberdatabase.UpdateTotalCharacterCount(dbConn, df.FileId, totalCharacterCount)
	liberdatabase.UpdateTotalRuneCount(dbConn, df.FileId, totalRuneCount)

	return nil
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

	totalWordCount := int64(0)
	tempWords := make([]liberdatabase.DocumentWord, 0, len(mapOfWords))
	for word, count := range mapOfWords {
		tempWord := liberdatabase.DocumentWord{
			FileId:    df.FileId,
			Word:      word,
			WordCount: count,
		}

		totalWordCount += count

		tempWords = append(tempWords, tempWord)

		if len(tempWords) >= 500 {
			liberdatabase.AddDocumentWord(dbConn, tempWords)
			tempWords = []liberdatabase.DocumentWord{}
		}
	}

	if len(tempWords) > 0 {
		liberdatabase.AddDocumentWord(dbConn, tempWords)
		tempWords = []liberdatabase.DocumentWord{}
	}

	liberdatabase.UpdateDocumentWordCount(dbConn, df.FileId, totalWordCount)

	calculateWordPercentages(dbConn, df.FileId, totalWordCount)

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

func calculateWordPercentages(dbConn *gorm.DB, fileId string, totalWordCount int64) {
	words := liberdatabase.GetDistinctWords(dbConn, fileId)
	percentages := make([]liberdatabase.DocumentWordStatistics, 0, len(words))

	for _, word := range words {
		wordPercent := (float64(word.WordCount) / float64(totalWordCount)) * 100

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

	if len(percentages) > 0 {
		liberdatabase.AddDocumentWordStatistics(dbConn, percentages)
		percentages = []liberdatabase.DocumentWordStatistics{}
	}
}
