package cipher

import (
	runelib "characterrepo"
	"fmt"
	"github.com/jdkato/prose/v2"
	"log"
	"math/big"
	"runer"
	"runtime"
	"strings"
	"sync"
	"time"
)

var processedCounter = big.NewInt(0)
var rateCounter = big.NewInt(0)

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
	resultsChan := make(chan string)
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
				posCounts, totalWords := analyzeText(latinText)
				probability := calculateSentenceProbability(posCounts, totalWords)

				if probability > 0 {
					resultsChan <- fmt.Sprintf("Decoded: %s\nKey: %s\nLatin: %s\nProbability: %.2f\n\n", decodedText, key, latinText, probability)
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
	for res := range resultsChan {
		result.WriteString(res)
	}

	return result.String(), nil
}

// generateCombinations generates all combinations of words from the word list.
func generateCombinations(wordList []string, length int) <-chan []string {
	combinations := make(chan []string)

	go func() {
		defer close(combinations)
		generate(wordList, length, []string{}, combinations)
	}()

	return combinations
}

// generate generates all combinations of words from the word list.
func generate(wordList []string, length int, current []string, combinations chan<- []string) {
	if length == 0 {
		combinations <- append([]string{}, current...)
		return
	}

	for i, word := range wordList {
		generate(wordList[i:], length-1, append(current, word), combinations)
	}
}

// EncodeVigenereCipher encodes the text using the Vigenere cipher.
func EncodeVigenereCipher(alphabet []string, key, text string) string {
	var encodedText strings.Builder
	keyIndex := 0

	for _, c := range text {
		if isLetter(c) {
			keyChar := rune(key[keyIndex])
			textIndex := indexOf(alphabet, strings.ToUpper(string(c)))
			keyIndexInAlphabet := indexOf(alphabet, strings.ToUpper(string(keyChar)))
			encodedCharIndex := (textIndex + keyIndexInAlphabet) % len(alphabet)
			encodedChar := alphabet[encodedCharIndex]
			if isUpper(c) {
				encodedText.WriteString(strings.ToUpper(encodedChar))
			} else {
				encodedText.WriteString(strings.ToLower(encodedChar))
			}
			keyIndex = (keyIndex + 1) % len(key)
		} else {
			encodedText.WriteRune(c)
		}
	}

	return encodedText.String()
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

// analyzeText analyzes the given text and returns the part-of-speech counts and total word count.
func analyzeText(text string) (map[string]int, int) {
	doc, err := prose.NewDocument(text)
	if err != nil {
		log.Fatalf("Failed to create document: %v", err)
	}

	posCounts := map[string]int{
		"Noun":        0,
		"Verb":        0,
		"Adjective":   0,
		"Adverb":      0,
		"Determiner":  0,
		"Conjunction": 0,
		"Preposition": 0,
		"Pronoun":     0,
		"Punctuation": 0,
		"NamedEntity": 0,
	}
	totalWords := 0

	for _, tok := range doc.Tokens() {
		switch tok.Tag {
		case "NN", "NNS", "NNP", "NNPS":
			posCounts["Noun"]++
		case "VB", "VBD", "VBG", "VBN", "VBP", "VBZ":
			posCounts["Verb"]++
		case "JJ", "JJR", "JJS":
			posCounts["Adjective"]++
		case "RB", "RBR", "RBS":
			posCounts["Adverb"]++
		case "DT":
			posCounts["Determiner"]++
		case "CC":
			posCounts["Conjunction"]++
		case "IN":
			posCounts["Preposition"]++
		case "PRP", "PRP$", "WP", "WP$":
			posCounts["Pronoun"]++
		case ".", ",", ":", ";", "!", "?":
			posCounts["Punctuation"]++
		}
		totalWords++
	}

	posCounts["NamedEntity"] = len(doc.Entities())

	return posCounts, totalWords
}

// calculateSentenceProbability calculates the probability of a sentence being a valid English sentence.
func calculateSentenceProbability(posCounts map[string]int, totalWords int) float64 {
	if totalWords == 0 {
		return 0.0
	}

	probability := 0.0
	if posCounts["Noun"] > 0 && posCounts["Verb"] > 0 {
		probability = 50.0
		if posCounts["Adjective"] > 0 {
			probability += 10.0
		}
		if posCounts["Adverb"] > 0 {
			probability += 10.0
		}
		if posCounts["Determiner"] > 0 {
			probability += 5.0
		}
		if posCounts["Conjunction"] > 0 {
			probability += 5.0
		}
		if posCounts["Preposition"] > 0 {
			probability += 5.0
		}
		if posCounts["Pronoun"] > 0 {
			probability += 5.0
		}
		if posCounts["Punctuation"] > 0 {
			probability += 10.0
		}
		if posCounts["NamedEntity"] > 0 {
			probability += 5.0
		}
	}

	return probability
}
