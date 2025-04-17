package cipher

import (
	"fmt"
	"math/big"
	"os"
	"runer"
	"runtime"
	"strings"
	"sync"
	"time"
)

// BulkDecryptAutokeyCipherRaw decodes the text using the Autokey cipher in a brute force fashion.
func BulkDecryptAutokeyCipherRaw(alphabet, wordList, latinList []string, text string, maxDepth int, file *os.File) (string, error) {
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
				decodedText := DecodeVigenereCipher(alphabet, keyArray, strings.Split(text, ""))

				if decodedText == "" {
					continue
				}

				decText := DecipheredText{
					Count: 0,
					Text:  decodedText,
					Key:   key,
				}
				resultsChan <- decText
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
		_, err := file.WriteString(fmt.Sprintf("%s : %s\n", decText.Key, decText.Text))
		if err != nil {
			fmt.Printf("Failed to write to file: %v", err)
		}
	}

	return result.String(), nil
}

// BulkDecryptAutokeyCipher decodes the text using the Vigenere cipher in a brute force fashion.
func BulkDecryptAutokeyCipher(alphabet, wordList, latinList []string, text string, maxDepth int) (string, error) {
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
				decodedText := DecryptAutokeyCipher(alphabet, keyArray, strings.Split(text, ""))

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

// DecryptAutokeyCipher decrypts a given ciphertext using the Autokey cipher.
func DecryptAutokeyCipher(alphabet, ciphertext, keyStream []string) string {
	var plaintext strings.Builder

	for _, c := range ciphertext {
		if index := indexOf(alphabet, c); index != -1 {
			// Get the key character from the keystream
			keyChar := keyStream[0]
			keyIndex := indexOf(alphabet, keyChar)

			// Decrypt the character
			plainIndex := (index - keyIndex + len(alphabet)) % len(alphabet)
			plainChar := alphabet[plainIndex]

			// Append the decrypted character to the plaintext
			plaintext.WriteString(plainChar)

			// Extend the keystream with the decrypted character
			keyStream = append(keyStream[1:], plainChar)
		} else {
			// Non-alphabetic characters are added as-is
			plaintext.WriteString(c)
		}
	}

	return plaintext.String()
}
