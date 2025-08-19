package main

import (
	runelib "characterrepo"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	var repeatCounts = map[string]int64{}
	filePtr := flag.String("file", "", "The file to process")

	flag.Parse()

	if *filePtr == "" {
		flag.Usage()
		return
	}

	fileText, fileError := os.ReadFile(*filePtr)
	if fileError != nil {
		fmt.Println("Error reading file:", fileError)
		return
	}

	fmt.Println(string(fileText))
	lines := strings.Split(string(fileText), "\n")
	for _, line := range lines {
		words := getAllWords(line)
		for _, word := range words {
			if _, exists := repeatCounts[word]; exists {
				repeatCounts[word]++
			} else {
				repeatCounts[word] = 1
			}
		}
	}

	// Create a slice of key-value pairs for sorting
	type wordCount struct {
		word  string
		count int64
	}
	pairs := make([]wordCount, 0, len(repeatCounts))
	for word, count := range repeatCounts {
		pairs = append(pairs, wordCount{word, count})
	}

	// Sort pairs by count in descending order
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})

	// Print sorted results
	for _, pair := range pairs {
		fmt.Printf("%s: %d\n", pair.word, pair.count)
	}
}

func getAllWords(line string) []string {
	charRepo := runelib.NewCharacterRepo()
	lineArray := strings.Split(line, "")
	var wordBuilder strings.Builder

	words := make([]string, 0, 16)

	// Iterate through runes directly
	for _, r := range lineArray {
		if charRepo.IsDinkus(r) || charRepo.IsSeperator(r) {
			words = append(words, wordBuilder.String())
			wordBuilder.Reset()
		} else {
			wordBuilder.WriteString(r)
		}
	}

	// Add the last word if the line ends with a letter
	if wordBuilder.Len() > 0 {
		words = append(words, wordBuilder.String())
	}

	return words
}
