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
func BulkDecryptAutokeyCipherRaw(threadId int, scorelist, alphabet, wordList []string, text string, db *gorm.DB) error {
	counter := int64(0)
	id := uuid.NewString()
	fmt.Printf("List Length: %d\n", len(scorelist))

	for _, key := range wordList {
		keyArray := strings.Split(key, "")
		decodedText := DecryptAutokeyCipher(alphabet, strings.Split(text, ""), keyArray)
		latinText := runer.TransposeRuneToLatin(decodedText)

		outputText := fmt.Sprintf("Latin: %s\nKey: %s\nAlphabet: %v\nDecoded:%s\n", latinText, key, alphabet, decodedText)
		score := ScoreTextWithList(outputText, scorelist)

		output := liberdatabase.OutputData{
			DocId: id,
			Data:  outputText,
			Score: score,
		}
		db.Create(&output)

		counter++
		if counter%10000 == 0 {
			fmt.Printf("%d - Decoded %d/%d\n", threadId, counter, len(wordList))
		}
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
