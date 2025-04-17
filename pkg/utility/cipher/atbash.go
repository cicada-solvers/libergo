package cipher

import (
	runelib "characterrepo"
	"fmt"
	"runer"
	"strings"
)

func BulkDecodeAtbashStringRaw(alphabet []string, text string) (string, error) {
	var result strings.Builder

	for i := 0; i < len(alphabet); i++ {
		// Move the last character to the first position
		newAlphabet := append([]string{alphabet[len(alphabet)-1]}, alphabet[:len(alphabet)-1]...)
		// Decode the text with the new alphabet
		decoded := DecodeAtbashCipher(strings.Split(text, ""), newAlphabet)
		result.WriteString(fmt.Sprintf("%d : %s\n", i, decoded))

		// Update the alphabet for the next iteration
		alphabet = newAlphabet
	}

	return result.String(), nil
}

func BulkDecodeAtbashString(alphabet []string, text string, decodeToLatin bool) (string, error) {
	var result strings.Builder

	for i := 0; i < len(alphabet); i++ {
		// Move the last character to the first position
		newAlphabet := append([]string{alphabet[len(alphabet)-1]}, alphabet[:len(alphabet)-1]...)
		// Decode the text with the new alphabet
		decoded := DecodeAtbashCipher(strings.Split(text, ""), newAlphabet)
		result.WriteString(fmt.Sprintf("Shift: %d - %s\n", i, decoded))

		if decodeToLatin {
			// Decode the text to Latin if required
			decodedLatin := runer.TransposeRuneToLatin(decoded)
			result.WriteString(fmt.Sprintf("Decoded to Latin: %s\n", decodedLatin))
		}

		// Update the alphabet for the next iteration
		alphabet = newAlphabet
	}

	return result.String(), nil
}

// DecodeAtbashCipher decodes the given text using the Atbash cipher.
func DecodeAtbashCipher(text, alphabet []string) string {
	var result strings.Builder
	charRepo := runelib.NewCharacterRepo()

	for _, c := range text {
		if isLetter(c) || charRepo.IsRune(string(c), false) {
			index := indexOf(alphabet, string(c))
			reversedIndex := len(alphabet) - 1 - index
			reversedChar := alphabet[reversedIndex]
			if isUpper(c) {
				result.WriteString(strings.ToUpper(reversedChar))
			} else {
				result.WriteString(reversedChar)
			}
		} else {
			result.WriteString(c)
		}
	}

	return result.String()
}
