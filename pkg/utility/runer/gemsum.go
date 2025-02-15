package runer

import (
	runelib "characterrepo"
	"math/big"
	"strings"
)

// TextType represents the type of text
type TextType int

// TextType constants
const (
	Latin TextType = iota
	Runeglish
	Runes
)

// String returns a string representation of the TextType
func (t TextType) String() string {
	return [...]string{"Latin", "Runeglish", "Runes"}[t]
}

// CalculateGemSum calculates the gem sum of a given string
func CalculateGemSum(gem string, textType TextType) int64 {
	repo := runelib.NewCharacterRepo()
	var retval int64
	var runeText string

	switch textType {
	case Latin:
		prep := PrepLatinToRune(strings.ToUpper(gem))
		runeText = TransposeLatinToRune(prep)
	case Runeglish:
		runeText = TransposeLatinToRune(strings.ToUpper(gem))
	case Runes:
		runeText = gem
	}

	for _, runeCharacter := range runeText {
		retval += int64(repo.GetValueFromRune(string(runeCharacter)))
	}
	return retval
}

// CalculateGemProduct calculates the product of the gem values of the given text
func CalculateGemProduct(gem string, textType TextType) big.Int {
	repo := runelib.NewCharacterRepo()
	var retval big.Int
	var zero big.Int
	var runeText string

	switch textType {
	case Latin:
		prep := PrepLatinToRune(strings.ToUpper(gem))
		runeText = TransposeLatinToRune(prep)
	case Runeglish:
		runeText = TransposeLatinToRune(strings.ToUpper(gem))
	case Runes:
		runeText = gem
	}

	for _, runeCharacter := range runeText {
		runeValue := big.NewInt(int64(repo.GetValueFromRune(string(runeCharacter))))
		if retval.Cmp(&zero) == 0 {
			retval.Set(runeValue)
		} else {
			retval.Mul(&retval, runeValue)
		}
	}
	return retval
}
