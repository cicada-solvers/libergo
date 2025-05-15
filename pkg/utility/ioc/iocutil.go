package ioc

import "strings"

// AlphabetType represents different types of alphabets that can be used for IOC calculation
type AlphabetType int

const (
	// Latin represents the standard 26-letter Latin alphabet (a-z)
	Latin AlphabetType = iota
	// Runeglish represents the standard Latin alphabet plus common runes/symbols
	Runeglish
	// Rune represents all possible Unicode runes in the text
	Rune
)

// GetAlphabet returns the character set corresponding to the specified alphabet type
func GetAlphabet(alphabetType AlphabetType) []string {
	var retval string

	switch alphabetType {
	case Latin:
		retval = "abcdefghijklmnopqrstuvwxyz"
	case Runeglish:
		retval = "abcdefghijlmnoprstuwxy"
	case Rune:
		retval = "ᛝᛟᛇᛡᛠᚫᚦᚠᚢᚩᚱᚳᚷᚹᚻᚾᛁᛄᛈᛉᛋᛏᛒᛖᛗᛚᛞᚪᚣ"
	default:
		retval = "abcdefghijklmnopqrstuvwxyz"
	}

	return strings.Split(retval, "")
}

// CalcIOC calculates the incidence of coincidence for the given text using the provided alphabet.
// The incidence of coincidence is a measure used in cryptanalysis that
// reflects the likelihood of randomly selecting the same letter twice from a text.
// It returns a float64 value between 0 and 1.
// If an empty alphabet is provided, the function defaults to the standard English alphabet.
func CalcIOC(text string, alphabetType AlphabetType) float64 {
	// Get the alphabet
	alphabet := GetAlphabet(alphabetType)

	// Create a map for faster character lookup
	validChars := make(map[string]bool)
	for _, char := range alphabet {
		validChars[char] = true
	}

	// Create a map to count occurrences of each letter
	counts := make(map[string]int)

	// Count only characters in our alphabet
	totalLetters := 0
	for _, char := range strings.ToLower(text) {
		if validChars[string(char)] {
			counts[string(char)]++
			totalLetters++
		}
	}

	// If there are fewer than 2 letters, return 0
	if totalLetters <= 1 {
		return 0.0
	}

	// Calculate the sum of frequencies squared
	var sum float64 = 0.0
	for _, count := range counts {
		sum += float64(count) * float64(count)
	}

	// Calculate and return the IOC
	return sum / (float64(totalLetters) * float64(totalLetters-1))
}
