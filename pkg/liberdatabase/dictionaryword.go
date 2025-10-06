package liberdatabase

import (
	runelib "characterrepo"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// DictionaryWord represents a structured word entry with various properties for linguistic and numerical analysis.
type DictionaryWord struct {
	DictionaryWordText          string `gorm:"column:dict_word"`
	RuneglishWordText           string `gorm:"column:dict_runeglish"`
	RuneWordText                string `gorm:"column:dict_rune"`
	RuneWordTextNoDoublet       string `gorm:"column:dict_rune_no_doublet"`
	GemSum                      int64  `gorm:"column:gem_sum"`
	GemSumPrime                 bool   `gorm:"column:gem_sum_prime"`
	GemProduct                  string `gorm:"column:gem_product"`
	GemProductPrime             bool   `gorm:"column:gem_product_prime"`
	DictionaryWordLength        int    `gorm:"column:dict_word_length"`
	RuneglishWordLength         int    `gorm:"column:dict_runeglish_length"`
	RuneWordLength              int    `gorm:"column:dict_rune_length"`
	DictRuneNoDoubletLength     int    `gorm:"column:dict_rune_no_doublet_length"`
	RunePattern                 string `gorm:"column:rune_pattern"`
	RunePatternNoDoubletPattern string `gorm:"column:rune_pattern_no_doublet"`
	RuneDistancePattern         string `gorm:"column:rune_distance_pattern"`
	Language                    string `gorm:"column:language"`
}

func (DictionaryWord) TableName() string {
	return "dictionary_words"
}

func AddDictionaryWords(db *gorm.DB, dictionaryWords []DictionaryWord) {
	db.Create(&dictionaryWords)
	return
}

func GetDictionaryWords(db *gorm.DB) []string {
	var dictionaryWords []DictionaryWord
	var retval []string
	counter := 0

	db.Distinct().Find(&dictionaryWords)

	for _, word := range dictionaryWords {
		if !slices.Contains(retval, word.RuneglishWordText) {
			if counter == math.MaxInt16 {
				fmt.Printf("Loading %s - %s\n", word.DictionaryWordText, word.RuneglishWordText)
				counter = 0
			}
			counter++

			retval = append(retval, word.RuneglishWordText)
		}
	}

	return retval
}

func GetDictionaryWordsByRuneLength(db *gorm.DB, length int) []string {
	var dictionaryWords []DictionaryWord
	var retval []string
	db.Where("dict_rune_length = ?", length).Distinct().Find(&dictionaryWords)

	for _, word := range dictionaryWords {
		if !slices.Contains(retval, word.RuneWordText) {
			retval = append(retval, word.RuneWordText)
		}
	}

	return retval
}

func getAllWords(line string) []string {
	lettersArray := strings.Split("ᛝᛟᛇᛡᛠᚫᚦᚠᚢᚩᚱᚳᚷᚹᚻᚾᛁᛄᛈᛉᛋᛏᛒᛖᛗᛚᛞᚪᚣ'", "")
	letterMap := make(map[rune]bool, len(lettersArray))
	for _, letter := range lettersArray {
		letterMap[rune(letter[0])] = true
	}
	var words []string
	var wordBuilder strings.Builder

	// Pre-allocate space for words to reduce reallocations
	words = make([]string, 0, 16) // Assuming average of ~16 words per line

	// Iterate through runes directly
	for _, r := range line {
		if letterMap[r] {
			wordBuilder.WriteRune(r)
		} else if wordBuilder.Len() > 0 {
			words = append(words, wordBuilder.String())
			wordBuilder.Reset()
		}
	}

	// Add the last word if the line ends with a letter
	if wordBuilder.Len() > 0 {
		words = append(words, wordBuilder.String())
	}

	return words
}

func GetRuneLinePattern(line string) []int {
	patternArray := make([]int, 0)

	words := getAllWords(line)
	for _, word := range words {
		wordArray := strings.Split(word, "")

		if len(wordArray) == 0 {
			continue
		}

		patternArray = append(patternArray, len(wordArray))
	}

	return patternArray
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
func RemoveDoublets(word []string) string {
	if len(word) == 0 {
		return ""
	}

	var result strings.Builder
	result.WriteString(word[0])

	for i := 1; i < len(word); i++ {
		if word[i] != word[i-1] {
			result.WriteString(word[i])
		}
	}

	return result.String()
}

// CalculateWordDistances calculates the distances between corresponding runes of two words based on their positions in runeList.
// It returns a comma-separated string of absolute distance values for each rune comparison.
func CalculateWordDistances(word1, word2, runeList []string) string {
	var result strings.Builder
	var isFirstTime = true

	for i := 0; i < len(word1); i++ {
		pos1 := getRuneIndex(word1[i], runeList)
		pos2 := getRuneIndex(word2[i], runeList)

		currentDistance := pos1 - pos2
		distance := int(math.Abs(float64(currentDistance)))

		if isFirstTime {
			result.WriteString(fmt.Sprintf("%s", strconv.Itoa(distance)))
		} else {
			result.WriteString(fmt.Sprintf(", %s", strconv.Itoa(distance)))
		}

		isFirstTime = false
	}

	return result.String()
}

// GetRuneDistancePattern generates a string representing the distances between adjacent runes in a slice of strings.
func GetRuneDistancePattern(word []string) string {
	charRepo := runelib.NewCharacterRepo()
	gemRunes := charRepo.GetGematriaRunes()

	if len(word) == 0 {
		return ""
	}

	var result strings.Builder
	currentValue := getRuneIndex(word[0], gemRunes)

	result.WriteString(strconv.Itoa(0))

	for i := 1; i < len(word); i++ {
		currentDistance := currentValue - getRuneIndex(word[i], gemRunes)
		distance := int(math.Abs(float64(currentDistance)))
		result.WriteString(fmt.Sprintf(", %s", strconv.Itoa(distance)))
		currentValue = getRuneIndex(word[i], gemRunes)
	}

	return result.String()
}

// GetRuneComparisonDistancePattern calculates a distance pattern between corresponding runes in two word slices.
func GetRuneComparisonDistancePattern(wordOne, wordTwo []string) string {
	charRepo := runelib.NewCharacterRepo()
	gemRunes := charRepo.GetGematriaRunes()

	if len(wordOne) == 0 {
		return ""
	}

	var result strings.Builder
	result.WriteString(strconv.Itoa(0))

	for i := 1; i < len(wordOne); i++ {
		currentDistance := getRuneIndex(wordOne[i], gemRunes) - getRuneIndex(wordTwo[i], gemRunes)
		distance := int(math.Abs(float64(currentDistance)))
		result.WriteString(fmt.Sprintf(", %s", strconv.Itoa(distance)))
	}

	return result.String()
}

// getRuneIndex finds the index of a given rune in the provided alphabet slice.
// Returns the index if found, otherwise returns -1.
func getRuneIndex(rune string, alphabet []string) int {
	for i, r := range alphabet {
		if r == rune {
			return i
		}
	}

	return -1
}
