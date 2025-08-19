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

	repeatCounts = sortMapDesc(repeatCounts)

	for word, count := range repeatCounts {
		fmt.Printf("%s: %d\n", word, count)
	}
}

func sortMapDesc(m map[string]int64) map[string]int64 {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})
	sortedMap := make(map[string]int64)
	for _, k := range keys {
		sortedMap[k] = m[k]
	}

	return sortedMap
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
