package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"sync"
)

type WordEntry struct {
	gorm.Model
	Word       string
	WordLength int `gorm:"index:idx_word_length"`
}

var csvFiles []string
var dbConn *gorm.DB
var dbMutex sync.Mutex

func main() {
	length := flag.Int("length", 5, "The length of the word")
	text := flag.String("text", "", "The text to decode")
	contains := flag.String("contains", "", "comma separated list of letters to check for")
	exclude := flag.String("exclude", "", "comma separated list of letters to exclude")
	mode := flag.String("mode", "solve", "load or solve")

	// Parse the flags
	flag.Parse()

	// remove spaces and uppercase flags
	*text = strings.ReplaceAll(*text, " ", "")
	*contains = strings.ReplaceAll(*contains, " ", "")
	*exclude = strings.ReplaceAll(*exclude, " ", "")
	*text = strings.ToUpper(*text)
	*contains = strings.ToUpper(*contains)
	*exclude = strings.ToUpper(*exclude)

	if *mode == "load" {
		if !DetectAndLoadTheFiles() {
			return
		}
	} else {
		_, err := os.Stat("words.db")
		if os.IsNotExist(err) {
			if !DetectAndLoadTheFiles() {
				return
			}
		} else {
			InitSQLiteConnection()
		}

		// This is for solving
		words, _ := GetWordsByLength(*length)
		words, _ = ParseWords(strings.Split(*text, ""), strings.Split(*contains, ","), strings.Split(*exclude, ","), words)

		fmt.Printf("Found %d words\n", len(words))
		for _, word := range words {
			fmt.Printf("%s\n", word.Word)
		}
	}
}

func DetectAndLoadTheFiles() bool {
	// Get the CSV files
	csvFiles = GetCSVFiles()
	fmt.Printf("CSV files: %v\n", csvFiles)

	_, err := os.Stat("words.db")
	if os.IsExist(err) {
		// remove words.db file
		removeErr := os.Remove("words.db")
		if removeErr != nil {
			fmt.Printf("Error removing words.db: %v\n", removeErr)
			return false
		}
	}
	fmt.Printf("DB does not exist, loading words\n")
	InitSQLiteConnection()
	initErr := InitSQLiteTables()
	if initErr != nil {
		fmt.Printf("Error initializing SQLite tables: %v\n", initErr)
		return false
	}

	// Load the CSV files to the database
	var lwg sync.WaitGroup
	for _, file := range csvFiles {
		lwg.Add(1)
		go func() {
			defer lwg.Done()
			fmt.Printf("Loading file: %s\n", file)
			loadErr := LoadFilesToDb(file)
			if loadErr != nil {
				fmt.Printf("Error loading CSV file: %v\n", loadErr)
			}
		}()
	}
	lwg.Wait()
	fmt.Printf("Finished loading files\n")
	return true
}

func GetCSVFiles() []string {
	// Get the CSV files from the current directory
	var fileNames []string
	files, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".csv") {
			fileNames = append(fileNames, file.Name())
		}
	}

	return fileNames
}

func InitSQLiteConnection() {
	foldrPath := "."

	databasePath := filepath.Join(foldrPath, "/words.db")

	dbConn, _ = gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
}

func InitSQLiteTables() error {
	// Migrate the schemas
	dbCreateError := dbConn.AutoMigrate(&WordEntry{})
	if dbCreateError != nil {
		fmt.Printf("Error creating table: %v\n", dbCreateError)
		return dbCreateError
	}

	return nil
}

func LoadFilesToDb(filePath string) error {
	var words []WordEntry
	fmt.Printf("Loading words to DB\n")

	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("failed to open file: %w", err)
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
		fmt.Printf("Failed to read CSV: %v\n", err)
	}

	// Extract words from the specified column
	for _, row := range rows {
		// Add the word to the database
		instance := WordEntry{
			Word:       row[0],
			WordLength: len(strings.Split(row[0], "")),
		}

		words = append(words, instance)

		if len(words) >= 250 {
			AddWordEntryToDatabase(words)
			words = []WordEntry{}
		}
	}

	AddWordEntryToDatabase(words)

	return nil
}

func AddWordEntryToDatabase(words []WordEntry) {
	dbMutex.Lock()
	result := dbConn.CreateInBatches(&words, len(words))
	if result.Error != nil {
		fmt.Printf("error inserting word: %v\n", result.Error)
	}
	dbMutex.Unlock()
}

func GetWordsByLength(length int) ([]WordEntry, error) {
	var words []WordEntry
	result := dbConn.Where("word_length = ?", length).Find(&words)
	if result.Error != nil {
		fmt.Printf("error querying words: %v\n", result.Error)
		return words, result.Error
	}
	words = SortWordsAlphabetically(words)
	words = UniqueSort(words)
	return words, nil
}

func SortWordsAlphabetically(words []WordEntry) []WordEntry {
	sort.Slice(words, func(i, j int) bool {
		return words[i].Word < words[j].Word
	})
	return words
}

func UniqueSort(words []WordEntry) []WordEntry {
	var uniqueWords []WordEntry
	for _, word := range words {
		if !slices.Contains(uniqueWords, word) {
			uniqueWords = append(uniqueWords, word)
		}
	}
	return uniqueWords
}

func ParseWords(text, contains, exclude []string, words []WordEntry) ([]WordEntry, error) {
	var preFilteredWords []WordEntry
	var filteredWords []WordEntry

	for _, word := range words {
		include := true

		if len(contains) > 0 {
			for _, letter := range contains {
				if !strings.Contains(word.Word, letter) {
					include = false
					break
				}
			}
		}

		if len(exclude) > 0 {
			for _, letter := range exclude {
				if strings.Contains(word.Word, letter) {
					include = false
					break
				}
			}
		}

		if include {
			preFilteredWords = append(preFilteredWords, word)
		}
	}

	if len(text) == 0 {
		return preFilteredWords, nil
	}

	for _, word := range preFilteredWords {
		include := true
		wordArray := strings.Split(word.Word, "")
		for position, letter := range text {
			if letter != "%" {
				if letter != wordArray[position] {
					include = false
				}
			}
		}

		if include {
			filteredWords = append(filteredWords, word)
		}
	}

	return filteredWords, nil
}
