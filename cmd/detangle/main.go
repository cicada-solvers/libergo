package main

import (
	runelib "characterrepo"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
	"runer"
	"strings"
)

type WordStruct struct {
	ID     string `gorm:"column:id"`
	Word   string `gorm:"column:word"`
	WordId string `gorm:"column:wordid"`
}

type WordEntry struct {
	gorm.Model
	Word                string `gorm:"column:word"`
	WordLength          int    `gorm:"column:word_length"`
	RuneglishWord       string `gorm:"column:runeglish_word"`
	RuneglishWordLength int    `gorm:"column:runeglish_word_length"`
	RuneWord            string `gorm:"column:rune_word"`
	RuneWordLength      int    `gorm:"column:rune_word_length"`
}

var repo = runelib.NewCharacterRepo()
var alphabetArray []string
var output string
var csvFiles []string
var dbConn *gorm.DB

func main() {
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

		for _, file := range csvFiles {
			fmt.Printf("Loading file: %s\n", file)
			loadErr := LoadFilesToDb(file)
			if loadErr != nil {
				fmt.Printf("Error loading CSV file: %v\n", loadErr)
			}
		}
	} else {
		InitSQLiteConnection()
	}

	// Process the text
	sentence := ""
	clonedText := CloneArray(textCharacters)
	ProcessText(countPatterns, clonedText, 0, sentence)
}

func ProcessText(countPatterns []int, textCharacters []string, position int, sentence string) {
	levelPrefix := fmt.Sprintf("Processing level: %d/%d", position+1, len(countPatterns))
	fmt.Printf("%s\n", levelPrefix)
	myCharacters := CloneArray(textCharacters)
	uuidString := uuid.New().String()

	if (position) < len(countPatterns)-1 {
		// First, we need to get the length of the current pattern
		patternLength := countPatterns[position]

		// Get the words from the CSV file
		readCount, loadErr := LoadWordsFromTable(uuidString, patternLength)
		if loadErr != nil {
			fmt.Printf("[%s] Error loading words from table: %v\n", levelPrefix, loadErr)
		}
		fmt.Printf("[%s] Loaded %d words from table\n", levelPrefix, readCount)

		isComplete := false
		for isComplete == false {
			word, getErr := GetFirstWordFromDatabase(uuidString)
			if getErr != nil {
				fmt.Printf("[%s] Error getting word from database: %v\n", levelPrefix, getErr)
				isComplete = true
			}
			if word.Word == "" {
				isComplete = true
			}

			if isComplete {
				break
			}

			fmt.Printf("[%s] Processing word: %s\n", levelPrefix, word.Word)
			newArrayWithRemoved, removedCount := RemoveLettersFromArray(myCharacters, strings.Split(word.Word, ""))
			fmt.Printf("[%s] Removed count: %d\n", levelPrefix, removedCount)
			fmt.Printf("[%s] New array size: %d\n", levelPrefix, len(newArrayWithRemoved))
			if removedCount == len(strings.Split(word.Word, "")) {
				ProcessText(countPatterns, newArrayWithRemoved, position+1, sentence)
				sentence += word.Word + "•"
			} else {
				fmt.Printf("[%s] Word (%s) skipped, not enough of the right letters.\n", levelPrefix, word.Word)
			}

			myCharacters = CloneArray(textCharacters)
			DeleteWordFromDatabase(word.ID)
		}
	} else {
		// First, we need to get the length of the current pattern
		patternLength := countPatterns[position]

		readCount, loadErr := LoadWordsFromTable(uuidString, patternLength)
		if loadErr != nil {
			fmt.Printf("[%s] Error loading words from CSV: %v\n", levelPrefix, loadErr)
		}
		fmt.Printf("[%s] Loaded %d words from table\n", levelPrefix, readCount)

		isComplete := false
		for isComplete == false {
			word, getErr := GetFirstWordFromDatabase(uuidString)
			if getErr != nil {
				fmt.Printf("[%s] Error getting word from database: %v\n", levelPrefix, getErr)
				isComplete = true
			}
			if word.Word == "" {
				isComplete = true
			}

			if isComplete {
				break
			}

			_, removedCount := RemoveLettersFromArray(myCharacters, strings.Split(word.Word, ""))
			if removedCount == len(strings.Split(word.Word, "")) {
				sentence += word.Word + "•"
				WriteToFile(sentence)
			}

			myCharacters = CloneArray(textCharacters)
			DeleteWordFromDatabase(uuidString)
		}
	}
}

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

func WriteToFile(text string) {
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
}

func CloneArray(array []string) []string {
	var retval []string
	for _, value := range array {
		retval = append(retval, value)
	}
	return retval
}

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

func IsLetterInAlphabet(character string) bool {
	for _, char := range alphabetArray {
		if char == character {
			return true
		}
	}
	return false
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

// LoadWordsFromTable reads all the words from a specific column in a CSV file.
func LoadWordsFromTable(guid string, length int) (int64, error) {
	fmt.Printf("Loading words from with a length of %d\n", length)
	readCount := int64(0)

	rows, _ := GetWordsFromDatabaseByLength(length)
	var words []WordStruct

	// Extract words from the specified column
	for _, row := range rows {
		instance := WordStruct{
			ID:     uuid.New().String(),
			Word:   row.RuneWord,
			WordId: guid,
		}

		words = append(words, instance)

		if len(words) >= 250 {
			AddWordToDatabase(words)
			words = []WordStruct{}
		}

		readCount++
	}

	AddWordToDatabase(words)

	return readCount, nil
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

func AddWordEntryToDatabase(words []WordEntry) {
	result := dbConn.CreateInBatches(&words, len(words))
	if result.Error != nil {
		fmt.Printf("error inserting word: %v\n", result.Error)
	}
}

func AddWordToDatabase(words []WordStruct) {
	result := dbConn.CreateInBatches(&words, len(words))
	if result.Error != nil {
		fmt.Printf("error inserting word: %v\n", result.Error)
	}
}

func GetFirstWordFromDatabase(uuid string) (WordStruct, error) {
	var word WordStruct
	result := dbConn.Where("wordid = ?", uuid).First(&word)
	if result.Error != nil {
		fmt.Printf("error querying words: %v\n", result.Error)
		return word, result.Error
	}

	return word, nil
}

func GetWordsFromDatabaseByLength(length int) ([]WordEntry, error) {
	var words []WordEntry
	result := dbConn.Where("rune_word_length = ?", length).Find(&words)
	if result.Error != nil {
		fmt.Printf("error querying words: %v\n", result.Error)
		return words, result.Error
	}
	return words, nil
}

func DeleteWordFromDatabase(uuid string) {
	result := dbConn.Exec(fmt.Sprintf("DELETE FROM `word_structs` WHERE id = '%s'", uuid))
	if result.Error != nil {
		fmt.Printf("error deleting word: %v\n", result.Error)
	}
}
