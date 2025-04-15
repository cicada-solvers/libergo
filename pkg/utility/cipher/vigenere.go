package cipher

import (
	runelib "characterrepo"
	"fmt"
	"math/big"
	"runer"
	"runtime"
	"strings"
	"sync"
	"time"
)

var processedCounter = big.NewInt(0)
var rateCounter = big.NewInt(0)
var latinWordList []string
var topResults []DecipheredText

// BulkDecodeVigenereCipher decodes the text using the Vigenere cipher in a brute force fashion.
func BulkDecodeVigenereCipher(alphabet, wordList, latinList []string, text string, maxDepth int) (string, error) {
	latinWordList = latinList
	if maxDepth > 10 {
		return "", fmt.Errorf("max depth of %d is not allowed, the maximum allowed depth is 10", maxDepth)
	}

	// We are going to put timer to see how many we have processed.
	processedTicker := time.NewTicker(time.Minute)
	defer processedTicker.Stop()

	go func() {
		for range processedTicker.C {
			fmt.Printf("Rate: %s/min - Processed %s items\n", rateCounter.String(), processedCounter.String())
			rateCounter.SetInt64(int64(0))
		}
	}()

	var result strings.Builder
	combinations := generateCombinations(wordList, maxDepth)
	combinationChan := make(chan []string)
	resultsChan := make(chan DecipheredText)
	var wg sync.WaitGroup

	// Start worker goroutines
	numWorkers := runtime.NumCPU() + (runtime.NumCPU() / 2) // Adjust based on your system's capabilities
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for combination := range combinationChan {
				processedCounter.Add(processedCounter, big.NewInt(1))
				rateCounter.Add(rateCounter, big.NewInt(1))

				key := strings.Join(combination, "")
				keyArray := strings.Split(key, "")
				decodedText := DecodeVigenereCipher(alphabet, keyArray, text)

				if decodedText == "" {
					continue
				}

				latinText := runer.TransposeRuneToLatin(decodedText)
				totalWords := countWords(latinText)

				if totalWords > 0 {
					totalText := fmt.Sprintf("Decoded: %s\nKey: %s\nLatin: %s\nCount: %d\n\n", decodedText, key, latinText, totalWords)
					decText := DecipheredText{
						Count: totalWords,
						Text:  totalText,
						Key:   key,
					}
					resultsChan <- decText
				}
			}
		}()
	}

	// Send combinations to workers
	go func() {
		for combination := range combinations {
			combinationChan <- combination
		}
		close(combinationChan)
	}()

	// Close results channel when workers are done
	go func() {
		wg.Wait()

		close(resultsChan)
	}()

	// Collect results
	for decText := range resultsChan {
		topResults = append(topResults, decText)
		topResults = sortTopResults(topResults)
	}

	result.Reset()
	for _, key := range topResults {
		result.WriteString(key.Text)
	}

	return result.String(), nil
}

// DecodeVigenereCipher decodes the text using the Vigenere cipher.
func DecodeVigenereCipher(alphabet, key []string, text string) string {
	if len(key) == 0 {
		fmt.Printf("Key cannot be empty")
		return ""
	}

	var decodedText strings.Builder
	charRepo := runelib.NewCharacterRepo()
	keyIndex := 0

	for _, c := range text {
		if isLetter(c) || charRepo.IsRune(string(c), false) {
			keyChar := key[keyIndex]
			textIndex := indexOf(alphabet, string(c))
			keyIndexInAlphabet := indexOf(alphabet, keyChar)
			decodedCharIndex := (textIndex - keyIndexInAlphabet + len(alphabet)) % len(alphabet)
			decodedChar := alphabet[decodedCharIndex]

			if charRepo.IsRune(string(c), false) {
				decodedText.WriteString(decodedChar)
			} else {
				if isUpper(c) {
					decodedText.WriteString(strings.ToUpper(decodedChar))
				} else {
					decodedText.WriteString(strings.ToLower(decodedChar))
				}
			}
			keyIndex = (keyIndex + 1) % len(key)
		} else {
			decodedText.WriteRune(c)
		}
	}

	return decodedText.String()
}
