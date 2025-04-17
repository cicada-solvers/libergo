package cipher

import (
	runelib "characterrepo"
	"fmt"
	"runer"
	"strings"
)

func BulkDecodeTrithemiusStringRaw(alphabet []string, text string) (string, error) {
	var result strings.Builder

	for i := 0; i < len(alphabet); i++ {
		// Move the last character to the first position
		newAlphabet := append([]string{alphabet[len(alphabet)-1]}, alphabet[:len(alphabet)-1]...)
		// Decode the text with the new alphabet
		decoded := DecryptTrithemiusCipher(newAlphabet, strings.Split(text, ""))
		result.WriteString(fmt.Sprintf("Shift: %d : %s\n", i, decoded))

		// Update the alphabet for the next iteration
		alphabet = newAlphabet
	}

	return result.String(), nil
}

func BulkDecodeTrithemiusString(alphabet []string, text string, decodeToLatin bool) (string, error) {
	var result strings.Builder

	for i := 0; i < len(alphabet); i++ {
		// Move the last character to the first position
		newAlphabet := append([]string{alphabet[len(alphabet)-1]}, alphabet[:len(alphabet)-1]...)
		// Decode the text with the new alphabet
		decoded := DecryptTrithemiusCipher(newAlphabet, strings.Split(text, ""))
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

// DecryptTrithemiusCipher decrypts the text using the Trithemius cipher.
func DecryptTrithemiusCipher(alphabet, text []string) string {
	var decryptedText strings.Builder
	charRepo := runelib.NewCharacterRepo()

	for i, c := range text {
		if isLetter(c) || charRepo.IsRune(c, false) {
			textIndex := indexOf(alphabet, strings.ToUpper(string(c)))
			shift := i % len(alphabet) // Reverse the shift based on position
			decryptedCharIndex := (textIndex - shift + len(alphabet)) % len(alphabet)
			decryptedChar := alphabet[decryptedCharIndex]
			if isUpper(c) {
				decryptedText.WriteString(strings.ToUpper(decryptedChar))
			} else {
				decryptedText.WriteString(strings.ToLower(decryptedChar))
			}
		} else {
			decryptedText.WriteString(c)
		}
	}

	return decryptedText.String()
}
