package cipher

import (
	"fmt"
	"runer"
	"strings"
)

// BulkDecodeCaesarString decodes the given text using the Caesar cipher with the provided alphabet and key.
func BulkDecodeCaesarString(alphabet []string, text []string, decodeToLatin bool) (string, error) {
	var result strings.Builder

	for i := 0; i < len(alphabet); i++ {
		fmt.Printf("Trying %d:\n", i)
		decoded := DecodeCaesarCipher(alphabet, text, []int{i})
		result.WriteString(fmt.Sprintf("Shift: %d - %s\n", i, strings.Join(decoded, "")))

		if decodeToLatin {
			decodedLatin := runer.TransposeRuneToLatin(strings.Join(decoded, ""))
			result.WriteString(fmt.Sprintf("Shift: %d - %s\n", i, decodedLatin))
		}
	}

	return result.String(), nil
}

func EncodeCaesarCipher(alphabet, text []string, key []int) []string {
	var result []string

	for i, char := range text {
		alphabetIndex := indexOf(alphabet, char)

		if alphabetIndex != -1 {
			shift := key[0]
			if len(key) > 1 {
				shift = key[i]
			}
			newIndex := (alphabetIndex + shift) % len(alphabet)
			result = append(result, alphabet[newIndex])
		} else {
			result = append(result, char)
		}
	}

	return result
}

func DecodeCaesarCipher(alphabet, text []string, key []int) []string {
	var result []string

	for i, char := range text {
		alphabetIndex := indexOf(alphabet, char)

		if alphabetIndex != -1 {
			shift := key[0]
			if len(key) > 1 {
				shift = key[i]
			}
			shift = shift % len(alphabet) // Ensure shift is within bounds
			newIndex := (alphabetIndex - shift + len(alphabet)) % len(alphabet)
			result = append(result, alphabet[newIndex])
		} else {
			result = append(result, char)
		}
	}

	return result
}
