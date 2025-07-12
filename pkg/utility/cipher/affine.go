package cipher

import (
	"errors"
	"fmt"
	"runer"
	"strings"
)

// BulkDecodeAffineCipherRaw attempts to decode the given text for all affine cipher key combinations within the specified alphabet.
// It returns a concatenated string of all possible decodings or an error if decoding fails.
func BulkDecodeAffineCipherRaw(alphabet []string, text string) (string, error) {
	var result strings.Builder

	for a := 0; a < len(alphabet)+1; a++ {
		for b := 0; b < len(alphabet)+1; b++ {
			fmt.Printf("Trying %d, %d:\n", a, b)
			decoded, err := DecodeAffineCipher(text, a, b, alphabet)
			if err != nil {
				continue
			}

			result.WriteString(fmt.Sprintf("%d/%d : %s\n", a, b, decoded))

			fmt.Println(decoded)
		}
	}

	return result.String(), nil
}

// BulkDecodeAffineCipher attempts to decode a text using all possible affine cipher keys within the given alphabet size.
// Returns the decoded text with all key combinations and may optionally transpose decoded results to Latin characters.
// Takes the alphabet, the text to decode, and a decodeToLatin flag as arguments. Returns the resulting string and an error.
func BulkDecodeAffineCipher(alphabet []string, text string, decodeToLatin bool) (string, error) {
	var result strings.Builder

	for a := 0; a < len(alphabet)+1; a++ {
		for b := 0; b < len(alphabet)+1; b++ {
			fmt.Printf("Trying %d, %d:\n", a, b)
			decoded, err := DecodeAffineCipher(text, a, b, alphabet)
			if err != nil {
				continue
			}

			result.WriteString(fmt.Sprintf("Multiplier: %d, Shift: %d - %s\n", a, b, decoded))

			if decodeToLatin {
				// Decode to Latin if needed
				decodedLatin := runer.TransposeRuneToLatin(decoded)
				result.WriteString(fmt.Sprintf("Multiplier: %d, Shift: %d - %s\n", a, b, decodedLatin))
			}

			fmt.Println(decoded)
		}
	}

	return result.String(), nil
}

// ModInverse calculates the modular inverse of a under modulo m.
func ModInverse(a, m int) (int, error) {
	for x := 1; x < m; x++ {
		if (a*x)%m == 1 {
			return x, nil
		}
	}
	return 0, errors.New("no modular inverse found")
}

// DecodeAffineCipher decodes the given text using the affine cipher with the given multiplier and shift.
func DecodeAffineCipher(text string, a, b int, alphabet []string) (string, error) {
	m := len(alphabet)
	inverseA, err := ModInverse(a, m)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	for _, c := range text {
		if strings.ContainsRune(strings.Join(alphabet, ""), c) {
			index := strings.Index(strings.Join(alphabet, ""), strings.ToLower(string(c)))
			decodedIndex := (inverseA * (index - b + m)) % m
			decodedChar := alphabet[decodedIndex]
			if strings.ToUpper(string(c)) == string(c) {
				result.WriteString(strings.ToUpper(decodedChar))
			} else {
				result.WriteString(decodedChar)
			}
		} else {
			result.WriteRune(c)
		}
	}
	return result.String(), nil
}
