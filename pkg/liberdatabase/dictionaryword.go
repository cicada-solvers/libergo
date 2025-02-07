package liberdatabase

import (
	"fmt"
	"gorm.io/gorm"

	"strings"
)

type DictionaryWord struct {
	gorm.Model
	DictionaryWordText   string `gorm:"column:dict_word"`
	RuneglishWordText    string `gorm:"column:dict_runeglish"`
	RuneWordText         string `gorm:"column:dict_rune"`
	GemSum               int64  `gorm:"column:gem_sum"`
	DictionaryWordLength int    `gorm:"column:dict_word_length"`
	RuneglishWordLength  int    `gorm:"column:dict_runeglish_length"`
	RuneWordLength       int    `gorm:"column:dict_rune_length"`
}

func (DictionaryWord) TableName() string {
	return "public.dictionary_words"
}

// GetRunePattern gets the rune pattern for the dictionary word
func (dw DictionaryWord) GetRunePattern() string {
	patternDictionary := make(map[int]string)
	runes := []string{}
	counter := 1

	for _, character := range dw.RuneWordText {
		if character == '\'' {
			runes = append(runes, "'")
			continue
		}

		found := false
		for key, value := range patternDictionary {
			if value == string(character) {
				runes = append(runes, string(rune(key)))
				found = true
				break
			}
		}

		if !found {
			runes = append(runes, string(rune(counter)))
			patternDictionary[counter] = string(character)
			counter++
		}
	}

	return strings.Join(runes, "")
}

// String returns a string representation of the dictionary word
func (dw DictionaryWord) String() string {
	return fmt.Sprintf("%s - %s - %s - %d", dw.DictionaryWordText, dw.RuneglishWordText, dw.RuneWordText, dw.GemSum)
}
