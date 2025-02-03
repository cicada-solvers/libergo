package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"sequences"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide the number of bits as an argument.")
		return
	}

	bits, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Invalid number of bits: %v\n", err)
		return
	}

	// Generate a random integer with the specified number of bits
	maxNum := new(big.Int).Lsh(big.NewInt(1), uint(bits))
	var n *big.Int
	for {
		n, err = rand.Int(rand.Reader, maxNum)
		if err != nil {
			fmt.Printf("Error generating random integer: %v\n", err)
			return
		}

		// Check if the number is prime
		if !sequences.IsPrime(n) {
			break
		}
	}

	fmt.Printf("Generated %d-bit integer: %s\n", bits, n.String())
}
