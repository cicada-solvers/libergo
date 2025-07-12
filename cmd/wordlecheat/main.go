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

// WordEntry represents a word and its length, stored in a database using GORM for ORM.
type WordEntry struct {
	gorm.Model
	Word       string
	WordLength int `gorm:"index:idx_word_length"`
}

// csvFiles holds the list of CSV file names to be processed or loaded into the database.
var csvFiles []string

// dbConn is a global variable representing the database connection instance, managed using the GORM library.
var dbConn *gorm.DB

// dbMutex is used to synchronize access to shared database resources, preventing race conditions during operations.
var dbMutex sync.Mutex

// main is the entry point of the program, handling command-line arguments and executing logic based on the specified mode.
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

// DetectAndLoadTheFiles initializes the database, removes existing data if present, and loads CSV file data into the database.
// It concurrently processes each CSV file, ensuring all files are loaded before completing.
// Returns true if the operation succeeds, false if any errors occur during the process.
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

// GetCSVFiles retrieves a slice of CSV file names from the current working directory. Logs a fatal error if directory read fails.
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

// InitSQLiteConnection initializes the SQLite database connection and sets the global dbConn variable using GORM.
func InitSQLiteConnection() {
	foldrPath := "."

	databasePath := filepath.Join(foldrPath, "/words.db")

	dbConn, _ = gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
}

// InitSQLiteTables initializes the database schema by creating tables using GORM's AutoMigrate function. Returns an error if schema creation fails.
func InitSQLiteTables() error {
	// Migrate the schemas
	dbCreateError := dbConn.AutoMigrate(&WordEntry{})
	if dbCreateError != nil {
		fmt.Printf("Error creating table: %v\n", dbCreateError)
		return dbCreateError
	}

	return nil
}

// LoadFilesToDb reads a CSV file, processes its contents, and adds WordEntry records to the database in batches.
// Accepts the filePath of the CSV as input and returns an error if the operation fails.
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

// AddWordEntryToDatabase adds a batch of WordEntry records to the database, ensuring thread-safe operation using dbMutex.
// It creates the entries using GORM's CreateInBatches method and logs errors if the insertion fails.
func AddWordEntryToDatabase(words []WordEntry) {
	dbMutex.Lock()
	result := dbConn.CreateInBatches(&words, len(words))
	if result.Error != nil {
		fmt.Printf("error inserting word: %v\n", result.Error)
	}
	dbMutex.Unlock()
}

// GetWordsByLength retrieves words from the database with the specified length, sorts them alphabetically, and ensures uniqueness.
// Returns the sorted list of WordEntry objects or an error if the query fails.
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

// SortWordsAlphabetically sorts a slice of WordEntry objects alphabetically based on the Word field and returns the sorted slice.
func SortWordsAlphabetically(words []WordEntry) []WordEntry {
	sort.Slice(words, func(i, j int) bool {
		return words[i].Word < words[j].Word
	})
	return words
}

// UniqueSort removes duplicate WordEntry objects from the input slice and returns a slice with unique entries.
func UniqueSort(words []WordEntry) []WordEntry {
	var uniqueWords []WordEntry
	for _, word := range words {
		if !slices.Contains(uniqueWords, word) {
			uniqueWords = append(uniqueWords, word)
		}
	}
	return uniqueWords
}

// ParseWords filters a list of WordEntry objects based on inclusion, exclusion, and positional criteria in the input parameters.
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
