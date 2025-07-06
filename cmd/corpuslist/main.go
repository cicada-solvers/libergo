package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"lgstructs"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"runer"
	"sequences"
	"strconv"
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
		word = filterNumbersOut(word)

		if len(word) <= 0 {
			continue
		}

		runeglish := runer.PrepLatinToRune(word)
		runeText := runer.TransposeLatinToRune(runeglish)
		runeTextNoDoublet := lgstructs.RemoveDoublets(strings.Split(runeText, ""))
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
			DictRuneNoDoubletLength:     utf8.RuneCountInString(runeTextNoDoublet),
			RuneWordLength:              utf8.RuneCountInString(runeText),
			RunePattern:                 lgstructs.GetRunePattern(word),
			RunePatternNoDoubletPattern: lgstructs.GetRunePattern(runeTextNoDoublet),
			RuneDistancePattern:         lgstructs.GetRuneDistancePattern(strings.Split(runeText, "")),
			Language:                    "English",
		}

		dictList = append(dictList, dictWord)
		delete(wordList, word) // Remove the word from wordList

		if len(dictList) >= math.MaxInt-1 {
			outputFile := fmt.Sprintf("%s_%05d.csv", outputBase, outputCounter)
			writeCsvFile(outputFile, dictList)
			outputCounter++
			dictList = dictList[:0]
		}
	}

	outputFile := fmt.Sprintf("%s_%05d.csv", outputBase, outputCounter)
	writeCsvFile(outputFile, dictList)
}

// writeCsvFile writes the dictionary words to a CSV file
func writeCsvFile(outputFile string, dictList []lgstructs.DictionaryWord) {
	// Create the CSV file
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

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header
	header := []string{
		"dict_word", "dict_runeglish", "dict_rune", "dict_rune_no_doublet",
		"gem_sum", "gem_sum_prime", "gem_product", "gem_product_prime",
		"dict_word_length", "dict_runeglish_length", "dict_rune_length",
		"rune_pattern", "rune_pattern_no_doublet", "dict_rune_no_doublet_length",
		"rune_distance_pattern", "language",
	}
	if err := writer.Write(header); err != nil {
		fmt.Printf("Failed to write header to output file: %v\n", err)
		return
	}

	// Write the data
	for _, dictWord := range dictList {
		pprime := 0
		if dictWord.GemProductPrime {
			pprime = 1
		}

		gprime := 0
		if dictWord.GemSumPrime {
			gprime = 1
		}

		record := []string{
			dictWord.DictionaryWordText,
			dictWord.RuneglishWordText,
			dictWord.RuneWordText,
			dictWord.RuneWordTextNoDoublet,
			strconv.Itoa(int(dictWord.GemSum)),
			strconv.Itoa(gprime),
			dictWord.GemProduct,
			strconv.Itoa(pprime),
			strconv.Itoa(dictWord.DictionaryWordLength),
			strconv.Itoa(dictWord.RuneglishWordLength),
			strconv.Itoa(dictWord.RuneWordLength),
			dictWord.RunePattern,
			dictWord.RunePatternNoDoubletPattern,
			strconv.Itoa(dictWord.DictRuneNoDoubletLength),
			dictWord.RuneDistancePattern,
			dictWord.Language,
		}
		if err := writer.Write(record); err != nil {
			fmt.Printf("Failed to write record to output file: %v\n", err)
			return
		}
	}

	fmt.Printf("CSV file created successfully: %s\n", outputFile)
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
	re := regexp.MustCompile(`[^\w]+`)
	for {
		line, readError := reader.ReadString('\n')
		if readError != nil {
			if readError.Error() != "EOF" {
				fmt.Printf("Error reading file: %v\n", err)
			}
			break
		}
		words := re.Split(line, -1)
		for _, word := range words {
			if word != "" {
				wordList[word] = true
			}
		}
	}
}

func filterNumbersOut(text string) string {
	wordArray := strings.Split(text, "")
	var newWordArray []string
	for _, character := range wordArray {
		if strings.ContainsAny(character, "ABCDEFGHIJKLMOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz'-") {
			newWordArray = append(newWordArray, character)
		}
	}
	return strings.Join(newWordArray, "")
}
