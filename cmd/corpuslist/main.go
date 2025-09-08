package main

import (
	"bufio"
	"flag"
	"fmt"
	"liberdatabase"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"runer"
	"sequences"
	"strings"
	"unicode/utf8"
)

// main reads text files from a directory and creates a SQL file with the dictionary words
func main() {
	// Define flags for the directory path and output SQL file
	reverseWords := flag.Bool("reverse", false, "Reverse the words in the sentence")
	dir := flag.String("dir", "./path/to/directory", "Directory to scan for text files")
	flag.Parse()

	wordList := make(map[string]bool)
	dictList := make([]liberdatabase.DictionaryWord, 0, 16384)

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

	_, _ = liberdatabase.InitTables()
	dbConn, _ := liberdatabase.InitConnection()

	// Read every word from the word list
	for word := range wordList {
		word = strings.ToUpper(word)
		word = filterNumbersOut(word)

		if len(word) <= 0 {
			continue
		}

		runeglish := runer.PrepLatinToRune(word)
		runeText := runer.TransposeLatinToRune(runeglish, *reverseWords)
		runeTextNoDoublet := liberdatabase.RemoveDoublets(strings.Split(runeText, ""))
		gemSum := runer.CalculateGemSum(runeText, runer.Runes, false)
		gemProd := runer.CalculateGemProduct(runeText, runer.Runes, false)

		dictWord := liberdatabase.DictionaryWord{
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
			RunePattern:                 liberdatabase.GetRunePattern(word),
			RunePatternNoDoubletPattern: liberdatabase.GetRunePattern(runeTextNoDoublet),
			RuneDistancePattern:         liberdatabase.GetRuneDistancePattern(strings.Split(runeText, "")),
			Language:                    "English",
		}

		dictList = append(dictList, dictWord)
		delete(wordList, word) // Remove the word from wordList

		if len(dictList) >= 500 {
			liberdatabase.AddDictionaryWords(dbConn, dictList)
			dictList = dictList[:0]
		}
	}

	liberdatabase.AddDictionaryWords(dbConn, dictList)

	_ = liberdatabase.CloseConnection(dbConn)
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

// filterNumbersOut removes all numeric characters from the input string, retaining only alphabetic characters and valid symbols.
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
