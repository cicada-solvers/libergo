package main

import (
	"decoder"
	"flag"
	"fmt"
	"log"
	"os"
	"titler"
)

// main is the entry point of the program, orchestrating Base64 decoding based on provided command-line flags.
func main() {
	titler.PrintTitle("Decode Base64")

	// Define command-line flags
	input := flag.String("input", "", "Base64 encoded string")
	hexFlag := flag.Bool("hex", false, "Decode as hex string")

	// Parse command-line flags
	flag.Parse()

	// Check if input flag is provided
	if *input == "" {
		fmt.Println("Error: input flag is required")
		flag.Usage()
		os.Exit(1)
	}

	// Create DecodeCommand
	cmd := &decoder.DecodeCommand{
		Input:    *input,
		Encoding: "PLAIN",
	}

	if *hexFlag {
		cmd.Encoding = "HEX"
	}

	// Call DecodeBase64String
	decoded, err := decoder.DecodeBase64String(cmd)
	if err != nil {
		log.Fatalf("Error decoding base64 string: %v", err)
	}

	// Output the decoded string
	fmt.Println(decoded)
}
