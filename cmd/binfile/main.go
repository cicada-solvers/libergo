package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"titler"
)

// main is the entry point of the program that converts binary string files to text or text files to binary strings.
func main() {
	titler.PrintTitle("Binary to String File Converter")
	inputFile := flag.String("inputfile", "", "Input file")
	outputFile := flag.String("outputfile", "", "Output file")
	flag.Parse()

	if *inputFile == "" || *outputFile == "" {
		fmt.Println("Error: Both inputfile and outputfile must be specified")
		flag.Usage()
		return
	}

	data, err := os.ReadFile(*inputFile)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}

	var output string
	if isBinaryString(string(data)) {
		output = binaryStringToText(string(data))
	} else {
		output = textToBinaryString(data)
	}

	err = os.WriteFile(*outputFile, []byte(output), 0644)
	if err != nil {
		fmt.Println("Error writing to output file:", err)
	}
}

// isBinaryString checks if a given string contains only binary digits ('0' and '1') or spaces.
func isBinaryString(s string) bool {
	for _, c := range s {
		if c != '0' && c != '1' && c != ' ' {
			return false
		}
	}
	return true
}

// binaryStringToText converts a binary string into its corresponding text representation.
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

// textToBinaryString converts a byte slice to a binary string representation separated by spaces.
func textToBinaryString(data []byte) string {
	var result strings.Builder
	for _, b := range data {
		result.WriteString(fmt.Sprintf("%08b", b))
		result.WriteString(" ")
	}
	return result.String()
}
