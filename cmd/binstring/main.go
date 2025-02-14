package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"titler"
)

func main() {
	titler.PrintTitle("Binary to String Converter")
	inputFile := flag.String("inputfile", "", "Input file")
	inputText := flag.String("inputtext", "", "Input text")
	outputFile := flag.String("outputfile", "", "Output file")
	flag.Parse()

	if *inputFile == "" && *inputText == "" {
		fmt.Println("Error: Either inputfile or inputtext must be specified")
		return
	}

	var input string
	if *inputFile != "" {
		data, err := os.ReadFile(*inputFile)
		if err != nil {
			fmt.Println("Error reading input file:", err)
			return
		}
		input = string(data)
	} else {
		input = *inputText
	}

	var output string
	if isBinaryString(input) {
		output = binaryStringToText(input)
	} else {
		output = textToBinaryString(input)
	}

	if *outputFile != "" {
		err := os.WriteFile(*outputFile, []byte(output), 0644)
		if err != nil {
			fmt.Println("Error writing to output file:", err)
		}
	} else {
		fmt.Println(output)
	}
}

func isBinaryString(s string) bool {
	for _, c := range s {
		if c != '0' && c != '1' && c != ' ' {
			return false
		}
	}
	return true
}

func binaryStringToText(s string) string {
	s = strings.ReplaceAll(s, " ", "")
	var result strings.Builder
	for i := 0; i < len(s); i += 8 {
		byteStr := s[i : i+8]
		var b byte
		for j := 0; j < 8; j++ {
			b = b<<1 + byte(byteStr[j]-'0')
		}
		result.WriteByte(b)
	}
	return result.String()
}

func textToBinaryString(s string) string {
	var result strings.Builder
	for i, c := range s {
		result.WriteString(fmt.Sprintf("%08b", c))
		if (i+1)%8 == 0 {
			result.WriteString(" ")
		}
	}
	return result.String()
}
