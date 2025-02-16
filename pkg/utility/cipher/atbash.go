package cipher

import (
	"fmt"
	"strings"
)

func BulkDecodeAtbashString(alphabet []string, text string) (string, error) {
	var result strings.Builder

	for i := 0; i < len(alphabet); i++ {
		// Move the last character to the first position
		newAlphabet := append([]string{alphabet[len(alphabet)-1]}, alphabet[:len(alphabet)-1]...)
		// Decode the text with the new alphabet
		decoded := DecodeAtbashCipher(text, newAlphabet)
		result.WriteString(fmt.Sprintf("Shift: %d - %s\n", i, decoded))
		// Update the alphabet for the next iteration
		alphabet = newAlphabet
	}

	return result.String(), nil
}

// EncodeAtbashCipher encodes the given text using the Atbash cipher.
func EncodeAtbashCipher(text string, alphabet []string) string {
	var result strings.Builder

	for _, c := range text {
		if isLetter(c) {
			index := indexOf(alphabet, string(c))
			reversedIndex := len(alphabet) - 1 - index
			reversedChar := alphabet[reversedIndex]
			if isUpper(c) {
				result.WriteString(strings.ToUpper(reversedChar))
			} else {
				result.WriteString(reversedChar)
			}
		} else {
			result.WriteRune(c)
		}
	}

	return result.String()
}

// DecodeAtbashCipher decodes the given text using the Atbash cipher.
func DecodeAtbashCipher(text string, alphabet []string) string {
	var result strings.Builder

	for _, c := range text {
		if isLetter(c) {
			index := indexOf(alphabet, string(c))
			reversedIndex := len(alphabet) - 1 - index
			reversedChar := alphabet[reversedIndex]
			if isUpper(c) {
				result.WriteString(strings.ToUpper(reversedChar))
			} else {
				result.WriteString(reversedChar)
			}
		} else {
			result.WriteRune(c)
		}
	}

	return result.String()
}

// isLetter checks if a rune is a letter.
func isLetter(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// isUpper checks if a rune is an uppercase letter.
func isUpper(c rune) bool {
	return c >= 'A' && c <= 'Z'
}
