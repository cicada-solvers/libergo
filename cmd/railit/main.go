package main

import (
	"bufio"
	"cipher"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Rail Fence Cipher Decoder")
	fmt.Println("-------------------------")

	reader := bufio.NewReader(os.Stdin)

	// Get the ciphertext
	fmt.Print("Enter the text to decode: ")
	ciphertext, _ := reader.ReadString('\n')
	ciphertext = strings.TrimSpace(ciphertext)

	// Get the number of rails
	fmt.Print("Enter the number of rails: ")
	railsStr, _ := reader.ReadString('\n')
	railsStr = strings.TrimSpace(railsStr)

	rails, err := strconv.Atoi(railsStr)
	if err != nil || rails <= 0 {
		fmt.Println("Invalid number of rails. Using default value of 3.")
		rails = 3
	}

	// Decode the text
	plaintext := cipher.DecodeRailFence(strings.Split(ciphertext, ""), rails)

	// Print the result
	fmt.Println("\nDecoded text:")
	fmt.Println(plaintext)
}
