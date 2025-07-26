package main

import (
	runelib "characterrepo"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var charRepo *runelib.CharacterRepo

// The main is the entry point of the application.
// It processes input text, performs transformations, and writes a binary file.
func main() {
	builder := strings.Builder{}
	textPtr := flag.String("text", "", "The text to process")
	outputPtr := flag.String("output", "", "The output file name")

	// Parse the flags
	flag.Parse()

	// Validate input
	if *textPtr == "" {
		fmt.Println("Error: Text is required")
		flag.Usage()
		os.Exit(1)
	}

	if *outputPtr == "" {
		fmt.Println("Error: Output file name is not supported")
		flag.Usage()
		os.Exit(1)
	}

	charRepo = runelib.NewCharacterRepo()

	replacer := strings.NewReplacer(
		"•", "|",
		"␍", "|",
		"␊", "|",
		"␗", "|",
		",", "|",
		" ", "|",
		".", "|",
		"\"", "|",
		"'", "|",
		"!", "|",
		"@", "|",
		"#", "|",
		"$", "|",
		"%", "|",
		"^", "|",
		"&", "|",
		"*", "|",
		"(", "|",
		")", "|",
		"-", "|",
		"_", "|",
		"=", "|",
		"+", "|",
		"[", "|",
		"]", "|",
		"{", "|",
		"}", "|",
	)
	processed := replacer.Replace(*textPtr)

	// Then split on the common separator
	runeArray := strings.Split(processed, ",")

	// Calculate each gem sum, then write it out to a binary string
	for _, runeStr := range runeArray {
		gumSum := charRepo.CalculateGemSum(runeStr)

		// Convert the integer value to binary string
		binStr := padBinaryString(fmt.Sprintf("%08b", gumSum))
		builder.WriteString(binStr)
	}

	// Now we are going to break the string into 8-piece chunks
	bytes := binStringToBytes(builder.String())

	// now write the byte array out to a file
	err := os.WriteFile(*outputPtr, bytes, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		os.Exit(1)
	}
}

// padBinaryString ensures a binary string is a multiple of 8 bits by adding leading zeros if necessary.
func padBinaryString(binStr string) string {
	// Calculate how many bits we need to add to make it a multiple of 8
	remainder := len(binStr) % 8

	// If already a multiple of 8, no padding needed
	if remainder == 0 {
		return binStr
	}

	// Add the required number of padding bits (0s) to make it a multiple of 8
	padding := 8 - remainder
	paddedBinStr := strings.Repeat("0", padding) + binStr

	return paddedBinStr
}

// binStringToBytes converts a binary string into a byte slice, grouping bits into 8-bit segments.
func binStringToBytes(binStr string) []byte {
	stringArray := strings.Split(binStr, "")
	tmpBuilder := strings.Builder{}
	var retval []byte

	for i := 0; i < len(stringArray); i++ {
		tmpBuilder.WriteString(stringArray[i])
		if len(tmpBuilder.String())%8 == 0 {
			byteVal, parseError := strconv.ParseUint(tmpBuilder.String(), 2, 8)
			if parseError != nil {
				fmt.Println("Error parsing binary string:", parseError)
				os.Exit(1)
			}
			retval = append(retval, byte(byteVal))
			tmpBuilder.Reset()
		}
	}

	return retval
}
