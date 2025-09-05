package main

import (
	runelib "characterrepo"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

// charRepo is a global variable that stores the character repository.
var charRepo *runelib.CharacterRepo

// main is the entry point of the program.
func main() {
	textFlag := flag.String("text", "", "Text to shift")
	shiftFlag := flag.Int("shift", 0, "Number of positions to shift to the left")
	textDirectionFlag := flag.String("direction", "left", "Direction to shift text (left or right)")
	flag.Parse()

	if *textFlag == "" || *shiftFlag == 0 {
		flag.Usage()
		return
	}

	var output strings.Builder
	charRepo = runelib.NewCharacterRepo()
	words := getWordsFromText(*textFlag)
	for _, word := range words {
		letters := strings.Split(word, "")
		direction := *textDirectionFlag
		directionValue := *shiftFlag

		if strings.Contains(word, "|") {
			wordShiftCombo := strings.Split(word, "|")
			directionValue, _ = strconv.Atoi(wordShiftCombo[1])
			letters = strings.Split(wordShiftCombo[0], "")

			if directionValue >= 0 {
				direction = "right"
			} else {
				direction = "left"
			}
		}

		fmt.Printf("Word: %s - Direction %s - Shift %d \n", word, direction, directionValue)

		if direction == "right" {
			result := shiftLettersRight(letters, directionValue)
			output.WriteString(fmt.Sprintf("%s•", result))
		} else {
			result := shiftLettersLeft(letters, directionValue)
			output.WriteString(fmt.Sprintf("%s•", result))
		}
	}

	fmt.Println(output.String())
}

// ShiftLettersLeft shifts the letters in the text to the left by the specified shift.
func shiftLettersLeft(text []string, shift int) string {
	n := len(text)
	if n == 0 {
		return strings.Join(text, "")
	}
	// Normalize shift to [0, n)
	s := shift % n
	if s < 0 {
		s += n
	}
	rotated := append(append([]string{}, text[s:]...), text[:s]...)
	return strings.Join(rotated, "")
}

// ShiftLettersRight shifts the letters in the text to the right by the specified shift.
func shiftLettersRight(text []string, shift int) string {
	n := len(text)
	if n == 0 {
		return strings.Join(text, "")
	}
	// Normalize shift to [0, n)
	s := shift % n
	if s < 0 {
		s += n
	}
	if s == 0 {
		return strings.Join(text, "")
	}
	// Right rotation by s == left rotation by n - s
	left := n - s
	rotated := append(append([]string{}, text[left:]...), text[:left]...)
	return strings.Join(rotated, "")
}

// getWordsFromText splits the text into words based on the predefined character set.
func getWordsFromText(text string) []string {
	textArray := strings.Split(text, "")
	var words []string
	var currentWord strings.Builder

	for _, char := range textArray {
		if charRepo.IsDinkus(char) || charRepo.IsLineSeperator(char) || char == " " {
			words = append(words, currentWord.String())
			currentWord.Reset()
		} else {
			currentWord.WriteString(char)
		}
	}

	if currentWord.Len() > 0 {
		words = append(words, currentWord.String())
		currentWord.Reset()
	}

	return words
}
