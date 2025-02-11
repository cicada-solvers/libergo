package main

import (
	"flag"
	"fmt"
	"liberdatabase"
	"log"
	"runer"
	"titler"
)

func main() {
	titler.PrintTitle("Get Words")

	// Define flags
	textTypeFlag := flag.String("textType", "runes", "Type of text: latin, runeglish, or runes")
	numberFlag := flag.Int("number", 0, "Number value")
	byFlag := flag.String("by", "", "Criteria: length, sum, or pattern")
	patternFlag := flag.String("pattern", "", "Pattern value")

	// Parse flags
	flag.Parse()

	// Map textTypeFlag to runer.TextType
	var textType runer.TextType
	switch *textTypeFlag {
	case "latin":
		textType = runer.Latin
	case "runeglish":
		textType = runer.Runeglish
	case "runes":
		textType = runer.Runes
	default:
		log.Fatalf("Invalid textType: %s", *textTypeFlag)
	}

	// Database connection (example using PostgreSQL)
	db, _connError := liberdatabase.InitConnection()
	if _connError != nil {
		log.Fatalf("Error connecting to database: %v", _connError)
	}

	// Handle the -by flag
	switch *byFlag {
	case "length":
		words, err := liberdatabase.GetWordsByLength(db, *numberFlag, textType)
		if err != nil {
			log.Fatalf("Error retrieving words by length: %v", err)
		}

		for _, word := range words {
			fmt.Println(word)
		}
	case "gemsum":
		words, err := liberdatabase.GetWordsByGemSum(db, int64(*numberFlag))
		if err != nil {
			log.Fatalf("Error retrieving words by gem sum: %v", err)
		}

		for _, word := range words {
			fmt.Println(word)
		}
	case "pattern":
		words, err := liberdatabase.GetWordsByPattern(db, *patternFlag)
		if err != nil {
			log.Fatalf("Error retrieving words by pattern: %v", err)
		}

		for _, word := range words {
			fmt.Println(word)
		}
	default:
		log.Fatalf("Invalid criteria: %s", *byFlag)
	}
}
