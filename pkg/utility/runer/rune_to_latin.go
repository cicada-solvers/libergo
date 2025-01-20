package runer

import (
	runelib "characterrepo"
	"strings"
)

func TransposeRuneToLatin(text string) string {
	var sb strings.Builder
	repo := runelib.NewCharacterRepo()

	for _, runeCharacter := range text {
		character := repo.GetCharFromRune(string(runeCharacter))
		if character != "" {
			sb.WriteString(character)
		} else {
			sb.WriteRune(runeCharacter)
		}
	}

	return sb.String()
}
