package main

import (
	"decoder"
	"flag"
	"fmt"
	"log"
	"titler"
)

func main() {
	titler.PrintTitle("Decode Base64")

	// Define command-line flags
	input := flag.String("input", "", "Base64 encoded string")
	hexFlag := flag.Bool("hex", false, "Decode as hex string")

	// Parse command-line flags
	flag.Parse()

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
