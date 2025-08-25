package main

import (
	"bufio"
	runelib "characterrepo"
	"flag"
	"fmt"
	"os"
	"strings"
)

var charRepo *runelib.CharacterRepo

func main() {
	charRepo = runelib.NewCharacterRepo()
	file := flag.String("file", "", "File to delexify")
	flag.Parse()

	if *file == "" {
		flag.Usage()
		return
	}

}

func processTextFile(path string) error {
	mapOfWords := make(map[string]int64)
	fmt.Printf("Processing file %s\n", path)

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			continue
		}

		line = strings.ToLower(line)
		words := getAllWords(line)

		for _, word := range words {
			_, keyExists := mapOfWords[word]
			if keyExists {
				mapOfWords[word]++
			} else {
				mapOfWords[word] = 1
			}
		}
	}

	if scanError := scanner.Err(); scanError != nil {
		return fmt.Errorf("error reading file %s: %w", path, scanError)
	}

	return nil
}

func getAllWords(line string) []string {
	lineArray := strings.Split(line, " ")
	var words []string
	var wordBuilder strings.Builder

	// Pre-allocate space for words to reduce reallocations
	words = make([]string, 0, 16) // Assuming average of ~16 words per line

	// Iterate through runes directly
	for _, r := range lineArray {
		if charRepo.IsSeperator(string(r)) || charRepo.IsDinkus(string(r)) || charRepo.IsLineSeperator(string(r)) {
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
