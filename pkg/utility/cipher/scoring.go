package cipher

import (
	"liberdatabase"
	"strings"

	"gorm.io/gorm"
)

func ScoreText(db *gorm.DB, text string) int64 {
	words := liberdatabase.GetDictionaryWords(db)
	score := int64(0)
	for _, word := range words {
		if strings.Contains(text, word) {
			score += int64(len(word))
		}
	}

	return score
}

func ScoreTextWithList(db *gorm.DB, text string, list []string) int64 {
	score := int64(0)
	for _, word := range list {
		if strings.Contains(text, word) {
			score += int64(len(word))
		}
	}

	return score
}
