package main

import (
	"characterrepo"
	"flag"
	"fmt"
	"os"
	"strings"
)

var repo = runelib.NewCharacterRepo()

func main() {
	// Parse the flags
	flag.Parse()

	// Get non-flag arguments
	args := flag.Args()

	// Check if the text was provided as a parameter
	if len(args) < 1 {
		fmt.Println("Error: No text provided for analysis")
		fmt.Println("Usage: ziggy <text>")
		flag.Usage()
		os.Exit(1)
	}

	// Get the text from the first argument
	text := args[0]
	var textArray []string

	// Remove any spaces from the text
	fmt.Println("Remove spaces from text? (Y/N):")
	removeSpaces := "N"
	_, err := fmt.Scanln(&removeSpaces)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	if removeSpaces == "Y" {
		textArray = RemoveSpaces(text)
	} else {
		textArray = strings.Split(text, "")
	}

	fmt.Printf("The number of characters is %d\n", len(textArray))

	var possibleRails []int

	for i := 3; i <= 100; i++ {
		// Get the remainder of I against the number of characters
		remainder := len(textArray) % i
		if remainder == 0 {
			fmt.Printf("Possible rails: %d\n", i)
			possibleRails = append(possibleRails, i)
		}
	}

	if len(possibleRails) == 0 {
		fmt.Println("No possible rails found")
		return
	}

	// Prompt the user to select a rail
	fmt.Println("Select a rail:")
	var rail int
	_, scanErr := fmt.Scanln(&rail)
	if scanErr != nil {
		fmt.Println("Error: ", scanErr)
		return
	}

	// Divide the text into rail pieces
	fmt.Println("Divide long or tall rail? (L/T):")
	longTall := "L"
	_, ltError := fmt.Scanln(&longTall)
	if ltError != nil {
		fmt.Println("Error: ", ltError)
		return
	}
	var rails map[int][]string
	if longTall == "L" {
		rails = divideTextLong(textArray, rail)
	} else {
		rails = divideTextTall(textArray, rail)
	}

	fmt.Printf("Rails: \n %v\n", rails)

	// Decode the text using the rail fence cipher
	decodedText := DecodeUsingRailFenceCipher(rails)
	fmt.Printf("Decoded text: \n %s\n", strings.Join(decodedText, ""))

	fmt.Println("Done")
}

// divideText divides the text into rail pieces
func divideTextLong(textArray []string, rail int) map[int][]string {
	retval := make(map[int][]string)
	railLength := len(textArray) / rail
	railLengthCounter := 0
	currentRail := 0
	var railText []string
	for i := 0; i < len(textArray); i++ {
		character := textArray[i]
		railText = append(railText, character)
		if railLengthCounter == railLength-1 {
			fmt.Println("Adding rail")
			retval[currentRail] = railText
			railLengthCounter = 0
			railText = []string{}
			currentRail++
		} else {
			railLengthCounter++
		}
	}

	if len(railText) > 0 {
		fmt.Println("Adding rail")
		retval[currentRail] = railText
		railLengthCounter = 0
		railText = []string{}
	}

	return retval
}

func divideTextTall(textArray []string, rail int) map[int][]string {
	retval := make(map[int][]string)
	currentRail := 0

	for i := 0; i < len(textArray); i++ {
		character := textArray[i]
		retval[currentRail] = append(retval[currentRail], character)
		if currentRail == rail-1 {
			currentRail = 0
		} else {
			currentRail++
		}
	}

	return retval
}

func RemoveSpaces(text string) []string {
	var retval []string
	textArray := strings.Split(text, "")

	for _, character := range textArray {
		if repo.IsLetterInAlphabet(character) || repo.IsRune(character, false) {
			retval = append(retval, character)
		}
	}

	return retval
}

func DecodeUsingRailFenceCipher(rails map[int][]string) []string {
	var retArray []string
	currentRailPosition := 0
	downDirection := true
	firstTime := true

	for i := 0; i < len(rails[0]); i++ {
		character := rails[currentRailPosition][i]
		retArray = append(retArray, character)

		if !firstTime {
			if currentRailPosition == len(rails)-1 || currentRailPosition == 0 {
				downDirection = !downDirection
			}
		} else {
			firstTime = false
		}

		if downDirection {
			currentRailPosition++
		} else {
			currentRailPosition--
		}
	}

	return retArray
}
