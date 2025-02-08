package runer

import (
	runelib "characterrepo"
	"strings"
)

type TextType int

const (
	Latin TextType = iota
	Runeglish
	Runes
)

func (t TextType) String() string {
	return [...]string{"Latin", "Runeglish", "Runes"}[t]
}

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
