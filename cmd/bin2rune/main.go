package main

import (
	runelib "characterrepo"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

var charRepo *runelib.CharacterRepo

func main() {
	charRepo = runelib.NewCharacterRepo()
	outputBuilder := &strings.Builder{}

	// Parsing the flags
	inputFile := flag.String("input", "", "Input file")
	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Error: input must be specified")
		flag.Usage()
		return
	}

	// Reading the file data to get integers.
	data, err := os.ReadFile(*inputFile)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}

	// Convert the data to an int array
	var intArray []int
	for _, b := range data {
		// Convert the byte int an integer value.
		intArray = append(intArray, int(b))
	}

	// Now we loop through it to see if we need to reduce it or not
	for counter, intValue := range intArray {
		if counter > 0 {
			outputBuilder.WriteString("•")
		}

		if charRepo.IsPrimer(intValue) {
			outputBuilder.WriteString(charRepo.GetRuneFromValue(intValue))
		} else {
			breakValues := breakNonPrimerIntoPrimers(intValue)
			for _, breakValue := range breakValues {
				outputBuilder.WriteString(charRepo.GetRuneFromValue(breakValue))
			}
		}
	}

	outputBuilder.WriteString("⊹")
	fmt.Println(outputBuilder.String())
}

// breakNonPrimerIntoPrimers breaks a given integer into a series of prime numbers and returns them as a slice.
// The function uses randomness and recursively reprocesses the input if results are invalid.
func breakNonPrimerIntoPrimers(value int) []int {
	isGood := false
	var retval []int

	for !isGood {
		tmpValue := value

		// Get random number 1 through 29
		for tmpValue > 0 {
			_, maxPrimerPosition := charRepo.GetMaxPrimerAndPositionFromValue(tmpValue)

			random := 0
			if maxPrimerPosition == 0 {
				random = charRepo.GetPrimerFromPosition(0)
			} else {
				random = charRepo.GetPrimerFromPosition(rand.Intn(maxPrimerPosition))
			}

			tmpValue -= random
			retval = append(retval, random)
		}

		if tmpValue < 0 {
			isGood = false
			retval = retval[:0]
		} else {
			isGood = true
		}
	}

	return retval
}
