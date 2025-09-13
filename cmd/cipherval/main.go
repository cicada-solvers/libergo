package main

import (
	"bufio"
	runelib "characterrepo"
	"cipher"
	"flag"
	"fmt"
	"liberdatabase"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// main initializes the program, parses input flags, validates them, and performs decoding based on the cipher type provided.
func main() {
	// Define the flags
	text := flag.String("text", "", "The text to decode")
	alphabet := flag.String("alphabet", "rune", "The alphabet to use (rune or english)")
	wordFile := flag.String("wordfile", "", "The text file of words to try for brute force decoding")
	cipherType := flag.String("ciphertype", "caesar", "The cipher to use (vigenere, atbash, affine, autokey, caesar, trithemius)")

	// Parse the flags
	flag.Parse()

	// Validate required flags
	if *text == "" {
		log.Fatal("The -text flag is required")
	}

	// Print the parsed flags (for debugging or further processing)
	fmt.Printf("Text: %s\n", *text)
	fmt.Printf("Alphabet: %s\n", *alphabet)
	fmt.Printf("Word File: %s\n", *wordFile)
	fmt.Printf("Cipher: %s\n", *cipherType)

	// Add your decoding logic here
	// Determine the alphabet to use
	var alphabetSet []string

	if strings.ToLower(*alphabet) == "rune" {
		repo := runelib.NewCharacterRepo()
		alphabetSet = repo.GetGematriaRunes()
	} else {
		alphabetSet = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
			"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	}

	// Database stuff
	id := uuid.NewString()
	_, _ = liberdatabase.InitTables()
	dbConn, _ := liberdatabase.InitConnection()

	// Now we are going to decode the text based on the cipher type
	var decodedText string
	var decodeErr error

	switch strings.ToLower(*cipherType) {
	case "caesar":
		list := liberdatabase.GetDictionaryWords(dbConn)
		fmt.Printf("List Length: %d\n", len(list))
		decodedText, decodeErr = cipher.BulkDecodeCaesarStringRaw(alphabetSet, strings.Split(*text, ""))
		if decodeErr != nil {
			fmt.Printf("Failed to decode using Caesar cipher: %v", decodeErr)
		}
		// Write the decoded text to the output file
		outputText := fmt.Sprintf("%s\n", decodedText)
		score := cipher.ScoreTextWithList(outputText, list)
		output := liberdatabase.OutputData{
			DocId: id,
			Data:  outputText,
			Score: score,
		}
		liberdatabase.AddOutputData(dbConn, output)
	case "affine":
		list := liberdatabase.GetDictionaryWords(dbConn)
		fmt.Printf("List Length: %d\n", len(list))
		decodedText, decodeErr = cipher.BulkDecodeAffineCipherRaw(alphabetSet, *text)
		if decodeErr != nil {
			fmt.Printf("Failed to decode using Affine cipher: %v", decodeErr)
		}
		// Write the decoded text to the output file
		outputText := fmt.Sprintf("%s\n", decodedText)
		score := cipher.ScoreTextWithList(outputText, list)
		output := liberdatabase.OutputData{
			DocId: id,
			Data:  outputText,
			Score: score,
		}
		liberdatabase.AddOutputData(dbConn, output)
	case "atbash":
		list := liberdatabase.GetDictionaryWords(dbConn)
		fmt.Printf("List Length: %d\n", len(list))
		decodedText, decodeErr = cipher.BulkDecodeAtbashStringRaw(alphabetSet, *text)
		if decodeErr != nil {
			fmt.Printf("Failed to decode using Atbash cipher: %v", decodeErr)
		}
		// Write the decoded text to the output file
		outputText := fmt.Sprintf("%s\n", decodedText)
		score := cipher.ScoreTextWithList(outputText, list)
		output := liberdatabase.OutputData{
			DocId: id,
			Data:  outputText,
			Score: score,
		}
		liberdatabase.AddOutputData(dbConn, output)
	case "trithemius":
		list := liberdatabase.GetDictionaryWords(dbConn)
		fmt.Printf("List Length: %d\n", len(list))
		decodedText, decodeErr = cipher.BulkDecodeTrithemiusStringRaw(alphabetSet, *text)
		if decodeErr != nil {
			fmt.Printf("Failed to decode using Trithemius cipher: %v", decodeErr)
		}
		// Write the decoded text to the output file
		outputText := fmt.Sprintf("%s\n", decodedText)
		score := cipher.ScoreTextWithList(outputText, list)
		output := liberdatabase.OutputData{
			DocId: id,
			Data:  outputText,
			Score: score,
		}
		liberdatabase.AddOutputData(dbConn, output)
	case "vigenere":
		if *wordFile == "" {
			log.Fatal("The -wordfile flag is required for Vigenere cipher")
		}

		scoreList := liberdatabase.GetDictionaryWords(dbConn)

		// Read words from the CSV file
		wordlist, csvErr := ReadWordsFromTextFile(*wordFile)
		if csvErr != nil {
			return
		}

		// The wait group is used to wait for all the goroutines to finish
		numWorkers := runtime.NumCPU()
		wg := sync.WaitGroup{}
		dbConns := make([]*gorm.DB, numWorkers)
		wordArrays := SplitStringSliceIntoParts(wordlist, numWorkers)

		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func() {
				dbConns[i], _ = liberdatabase.InitConnection()
				decodeErr = cipher.BulkDecodeVigenereCipherRaw(i, scoreList, alphabetSet, wordArrays[i], *text, dbConns[i])
				if decodeErr != nil {
					fmt.Printf("Failed to decode using Vigenere cipher: %v", decodeErr)
				}

				wg.Done()
			}()
		}

		wg.Wait()

		for i := 0; i < numWorkers; i++ {
			_ = liberdatabase.CloseConnection(dbConns[i])
		}
	case "autokey":
		if *wordFile == "" {
			log.Fatal("The -wordfile flag is required for autokey cipher")
		}

		scoreList := liberdatabase.GetDictionaryWords(dbConn)

		// Read words from the CSV file
		wordlist, csvErr := ReadWordsFromTextFile(*wordFile)
		if csvErr != nil {
			return
		}

		// The wait group is used to wait for all the goroutines to finish
		numWorkers := runtime.NumCPU()
		wg := sync.WaitGroup{}
		dbConns := make([]*gorm.DB, numWorkers)
		wordArrays := SplitStringSliceIntoParts(wordlist, numWorkers)

		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func() {
				dbConns[i], _ = liberdatabase.InitConnection()
				decodeErr = cipher.BulkDecryptAutokeyCipherRaw(i, scoreList, alphabetSet, wordArrays[i], *text, dbConns[i])
				if decodeErr != nil {
					fmt.Printf("Failed to decode using autokey cipher: %v", decodeErr)
				}

				wg.Done()
			}()
		}

		wg.Wait()

		for i := 0; i < numWorkers; i++ {
			_ = liberdatabase.CloseConnection(dbConns[i])
		}
	}

	_ = liberdatabase.CloseConnection(dbConn)
}

// ReadWordsFromTextFile reads all the words from a text file.
// Words are parsed using whitespace as separators, and common punctuation is trimmed.
func ReadWordsFromTextFile(filePath string) ([]string, error) {
	// Open the text file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			fmt.Printf("Failed to close file: %v", closeErr)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	// Increase the buffer and max token size to support long words/lines
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	scanner.Split(bufio.ScanWords)

	var words []string
	for scanner.Scan() {
		w := strings.Trim(scanner.Text(), " \t\r\n,;:.!?\"'()[]{}<>")
		if w == "" {
			continue
		}
		words = append(words, w)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return words, nil
}

// SplitStringSliceIntoParts splits items into exactly `parts` contiguous slices,
// distributing the elements as evenly as possible. Earlier parts receive one
// extra element if the division is not even. If parts <= 0, it returns an empty result.
// If parts > len(items), some parts will be empty.
func SplitStringSliceIntoParts(items []string, parts int) [][]string {
	if parts <= 0 {
		return [][]string{}
	}

	n := len(items)
	result := make([][]string, parts)

	q := 0 // base size for each part
	r := 0 // remainder (number of parts that get an extra item)
	if parts > 0 {
		q = n / parts
		r = n % parts
	}

	idx := 0
	for i := 0; i < parts; i++ {
		size := q
		if i < r {
			size++
		}

		if size > 0 {
			end := idx + size
			result[i] = items[idx:end]
			idx = end
		} else {
			// Ensure it's an empty (nil) slice; use []string{} if you prefer non-nil.
			result[i] = nil
		}
	}

	return result
}
