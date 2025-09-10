package cipher

import (
	"fmt"
	"liberdatabase"
	"runer"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BulkDecryptAutokeyCipherRaw decodes the text using the Autokey cipher in a brute force fashion.
func BulkDecryptAutokeyCipherRaw(alphabet, wordList []string, text string, db *gorm.DB) error {
	id := uuid.NewString()
	list := liberdatabase.GetDictionaryWords(db)

	for _, key := range wordList {
		keyArray := strings.Split(key, "")
		decodedText := DecryptAutokeyCipher(alphabet, strings.Split(text, ""), keyArray)
		latinText := runer.TransposeRuneToLatin(decodedText)

		outputText := fmt.Sprintf("Decoded: %s\nKey: %s\nLatin:%s\n\n", decodedText, key, latinText)
		fmt.Println(outputText)
		score := ScoreTextWithList(db, outputText, list)
		output := liberdatabase.OutputData{
			DocId: id,
			Data:  outputText,
			Score: score,
		}
		db.Create(&output)
	}

	return nil
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
