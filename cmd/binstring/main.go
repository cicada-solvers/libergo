package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"titler"
)

// main is the entry point of the application, responsible for handling flags, input/output, and invoking core functions.
func main() {
	titler.PrintTitle("Binary to String Converter")
	inputFile := flag.String("inputfile", "", "Input file")
	inputText := flag.String("inputtext", "", "Input text")
	outputFile := flag.String("outputfile", "", "Output file")
	flag.Parse()

	if *inputFile == "" && *inputText == "" {
		fmt.Println("Error: Either inputfile or inputtext must be specified")
		flag.Usage()
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

// isBinaryString checks if the given string contains only binary digits ('0' and '1') and spaces. Returns true if valid.
func isBinaryString(s string) bool {
	for _, c := range s {
		if c != '0' && c != '1' && c != ' ' {
			return false
		}
	}
	return true
}

// binaryStringToText converts a binary string representation into its corresponding text by decoding 8-bit binary groups.
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

// textToBinaryString converts a given string into its binary representation, with each character encoded as an 8-bit binary value.
func textToBinaryString(s string) string {
	var result strings.Builder
	for _, c := range s {
		result.WriteString(fmt.Sprintf("%08b", c))
		result.WriteString(" ")
	}
	return result.String()
}
