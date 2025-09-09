package cipher

import (
	runelib "characterrepo"
	"fmt"
	"math/big"
	"os"
	"runer"
	"runtime"
	"strings"
	"sync"
	"time"
)

// BulkDecodeVigenereCipherRaw decodes the text using the Vigenere cipher in a brute force fashion.
func BulkDecodeVigenereCipherRaw(alphabet, wordList []string, text string, maxDepth int, file *os.File) (string, error) {
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

// BulkDecodeVigenereCipher decodes the text using the Vigenere cipher in a brute force fashion.
func BulkDecodeVigenereCipher(alphabet, wordList []string, text string, maxDepth int) (string, error) {
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

				latinText := runer.TransposeRuneToLatin(decodedText)

				totalText := fmt.Sprintf("Decoded: %s\nKey: %s\nLatin: %s\n\n", decodedText, key, latinText)
				decText := DecipheredText{
					Count: 0,
					Text:  totalText,
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
		topResults = append(topResults, decText)
		topResults = sortTopResults(topResults)
	}

	result.Reset()
	for _, key := range topResults {
		result.WriteString(key.Text)
	}

	return result.String(), nil
}

func DecodeVigenereCipher(alphabet, key, text []string) string {
	if len(key) == 0 {
		fmt.Printf("Key cannot be empty")
		return ""
	}
	if len(alphabet) == 0 {
		fmt.Printf("Alphabet cannot be empty")
		return ""
	}

	// Build a fast lookup for the alphabet
	alphaIndex := make(map[string]int, len(alphabet))
	for i, a := range alphabet {
		alphaIndex[a] = i
	}

	// Pre-clean the key: drop empties and normalize for lookup
	cleanKey := make([]string, 0, len(key))
	for _, k := range key {
		if k == "" {
			continue
		}
		kl := strings.ToLower(k)
		if _, ok := alphaIndex[kl]; ok {
			cleanKey = append(cleanKey, kl)
		}
	}
	if len(cleanKey) == 0 {
		// No usable key symbols in the given alphabet
		return ""
	}

	var decodedText strings.Builder
	charRepo := runelib.NewCharacterRepo()
	keyIndex := 0

	for _, c := range text {
		if c == "" {
			// Ignore empty splits from strings.Split(s, "")
			continue
		}

		if isLetter(c) || charRepo.IsRune(c, false) {
			// Normalize for lookup if this is a Latin letter (case-sensitive alphabets should adapt this)
			inLookup := c
			restoreUpper := false

			if !charRepo.IsRune(c, false) {
				if isUpper(c) {
					restoreUpper = true
				}
				inLookup = strings.ToLower(c)
			}

			ti, okT := alphaIndex[inLookup]
			ki, okK := alphaIndex[cleanKey[keyIndex]]
			if !okT || !okK {
				// Character or key symbol not in alphabet: keep as-is and do not advance key
				decodedText.WriteString(c)
				continue
			}

			decodedChar := alphabet[(ti-ki+len(alphabet))%len(alphabet)]

			if charRepo.IsRune(c, false) {
				decodedText.WriteString(decodedChar)
			} else if restoreUpper {
				decodedText.WriteString(strings.ToUpper(decodedChar))
			} else {
				decodedText.WriteString(strings.ToLower(decodedChar))
			}

			keyIndex = (keyIndex + 1) % len(cleanKey)
		} else {
			decodedText.WriteString(c)
		}
	}

	return decodedText.String()
}
