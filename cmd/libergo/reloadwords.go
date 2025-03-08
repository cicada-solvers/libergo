package main

import (
	"bufio"
	"config"
	"fmt"
	"gorm.io/gorm"
	"liberdatabase"
	"math/big"
	"os"
	"path/filepath"
	"runer"
	"sequences"
	"strings"
	"sync"
	"unicode/utf8"
)

func ReloadWords() {
	// Get the configuration folder path
	configFolderPath, err := config.GetConfigFolderPath()
	if err != nil {
		_, err2 := fmt.Fprintf(os.Stderr, "Error getting config folder path: %v\n", err)
		if err2 != nil {
			return
		}
		os.Exit(1)
	}

	// Construct the path to the words.txt file
	wordsFilePath := filepath.Join(configFolderPath, "words.txt")

	// Open the words.txt file
	file, err := os.Open(wordsFilePath)
	if err != nil {
		_, err3 := fmt.Fprintf(os.Stderr, "Error opening words.txt file: %v\n", err)
		if err3 != nil {
			return
		}
		os.Exit(1)
	}
	defer func(file *os.File) {
		fileErr := file.Close()
		if fileErr != nil {
			_, err4 := fmt.Fprintf(os.Stderr, "Error closing words.txt file: %v\n", err)
			if err4 != nil {
				return
			}
			os.Exit(1)
		}
	}(file)

	// Read all lines into memory
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		_, err5 := fmt.Fprintf(os.Stderr, "Error reading words.txt file: %v\n", err)
		if err5 != nil {
			return
		}
		os.Exit(1)
	}

	// Get the number of workers from the configuration
	configuration, _ := config.LoadConfig()
	numWorkers := configuration.NumWorkers

	// Create a wait group and channels for the words and errors
	var wg sync.WaitGroup
	wordChan := make(chan string, len(lines))
	errChan := make(chan error, len(lines))

	// Start a fixed number of worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Get the database connection for this worker
			db, err := liberdatabase.InitConnection()
			if err != nil {
				errChan <- fmt.Errorf("error initializing database connection: %v", err)
				return
			}
			defer func(db *gorm.DB) {
				closeErr := liberdatabase.CloseConnection(db)
				if closeErr != nil {
					errChan <- fmt.Errorf("error closing database connection: %v", err)
				}
			}(db)

			for latinText := range wordChan {
				runeglishText := runer.PrepLatinToRune(strings.ToUpper(latinText))
				runeText := runer.TransposeLatinToRune(runeglishText)
				gemSum := runer.CalculateGemSum(runeText, runer.Runes)
				gemSumBigInt := big.NewInt(gemSum)
				gemSumPrime := sequences.IsPrime(gemSumBigInt)
				gemProduct := runer.CalculateGemProduct(runeText, runer.Runes)
				gemProductPrime := sequences.IsPrime(&gemProduct)

				fmt.Printf("Loading Word: %s Runeglish: %s Runes: %s Gem Sum: %d\n", strings.ToUpper(latinText), runeglishText, runeText, gemSum)

				dictionaryWord := liberdatabase.DictionaryWord{
					DictionaryWordText:   strings.ToUpper(latinText),
					RuneglishWordText:    runeglishText,
					RuneWordText:         runeText,
					GemSum:               gemSum,
					GemSumPrime:          gemSumPrime,
					GemProduct:           gemProduct.String(),
					GemProductPrime:      gemProductPrime,
					DictionaryWordLength: len(latinText),
					RuneglishWordLength:  len(runeglishText),
					RuneWordLength:       utf8.RuneCountInString(runeText),
					Language:             "English",
				}

				dictionaryWord.RunePattern = dictionaryWord.GetRunePattern()
				dictionaryWord.RuneWordTextNoDoublet = dictionaryWord.GetRuneWordToRuneWordNoDoublet()
				dictionaryWord.RunePatternNoDoubletPattern = dictionaryWord.GetRunePatternNoDoublet()

				err = liberdatabase.InsertDictionaryWord(db, dictionaryWord)
				if err != nil {
					errChan <- fmt.Errorf("error inserting dictionary word: %v", err)
				}
			}
		}()
	}

	// Send the words to the channel
	for _, line := range lines {
		wordChan <- line
	}
	close(wordChan)

	// Wait for all goroutines to finish
	wg.Wait()
	close(errChan)

	// Handle errors
	for err := range errChan {
		_, charErr := fmt.Fprintf(os.Stderr, "%v\n", err)
		if charErr != nil {
			return
		}
	}
}

func GenerateSQLScript(words []liberdatabase.DictionaryWord, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating SQL script file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("error closing SQL script file: %v\n", err)
		}
	}(file)

	createTableSQL := `
CREATE TABLE IF NOT EXISTS dictionary_words (
 id BIGINT AUTO_INCREMENT PRIMARY KEY,
 dict_word VARCHAR(255) NOT NULL,
 dict_runeglish VARCHAR(255) NOT NULL,
 dict_rune VARCHAR(255) NOT NULL,
 dict_rune_no_doublet VARCHAR(255) NOT NULL,
 gem_sum BIGINT NOT NULL,
 gem_sum_prime BOOLEAN NOT NULL,
 gem_product VARCHAR(2048) NOT NULL,
 gem_product_prime BOOLEAN NOT NULL,
 dict_word_length INT NOT NULL,
 dict_runeglish_length INT NOT NULL,
 dict_rune_length INT NOT NULL,
 rune_pattern VARCHAR(255) NOT NULL,
 rune_pattern_no_doublet VARCHAR(255) NOT NULL,
 language VARCHAR(255) NOT NULL
);
`
	_, err = file.WriteString(createTableSQL)
	if err != nil {
		return fmt.Errorf("error writing create table statement to SQL script file: %v", err)
	}

	for _, word := range words {
		sql := fmt.Sprintf(
			"INSERT INTO dictionary_words (dict_word, dict_runeglish, dict_rune, dict_rune_no_doublet, gem_sum, gem_sum_prime, gem_product, gem_product_prime, dict_word_length, dict_runeglish_length, dict_rune_length, rune_pattern, rune_pattern_no_doublet, language) VALUES ('%s', '%s', '%s', '%s', %d, %t, '%s', %t, %d, %d, %d, '%s', '%s', '%s');\n",
			strings.ReplaceAll(word.DictionaryWordText, "'", "''"),
			strings.ReplaceAll(word.RuneglishWordText, "'", "''"),
			strings.ReplaceAll(word.RuneWordText, "'", "''"),
			strings.ReplaceAll(word.RuneWordTextNoDoublet, "'", "''"),
			word.GemSum,
			word.GemSumPrime,
			strings.ReplaceAll(word.GemProduct, "'", "''"),
			word.GemProductPrime,
			word.DictionaryWordLength,
			word.RuneglishWordLength,
			word.RuneWordLength,
			strings.ReplaceAll(word.RunePattern, "'", "''"),
			strings.ReplaceAll(word.RunePatternNoDoubletPattern, "'", "''"),
			strings.ReplaceAll(word.Language, "'", "''"),
		)
		_, err := file.WriteString(sql)
		if err != nil {
			return fmt.Errorf("error writing to SQL script file: %v", err)
		}
	}

	return nil
}

func GenerateMySQLScript(fileTypeInfoModels []liberdatabase.FileTypeInfoModel) error {
	file, err := os.Create("file_type_info_mysql.sql")
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	writer := bufio.NewWriter(file)
	defer func(writer *bufio.Writer) {
		err := writer.Flush()
		if err != nil {
			fmt.Println(err)
		}
	}(writer)

	_, err = writer.WriteString("CREATE TABLE IF NOT EXISTS file_type_info (\n" +
		"  id BIGINT AUTO_INCREMENT PRIMARY KEY,\n" +
		"  name VARCHAR(255),\n" +
		"  file_type VARCHAR(255),\n" +
		"  mime_type VARCHAR(255),\n" +
		"  header TEXT,\n" +
		"  alias TEXT,\n" +
		"  offset INT,\n" +
		"  sub_header TEXT\n" +
		");\n\n")
	if err != nil {
		return err
	}

	for _, model := range fileTypeInfoModels {
		_, err = writer.WriteString(fmt.Sprintf("INSERT INTO file_type_info (name, file_type, mime_type, header, alias, offset, sub_header) VALUES ('%s', '%s', '%s', '%s', '%s', %d, x'%s');\n",
			model.Name, model.FileType, model.MimeType, model.Header, model.Alias, model.Offset, model.SubHeader))
		if err != nil {
			return err
		}
	}

	return nil
}
