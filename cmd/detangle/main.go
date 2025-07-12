package main

import (
	runelib "characterrepo"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
	"runer"
	"runtime"
	"slices"
	"strings"
	"sync"
)

type WordStruct struct {
	ProcessId string `gorm:"index:idx_first_word"`
	Word      string
	Sequence  int64 `gorm:"index:idx_first_word"`
	Level     int   `gorm:"index:idx_first_word"`
}

type WordEntry struct {
	gorm.Model
	Word                string
	WordLength          int `gorm:"index:idx_word_length"`
	RuneglishWord       string
	RuneglishWordLength int `gorm:"index:idx_runeglish_length"`
	RuneWord            string
	RuneWordLength      int `gorm:"index:idx_rune_length"`
}

type PatternPosition struct {
	Pattern  int
	Position int
}

var repo = runelib.NewCharacterRepo()
var alphabetArray []string
var output string
var csvFiles []string
var dbConn *gorm.DB
var processId string
var dbMutex sync.Mutex

// main is the entry point of the program. It initializes flags, processes input, and handles text processing tasks.
func main() {
	processId = uuid.New().String()
	text := flag.String("text", "", "The text to decode")
	alphabet := flag.String("alphabet", "rune", "The alphabet to use (rune or english)")
	outputFile := flag.String("output", "", "The output file to write the results")

	// Parse the flags
	flag.Parse()

	// Validate required flags
	if *text == "" {
		log.Fatal("The -text flag is required")
	}
	if *outputFile == "" {
		log.Fatal("The -output flag is required")
	}

	output = *outputFile

	// Get the CSV files
	csvFiles = GetCSVFiles()
	fmt.Printf("CSV files: %v\n", csvFiles)

	// Get the alphabet
	if *alphabet == "rune" {
		alphabetArray = repo.GetGematriaRunes()
	} else {
		alphabetArray = repo.GetRunglishAlphabet()
	}

	// Get what characters are available based on the alphabet
	textCharacters := GetAvailableCharacters(*text)
	fmt.Printf("Text characters: %v\n", textCharacters)
	countPatterns := GetCountPatterns(*text)
	fmt.Printf("Count patterns: %v\n", countPatterns)

	// Load the words to the database
	// Only load if the db does not exist
	_, err := os.Stat("words.db")
	if os.IsNotExist(err) {
		fmt.Printf("DB does not exist, loading words\n")
		InitSQLiteConnection()
		initErr := InitSQLiteTables()
		if initErr != nil {
			fmt.Printf("Error initializing SQLite tables: %v\n", initErr)
			return
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
	} else {
		InitSQLiteConnection()
		VacuumDb()
	}

	// Load the words from the database
	positionChan := make(chan PatternPosition)
	var pwg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		pwg.Add(1)
		go func() {
			defer pwg.Done()
			for positionPattern := range positionChan {
				fmt.Printf("[%d/%d] loading words from table\n", positionPattern.Position, len(countPatterns))
				readCount, loadErr := LoadWordsFromTable(positionPattern.Pattern, positionPattern.Position)
				if loadErr != nil {
					fmt.Printf("[%d/%d] Error loading words from table: %v\n", positionPattern.Position, len(countPatterns), loadErr)
				}
				fmt.Printf("[%d/%d] Loaded %d words from table\n", positionPattern.Position, len(countPatterns), readCount)
			}
		}()
	}

	for position, pattern := range countPatterns {
		positionChan <- PatternPosition{
			Pattern:  pattern,
			Position: position,
		}
	}
	close(positionChan)
	pwg.Wait()

	// Process the text
	sentence := ""
	clonedText := CloneArray(textCharacters)
	ProcessText(countPatterns, clonedText, 0, sentence)
}

// ProcessText recursively processes text based on patterns, extracting words from a database and structuring sentences.
// It iterates over text patterns and database sequences, modifying the input text and appending results to sentences.
// Handles remaining letters and resets text state if conditions require further recursive processing.
func ProcessText(countPatterns []int, textCharacters []string, position int, passedSentence string) {
	sequence, _ := GetMaxSequenceFromDatabaseByLevel(position)
	levelPrefix := fmt.Sprintf("Processing level: %d/%d - %d", position+1, len(countPatterns), sequence)
	fmt.Printf("%s\n", levelPrefix)
	myCharactersLeft := CloneArray(textCharacters)
	sentence := CloneSentence(passedSentence)

	if (position) < len(countPatterns)-1 {
		// Get the words from the CSV file
		for sequence >= 1 {
			levelPrefix = fmt.Sprintf("Processing level: %d/%d - %d", position+1, len(countPatterns), sequence)
			fmt.Printf("[%s] Sequence %d\n", levelPrefix, sequence)
			word, getErr := GetFirstWordFromDatabase(&sequence, position, myCharactersLeft)
			if getErr != nil {
				fmt.Printf("[%s] Error getting word from database: %v\n", levelPrefix, getErr)
				continue
			}

			fmt.Printf("[%s] Processing word: %s\n", levelPrefix, word.Word)
			newArrayWithRemoved, removedCount := RemoveLettersFromArray(myCharactersLeft, strings.Split(word.Word, ""))
			fmt.Printf("[%s] Removed count: %d\n", levelPrefix, removedCount)
			fmt.Printf("[%s] New array size: %d\n", levelPrefix, len(newArrayWithRemoved))
			sentence += word.Word + "•"
			ProcessText(countPatterns, newArrayWithRemoved, position+1, sentence)

			myCharactersLeft = CloneArray(textCharacters)
			sentence = CloneSentence(passedSentence)
			sequence--
		}
	} else {
		for sequence >= 1 {
			levelPrefix = fmt.Sprintf("Processing level: %d/%d - %d", position+1, len(countPatterns), sequence)
			fmt.Printf("[%s] Sequence %d\n", levelPrefix, sequence)
			word, getErr := GetFirstWordFromDatabase(&sequence, position, myCharactersLeft)
			if getErr != nil {
				fmt.Printf("[%s] Error getting word from database: %v\n", levelPrefix, getErr)
				continue
			}

			RemoveLettersFromArray(myCharactersLeft, strings.Split(word.Word, ""))
			sentence += word.Word + "•"
			WriteToFile(sentence)

			myCharactersLeft = CloneArray(textCharacters)
			sentence = CloneSentence(passedSentence)
			sequence--
		}
	}
}

// RemoveLettersFromArray removes specified letters from the input array and returns the updated array and count of removed items.
func RemoveLettersFromArray(array []string, letters []string) ([]string, int) {
	result := make([]string, len(array))
	copy(result, array)
	removedCount := 0

	for _, letter := range letters {
		for i, val := range result {
			if val == letter {
				// Remove this element by shifting everything after it
				result = append(result[:i], result[i+1:]...)
				removedCount++
				break
			}
		}
	}

	return result, removedCount
}

// WriteToFile writes the provided text to a file, appending a Latin-transposed version of the text and a newline.
func WriteToFile(text string) {
	dbMutex.Lock()
	file, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}(file)

	_, err = file.WriteString(text)
	_, err = file.WriteString("\n")
	translated := runer.TransposeRuneToLatin(text)
	_, err = file.WriteString("\n")
	_, err = file.WriteString(translated)
	_, err = file.WriteString("\n")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Wrote to file: %s\n", output)
	dbMutex.Unlock()
}

// CloneArray creates a new copy of the input string slice and returns it without modifying the original slice.
func CloneArray(array []string) []string {
	var retval []string
	for _, value := range array {
		retval = append(retval, value)
	}
	return retval
}

// CloneSentence creates a clone of the given sentence by splitting it into characters and rejoining them into a new string.
func CloneSentence(sentence string) string {
	var tmpval []string

	for _, value := range strings.Split(sentence, "") {
		tmpval = append(tmpval, value)
	}

	return strings.Join(tmpval, "")
}

// GetCountPatterns processes a string and returns a slice of integers representing counts of valid sequential characters.
func GetCountPatterns(text string) []int {
	var retval []int
	textArray := strings.Split(text, "")
	counter := 0

	for _, character := range textArray {
		if IsLetterInAlphabet(character) && !repo.IsDinkus(character) && !repo.IsSeperator(character) {
			counter++
		} else {
			retval = append(retval, counter)
			counter = 0
		}
	}

	return retval
}

// GetAvailableCharacters filters and returns a list of characters from the input text that exist in the predefined alphabet.
func GetAvailableCharacters(text string) []string {
	var retval []string
	textArray := strings.Split(text, "")
	for _, character := range textArray {
		if IsLetterInAlphabet(character) {
			retval = append(retval, character)
		}
	}

	return retval
}

// IsLetterInAlphabet checks if the given character exists in the predefined alphabet array and returns true if found.
func IsLetterInAlphabet(character string) bool {
	for _, char := range alphabetArray {
		if char == character {
			return true
		}
	}
	return false
}

// GetCSVFiles retrieves the list of CSV file names from the current directory.
// The function scans the current directory and filters files ending with ".csv".
// It returns a slice containing names of all found CSV files.
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

// LoadWordsFromTable reads all the words from a specific column in a CSV file.
func LoadWordsFromTable(length, level int) (int64, error) {
	fmt.Printf("Loading words from with a length of %d\n", length)
	readCount := int64(0)
	sequence := int64(1)

	rows, _ := GetWordsFromDatabaseByLength(length)
	var words []WordStruct

	// Extract words from the specified column
	for _, row := range rows {
		instance := WordStruct{
			ProcessId: processId,
			Word:      row.RuneWord,
			Sequence:  sequence,
			Level:     level,
		}

		words = append(words, instance)

		if len(words) >= 250 {
			AddWordToDatabase(words)
			words = []WordStruct{}
		}

		readCount++
		sequence++
	}

	AddWordToDatabase(words)

	return readCount, nil
}

// LoadFilesToDb reads words from a CSV file and inserts them into the database in batches.
// Takes filePath as the path to the input CSV file.
// Returns an error if file operations, CSV parsing, or database operations fail.
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
			Word:                row[0],
			RuneglishWord:       row[1],
			RuneWord:            row[2],
			WordLength:          len(strings.Split(row[0], "")),
			RuneglishWordLength: len(strings.Split(row[1], "")),
			RuneWordLength:      len(strings.Split(row[2], "")),
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

// InitSQLiteConnection initializes the SQLite database
func InitSQLiteConnection() {
	foldrPath := "."

	databasePath := filepath.Join(foldrPath, "/words.db")

	dbConn, _ = gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
}

// InitSQLiteTables initializes and sets up the necessary SQLite database tables.
// It drops pre-existing tables if they exist and creates new ones based on the defined schemas.
// Returns an error if there is any issue during the table creation or migration process.
func InitSQLiteTables() error {
	// Remove the old table if it exists
	dropError := dbConn.Migrator().DropTable(&WordStruct{})
	if dropError != nil {
		fmt.Printf("Error dropping table: %v\n", dropError)
	}

	// Migrate the schemas
	dbCreateError := dbConn.AutoMigrate(&WordStruct{})
	if dbCreateError != nil {
		fmt.Printf("Error creating table: %v\n", dbCreateError)
	}
	dbCreateError = dbConn.AutoMigrate(&WordEntry{})
	if dbCreateError != nil {
		fmt.Printf("Error creating table: %v\n", dbCreateError)
	}

	return nil
}

// AddWordEntryToDatabase inserts a batch of WordEntry records into the database, ensuring thread safety with a mutex lock.
// Accepts a slice of WordEntry as input. Logs an error if database insertion fails.
func AddWordEntryToDatabase(words []WordEntry) {
	dbMutex.Lock()
	result := dbConn.CreateInBatches(&words, len(words))
	if result.Error != nil {
		fmt.Printf("error inserting word: %v\n", result.Error)
	}
	dbMutex.Unlock()
}

// AddWordToDatabase inserts a batch of words into the database and handles errors during the insertion process.
// It uses a mutex lock to ensure thread safety when accessing the database connection.
func AddWordToDatabase(words []WordStruct) {
	dbMutex.Lock()
	result := dbConn.CreateInBatches(&words, len(words))
	if result.Error != nil {
		fmt.Printf("error inserting word: %v\n", result.Error)
	}
	dbMutex.Unlock()
}

// GetFirstWordFromDatabase retrieves the first valid word from the database matching the provided sequence, level, and criteria.
// It locks database access, filters words against provided characters, and decrements the sequence if conditions aren't met.
func GetFirstWordFromDatabase(sequence *int64, level int, charactersLeft []string) (WordStruct, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	var word WordStruct

	for *sequence >= 1 {
		result := dbConn.Where("process_id = ? and sequence = ? and level = ?", processId, *sequence, level).First(&word)
		if result.Error != nil {
			fmt.Printf("error querying words: %v\n", result.Error)
			return word, result.Error
		}

		if !AreAllLettersInRemaining(word, charactersLeft) {
			*sequence--
		} else {
			// Return nil word and no error
			return word, nil
		}
	}

	// Return nil word and new error
	return word, errors.New("no words left in database")
}

// AreAllLettersInRemaining checks if all letters in the given word are present in the provided slice of remaining characters.
func AreAllLettersInRemaining(word WordStruct, charactersLeft []string) bool {
	lettersInWord := strings.Split(word.Word, "")
	for _, letter := range lettersInWord {
		if !slices.Contains(charactersLeft, letter) {
			return false
		}
	}

	return true
}

// GetWordsFromDatabaseByLength retrieves words from the database where the word length matches the specified value.
// Returns a slice of WordEntry and an error if the query fails.
// Thread-safe due to a mutex lock while accessing the database connection.
func GetWordsFromDatabaseByLength(length int) ([]WordEntry, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	var words []WordEntry
	result := dbConn.Where("rune_word_length = ?", length).Find(&words)
	if result.Error != nil {
		fmt.Printf("error querying words: %v\n", result.Error)
		return words, result.Error
	}

	return words, nil
}

// GetMaxSequenceFromDatabaseByLevel retrieves the maximum sequence value from the database for a specified level.
// It locks the database access, executes the query, and returns the maximum sequence value or an error if one occurs.
func GetMaxSequenceFromDatabaseByLevel(level int) (int64, error) {
	dbMutex.Lock()
	var maxSequence int64
	sql := fmt.Sprintf("SELECT MAX(sequence) FROM `word_structs` WHERE process_id = \"%s\" and level = %d", processId, level)
	dbConn.Raw(sql).Scan(&maxSequence)
	dbMutex.Unlock()
	return maxSequence, nil
}

// VacuumDb clears and reorganizes the SQLite database by truncating, vacuuming, and reindexing the `word_structs` table.
func VacuumDb() {
	dbMutex.Lock()
	result := dbConn.Exec("DELETE FROM word_structs")
	if result.Error != nil {
		fmt.Printf("error truncating: %v\n", result.Error)
	}

	result = dbConn.Exec("VACUUM main")
	if result.Error != nil {
		fmt.Printf("error vacuuming: %v\n", result.Error)
	}

	result = dbConn.Exec("REINDEX 'word_structs'")
	if result.Error != nil {
		fmt.Printf("error reindexing: %v\n", result.Error)
	}
	dbMutex.Unlock()
	return
}
