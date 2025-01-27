package main

import (
	"fmt"
	"math/big"
	"os"
	"sequences"
)

// main is the entry point of the program.
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: isprime <number>")
		os.Exit(1)
	}

	numberStr := os.Args[1]
	number, ok := new(big.Int).SetString(numberStr, 10)
	if !ok {
		fmt.Printf("Invalid number: %s\n", numberStr)
		os.Exit(1)
	}

	if sequences.IsPrime(number) {
		fmt.Printf("%s is a prime number.\n", numberStr)
	} else {
		fmt.Printf("%s is not a prime number.\n", numberStr)
	}
}
