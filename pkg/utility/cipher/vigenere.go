package cipher

import (
	runelib "characterrepo"
	"fmt"
	"liberdatabase"
	"runer"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BulkDecodeVigenereCipherRaw decodes the text using the Vigenere cipher in a brute force fashion.
func BulkDecodeVigenereCipherRaw(alphabet, wordList []string, text string, db *gorm.DB) error {
	id := uuid.NewString()
	list := liberdatabase.GetDictionaryWords(db)

	for _, key := range wordList {
		keyArray := strings.Split(key, "")
		decodedText := DecodeVigenereCipher(alphabet, keyArray, strings.Split(text, ""))
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

func DecodeVigenereCipher(alphabet, key, text []string) string {
	if len(key) == 0 {
		fmt.Printf("Key cannot be empty")
		return ""
	}
	if len(alphabet) == 0 {
		fmt.Printf("Alphabet cannot be empty")
		return ""
	}

	// Build a fast lookup for the alphabet
	alphaIndex := make(map[string]int, len(alphabet))
	for i, a := range alphabet {
		alphaIndex[a] = i
	}

	// Pre-clean the key: drop empties and normalize for lookup
	cleanKey := make([]string, 0, len(key))
	for _, k := range key {
		if k == "" {
			continue
		}
		kl := strings.ToLower(k)
		if _, ok := alphaIndex[kl]; ok {
			cleanKey = append(cleanKey, kl)
		}
	}
	if len(cleanKey) == 0 {
		// No usable key symbols in the given alphabet
		return ""
	}

	var decodedText strings.Builder
	charRepo := runelib.NewCharacterRepo()
	keyIndex := 0

	for _, c := range text {
		if c == "" {
			// Ignore empty splits from strings.Split(s, "")
			continue
		}

		if isLetter(c) || charRepo.IsRune(c, false) {
			// Normalize for lookup if this is a Latin letter (case-sensitive alphabets should adapt this)
			inLookup := c
			restoreUpper := false

			if !charRepo.IsRune(c, false) {
				if isUpper(c) {
					restoreUpper = true
				}
				inLookup = strings.ToLower(c)
			}

			ti, okT := alphaIndex[inLookup]
			ki, okK := alphaIndex[cleanKey[keyIndex]]
			if !okT || !okK {
				// Character or key symbol not in alphabet: keep as-is and do not advance key
				decodedText.WriteString(c)
				continue
			}

			decodedChar := alphabet[(ti-ki+len(alphabet))%len(alphabet)]

			if charRepo.IsRune(c, false) {
				decodedText.WriteString(decodedChar)
			} else if restoreUpper {
				decodedText.WriteString(strings.ToUpper(decodedChar))
			} else {
				decodedText.WriteString(strings.ToLower(decodedChar))
			}

			keyIndex = (keyIndex + 1) % len(cleanKey)
		} else {
			decodedText.WriteString(c)
		}
	}

	return decodedText.String()
}
