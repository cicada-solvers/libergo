package cipher

import (
	"math/big"
	"sort"
	"strings"
)

type AlphabetType int

const (
	Latin AlphabetType = iota
	Rune
)

type DecipheredText struct {
	Text  string
	Count int64
	Key   string
}

var processedCounter = big.NewInt(0)
var rateCounter = big.NewInt(0)
var latinWordList []string
var topResults []DecipheredText

// indexOf returns the index of the target string in the slice, or -1 if not found.
func indexOf(slice []string, target string) int {
	for i, v := range slice {
		if v == target {
			return i
		}
	}
	return -1
}

// isLetter checks if a string contains a single letter.
func isLetter(s string) bool {
	if len(s) != 1 {
		return false
	}
	c := rune(s[0])
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// isUpper checks if a string contains a single uppercase letter.
func isUpper(s string) bool {
	if len(s) != 1 {
		return false
	}
	c := rune(s[0])
	return c >= 'A' && c <= 'Z'
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

// countWords counts the number of words in the text.
func countWords(text string) int64 {
	var count int64
	for _, word := range latinWordList {
		if strings.Contains(text, word) {
			count = count + (int64(len(word)) * 2)
		}
	}
	return count
}

// sortTopResults sorts the top results based on the count of words.
func sortTopResults(results []DecipheredText) []DecipheredText {
	var sortedList []DecipheredText
	for _, v := range results {
		sortedList = append(sortedList, v)
	}

	sort.Slice(sortedList, func(i, j int) bool {
		return sortedList[i].Count > sortedList[j].Count
	})

	if len(sortedList) > 200 {
		sortedList = sortedList[:200]
	}

	return sortedList
}
