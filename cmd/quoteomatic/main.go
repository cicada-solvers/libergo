package main

import (
	"bufio"
	"flag"
	"fmt"
	"liberdatabase"
	"os"
	"runer"
	"strings"
)

func main() {
	textInFlag := flag.String("text", "", "Text to match")
	authorFlag := flag.String("author", "blake", "Author to use")
	patternOrGemSumFlag := flag.String("patternOrGemSum", "pattern", "Pattern or Gem Sum to use")

	flag.Parse()

	if *textInFlag == "" {
		fmt.Println("Error: No text provided for analysis")
		flag.Usage()
		os.Exit(1)
	}

	highestPercent := 0.0
	highestPercentQuote := ""

	quoteFile := fmt.Sprintf("%s.txt", *authorFlag)

	file, err := os.Open(quoteFile)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			continue
		}

		quoteParts := strings.Split(line, "|")
		quote := quoteParts[1]

		fmt.Printf("Checking Quote: %s - %s\n", quoteParts[0], quoteParts[1])

		quoteRunes := runer.PrepLatinToRune(quote)
		quoteRunes = runer.TransposeLatinToRune(quoteRunes, false)

		if *patternOrGemSumFlag == "pattern" {
			quotePattern := liberdatabase.GetRuneLinePattern(quoteRunes)
			textPattern := liberdatabase.GetRuneLinePattern(*textInFlag)
			matchCount := 0

			shortestLength := 0
			if len(quotePattern) < len(textPattern) {
				shortestLength = len(quotePattern)
			} else {
				shortestLength = len(textPattern)
			}

			for i := 0; i < shortestLength; i++ {
				if quotePattern[i] == textPattern[i] {
					matchCount++
				}
			}

			matchPercentage := float64(matchCount) / float64(shortestLength) * 100

			if matchPercentage > highestPercent {
				highestPercent = matchPercentage
				highestPercentQuote = fmt.Sprintf("Highest Percentage: %s - %f (%d)\n Quote: %s\n", quoteParts[0], matchPercentage, matchCount, quoteParts[1])
			}

		} else {
			fmt.Println("Gem Sum Not Implemented Yet")
		}

		fmt.Printf("%s\n", highestPercentQuote)
	}

	if scanError := scanner.Err(); scanError != nil {
		fmt.Printf("Error reading file: %v\n", scanError)
		os.Exit(1)
	}

}
