package lgstructs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
func GetRunePattern(dw DictionaryWord) string {
	patternDictionary := make(map[int]string)
	var runes []string
	counter := 1

	for _, character := range dw.RuneWordText {
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

func GetWordsFromApi(field, value string) ([]DictionaryWord, error) {
	url := fmt.Sprintf("https://cmbsolver.com/cmbsolver-api/runewords.php/%s/%s", field, value)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing response body")
		}
	}(resp.Body)

	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(resp)

	if resp.Header.Get("Content-Type") != "application/json" {
		return nil, fmt.Errorf("unexpected content type: %s", resp.Header.Get("Content-Type"))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var fileTypeInfoModels []DictionaryWord
	err = json.Unmarshal(body, &fileTypeInfoModels)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return fileTypeInfoModels, nil
}
