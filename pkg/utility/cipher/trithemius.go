package cipher

import (
	runelib "characterrepo"
	"fmt"
	"runer"
	"strings"
)

func BulkDecodeTrithemiusString(alphabet []string, text string, decodeToLatin bool) (string, error) {
	var result strings.Builder

	for i := 0; i < len(alphabet); i++ {
		// Move the last character to the first position
		newAlphabet := append([]string{alphabet[len(alphabet)-1]}, alphabet[:len(alphabet)-1]...)
		// Decode the text with the new alphabet
		decoded := DecryptTrithemiusCipher(newAlphabet, text)
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

// EncryptTrithemiusCipher encrypts the text using the Trithemius cipher.
func EncryptTrithemiusCipher(alphabet []string, text string) string {
	var encryptedText strings.Builder

	for i, c := range text {
		if isLetter(c) {
			textIndex := indexOf(alphabet, strings.ToUpper(string(c)))
			shift := i % len(alphabet) // Shift based on position
			encryptedCharIndex := (textIndex + shift) % len(alphabet)
			encryptedChar := alphabet[encryptedCharIndex]
			if isUpper(c) {
				encryptedText.WriteString(strings.ToUpper(encryptedChar))
			} else {
				encryptedText.WriteString(strings.ToLower(encryptedChar))
			}
		} else {
			encryptedText.WriteRune(c)
		}
	}

	return encryptedText.String()
}

// DecryptTrithemiusCipher decrypts the text using the Trithemius cipher.
func DecryptTrithemiusCipher(alphabet []string, text string) string {
	var decryptedText strings.Builder
	charRepo := runelib.NewCharacterRepo()

	for i, c := range text {
		if isLetter(c) || charRepo.IsRune(string(c), false) {
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
			decryptedText.WriteRune(c)
		}
	}

	return decryptedText.String()
}
