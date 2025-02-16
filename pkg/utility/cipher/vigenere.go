package cipher

import (
	"fmt"
	"gorm.io/gorm"
	"liberdatabase"
	"strings"
)

type AlphabetType int

const (
	Latin AlphabetType = iota
	Rune
)

func BulkDecodeVigenereCipher(alphabetType AlphabetType, alphabet []string, text string, maxDepth int) (string, error) {
	if maxDepth > 10 {
		return "", fmt.Errorf("max depth of %d is not allowed, the maximum allowed depth is 10", maxDepth)
	}

	connection, connError := liberdatabase.InitConnection()
	if connError != nil {
		return "", connError
	}
	defer func(db *gorm.DB) {
		err := liberdatabase.CloseConnection(db)
		if err != nil {
			// Handle error
		}
	}(connection)

	words, err := liberdatabase.GetAllWords(connection, 0)
	if err != nil {
		return "", err
	}

	var wordList []string
	for _, word := range words {
		switch alphabetType {
		case Latin:
			wordList = append(wordList, word.DictionaryWordText, word.RuneglishWordText)
		case Rune:
			wordList = append(wordList, word.RuneWordText)
		}
	}

	var result strings.Builder
	for depth := 1; depth <= maxDepth; depth++ {
		combinations := generateCombinations(wordList, depth)
		for combination := range combinations {
			decodedText := DecodeVigenereCipher(alphabet, strings.Join(combination, ""), text)
			result.WriteString(decodedText + "\n")
		}
	}

	return result.String(), nil
}

func generateCombinations(wordList []string, length int) <-chan []string {
	combinations := make(chan []string)

	go func() {
		defer close(combinations)
		generate(wordList, length, []string{}, combinations)
	}()

	return combinations
}

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
func DecodeVigenereCipher(alphabet []string, key, text string) string {
	var decodedText strings.Builder
	keyIndex := 0

	for _, c := range text {
		if isLetter(c) {
			keyChar := rune(key[keyIndex])
			textIndex := indexOf(alphabet, strings.ToUpper(string(c)))
			keyIndexInAlphabet := indexOf(alphabet, strings.ToUpper(string(keyChar)))
			decodedCharIndex := (textIndex - keyIndexInAlphabet + len(alphabet)) % len(alphabet)
			decodedChar := alphabet[decodedCharIndex]
			if isUpper(c) {
				decodedText.WriteString(strings.ToUpper(decodedChar))
			} else {
				decodedText.WriteString(strings.ToLower(decodedChar))
			}
			keyIndex = (keyIndex + 1) % len(key)
		} else {
			decodedText.WriteRune(c)
		}
	}

	return decodedText.String()
}
