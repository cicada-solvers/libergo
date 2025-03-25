package lgstructs

import (
	"fmt"
	"strings"
)

type DictionaryWord struct {
	DictionaryWordText          string `json:"dict_word"`
	RuneglishWordText           string `json:"dict_runeglish"`
	RuneWordText                string `json:"dict_rune"`
	RuneWordTextNoDoublet       string `json:"dict_rune_no_doublet"`
	GemSum                      int64  `json:"gem_sum"`
	GemSumPrime                 bool   `json:"gem_sum_prime"`
	GemProduct                  string `json:"gem_product"`
	GemProductPrime             bool   `json:"gem_product_prime"`
	DictionaryWordLength        int    `json:"dict_word_length"`
	RuneglishWordLength         int    `json:"dict_runeglish_length"`
	RuneWordLength              int    `json:"dict_rune_length"`
	RunePattern                 string `json:"rune_pattern"`
	RunePatternNoDoubletPattern string `json:"rune_pattern_no_doublet"`
	Language                    string `json:"language"`
}

// GetRunePattern gets the rune pattern for the dictionary word
func GetRunePattern(word string) string {
	patternDictionary := make(map[int]string)
	var runes []string
	counter := 1

	for _, character := range word {
		if character == '\'' {
			runes = append(runes, "'")
			continue
		}

		found := false
		for key, value := range patternDictionary {
			if value == string(character) {
				runes = append(runes, fmt.Sprintf("%d", key))
				found = true
				break
			}
		}

		if !found {
			runes = append(runes, fmt.Sprintf("%d", counter))
			patternDictionary[counter] = string(character)
			counter++
		}
	}

	return strings.Join(runes, ",")
}

// RemoveDoublets removes consecutive duplicate characters from a word
func RemoveDoublets(word string) string {
	if len(word) == 0 {
		return word
	}

	var result strings.Builder
	result.WriteByte(word[0])

	for i := 1; i < len(word); i++ {
		if word[i] != word[i-1] {
			result.WriteByte(word[i])
		}
	}

	return result.String()
}
