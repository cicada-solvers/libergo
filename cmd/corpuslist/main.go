package main

import (
	"bufio"
	"flag"
	"fmt"
	"lgstructs"
	"math/big"
	"os"
	"path/filepath"
	"runer"
	"sequences"
	"strings"
	"unicode/utf8"
)

// main reads text files from a directory and creates a SQL file with the dictionary words
func main() {
	// Define flags for the directory path and output SQL file
	dir := flag.String("dir", "./path/to/directory", "Directory to scan for text files")
	flag.Parse()

	wordList := make(map[string]bool)
	dictList := make([]lgstructs.DictionaryWord, 0, 16384)
	outputBase := "dictionary_words"
	outputCounter := 1

	// Scan the directory for text files
	err := filepath.WalkDir(*dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".txt") {
			processFile(path, wordList)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Failed to read directory: %v\n", err)
		return
	}

	// Read every word from the word list
	for word := range wordList {
		word = strings.ToUpper(word)
		runeglish := runer.PrepLatinToRune(word)
		runeText := runer.TransposeLatinToRune(runeglish)
		runeTextNoDoublet := lgstructs.RemoveDoublets(runeText)
		gemSum := runer.CalculateGemSum(runeText, runer.Runes)
		gemProd := runer.CalculateGemProduct(runeText, runer.Runes)

		dictWord := lgstructs.DictionaryWord{
			DictionaryWordText:          word,
			RuneglishWordText:           runeglish,
			RuneWordText:                runeText,
			RuneWordTextNoDoublet:       runeTextNoDoublet,
			GemSum:                      gemSum,
			GemSumPrime:                 sequences.IsPrime(big.NewInt(gemSum)),
			GemProduct:                  gemProd.String(),
			GemProductPrime:             sequences.IsPrime(&gemProd),
			DictionaryWordLength:        len(word),
			RuneglishWordLength:         len(runeglish),
			RuneWordLength:              utf8.RuneCountInString(runeText),
			RunePattern:                 lgstructs.GetRunePattern(word),
			RunePatternNoDoubletPattern: lgstructs.GetRunePattern(runeTextNoDoublet),
			Language:                    "English",
		}

		dictList = append(dictList, dictWord)

		if len(dictList) >= 250000 {
			outputFile := fmt.Sprintf("%s_%05d.sql", outputBase, outputCounter)
			writeSqlFile(outputFile, dictList)
			outputCounter++
			dictList = dictList[:0]
		}
	}

	outputFile := fmt.Sprintf("%s_%05d.sql", outputBase, outputCounter)
	writeSqlFile(outputFile, dictList)
}

// writeSqlFile writes the dictionary words to a SQL file
func writeSqlFile(outputFile string, dictList []lgstructs.DictionaryWord) {
	// Create the SQL file
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Failed to create output file: %v\n", err)
		return
	}
	defer func(file *os.File) {
		fileError := file.Close()
		if fileError != nil {
			fmt.Printf("Failed to close output file: %v\n", fileError)
		}
	}(file)

	writer := bufio.NewWriter(file)

	// Write the CREATE TABLE statement
	createTableSQL := `
CREATE TABLE IF NOT EXISTS dictionary_words (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    dict_word VARCHAR(255) COLLATE utf8mb4_general_ci NOT NULL,
    dict_runeglish VARCHAR(255) COLLATE utf8mb4_general_ci NOT NULL,
    dict_rune VARCHAR(255) COLLATE utf8mb4_general_ci NOT NULL,
    dict_rune_no_doublet VARCHAR(255) COLLATE utf8mb4_general_ci NOT NULL,
    gem_sum BIGINT NOT NULL,
    gem_sum_prime TINYINT(1) NOT NULL,
    gem_product VARCHAR(2048) COLLATE utf8mb4_general_ci NOT NULL,
    gem_product_prime TINYINT(1) NOT NULL,
    dict_word_length INT NOT NULL,
    dict_runeglish_length INT NOT NULL,
    dict_rune_length INT NOT NULL,
    rune_pattern VARCHAR(255) COLLATE utf8mb4_general_ci NOT NULL,
    rune_pattern_no_doublet VARCHAR(255) COLLATE utf8mb4_general_ci NOT NULL,
    language VARCHAR(255) COLLATE utf8mb4_general_ci NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
`
	_, err = writer.WriteString(createTableSQL + "\n")
	if err != nil {
		fmt.Printf("Failed to write to output file: %v\n", err)
		return
	}

	// Write the INSERT statements
	for _, dictWord := range dictList {
		insertSQL := fmt.Sprintf(`
INSERT INTO dictionary_words (
    dict_word, dict_runeglish, dict_rune, dict_rune_no_doublet, gem_sum, gem_sum_prime, gem_product, gem_product_prime, dict_word_length, dict_runeglish_length, dict_rune_length, rune_pattern, rune_pattern_no_doublet, language
) VALUES ('%s', '%s', '%s', '%s', %d, %t, '%s', %t, %d, %d, %d, '%s', '%s', '%s');
`,
			escapeString(dictWord.DictionaryWordText),
			escapeString(dictWord.RuneglishWordText),
			escapeString(dictWord.RuneWordText),
			escapeString(dictWord.RuneWordTextNoDoublet),
			dictWord.GemSum,
			dictWord.GemSumPrime,
			escapeString(dictWord.GemProduct),
			dictWord.GemProductPrime,
			dictWord.DictionaryWordLength,
			dictWord.RuneglishWordLength,
			dictWord.RuneWordLength,
			escapeString(dictWord.RunePattern),
			escapeString(dictWord.RunePatternNoDoubletPattern),
			escapeString(dictWord.Language),
		)
		_, err = writer.WriteString(insertSQL + "\n")
		if err != nil {
			fmt.Printf("Failed to write to output file: %v\n", err)
			return
		}
	}

	// Flush the writer
	err = writer.Flush()
	if err != nil {
		fmt.Printf("Failed to flush output file: %v\n", err)
		return
	}

	fmt.Printf("SQL file created successfully: %s\n", outputFile)
}

// escapeString escapes single quotes in a string for SQL insertion
func escapeString(str string) string {
	return strings.ReplaceAll(str, "'", "''")
}

// processFile reads a file and adds all words to the word list
func processFile(filePath string, wordList map[string]bool) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		return
	}
	defer func(file *os.File) {
		closeError := file.Close()
		if closeError != nil {
			fmt.Printf("Failed to close file: %v\n", closeError)
		}
	}(file)

	reader := bufio.NewReader(file)
	for {
		line, readError := reader.ReadString('\n')
		if readError != nil {
			if readError.Error() != "EOF" {
				fmt.Printf("Error reading file: %v\n", err)
			}
			break
		}
		words := strings.Fields(line)
		for _, word := range words {
			wordList[word] = true
		}
	}
}
