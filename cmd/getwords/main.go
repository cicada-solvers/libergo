package main

import (
	"flag"
	"fmt"
	"liberdatabase"
	"log"
	"math/big"
	"os"
	"runer"
	"strconv"
	"titler"
)

func main() {
	titler.PrintTitle("Get Words")

	// Define flags
	textTypeFlag := flag.String("textType", "runes", "Type of text: latin, runeglish, or runes")
	numberFlag := flag.String("number", "0", "Number value")
	byFlag := flag.String("by", "", "Criteria: length, gemsum, gemproduct, or pattern")
	patternFlag := flag.String("pattern", "", "Pattern value")

	// Parse flags
	flag.Parse()

	// Check if required flags are provided
	if *textTypeFlag == "" || *byFlag == "" {
		fmt.Println("Error: textType and by flags are required")
		flag.Usage()
		os.Exit(1)
	}

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
		number, err := strconv.Atoi(*numberFlag)
		if err != nil {
			log.Fatalf("Error converting numberFlag to int: %v", err)
		}
		words, err := liberdatabase.GetWordsByLength(db, number, textType)
		if err != nil {
			log.Fatalf("Error retrieving words by length: %v", err)
		}

		for _, word := range words {
			fmt.Println(word)
		}
	case "gemsum":
		number, err := strconv.ParseInt(*numberFlag, 10, 64)
		if err != nil {
			log.Fatalf("Error converting numberFlag to int64: %v", err)
		}
		words, err := liberdatabase.GetWordsByGemSum(db, number)
		if err != nil {
			log.Fatalf("Error retrieving words by gem sum: %v", err)
		}

		for _, word := range words {
			fmt.Println(word)
		}
	case "gemproduct":
		gemProduct := new(big.Int)
		gemProduct, ok := gemProduct.SetString(*patternFlag, 10)
		if !ok {
			log.Fatalf("Error converting pattern to big.Int: %s", *patternFlag)
		}
		words, err := liberdatabase.GetWordsByGemProduct(db, *gemProduct)
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
