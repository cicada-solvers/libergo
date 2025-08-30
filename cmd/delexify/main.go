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

	err := processTextFile(*file)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

}

func processTextFile(path string) error {
	var lines []string
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

		if charRepo.ContainsLineSeperator(line) {
			tmpLines := seperateLine(line)
			for _, tmpLine := range tmpLines {
				lines = append(lines, tmpLine)
			}
		} else {
			lines = append(lines, line)
		}
	}

	if scanError := scanner.Err(); scanError != nil {
		return fmt.Errorf("error reading file %s: %w", path, scanError)
	}

	for _, line := range lines {
		words := getAllWords(line)
		for _, word := range words {
			if _, ok := mapOfWords[word]; ok {
				mapOfWords[word]++
			} else {
				mapOfWords[word] = 1
			}
		}
	}

	return nil
}

func seperateLine(line string) []string {
	lineArray := strings.Split(line, " ")
	var stringBuilder strings.Builder
	var lines []string

	for _, char := range lineArray {
		if charRepo.IsLineSeperator(char) {
			lines = append(lines, stringBuilder.String())
			stringBuilder.Reset()
		} else {
			stringBuilder.WriteString(char)
		}
	}

	return lines
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
