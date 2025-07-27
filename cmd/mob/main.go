package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

// main is the entry point of the program. It validates the input, calculates the Möbius function, and prints the result.
func main() {
	// Validate the number of arguments
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <number>")
	}

	// Parse the number from the arguments
	number, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil || number <= 0 {
		log.Fatal("The input number must be a positive integer")
	}

	// Calculate the Möbius function
	result := mobiusFunction(number)

	// Print the result
	fmt.Printf("The Möbius function μ(%d) = %d\n", number, result)
}

// mobiusFunction calculates the Möbius function μ(n)
func mobiusFunction(n int64) int {
	if n == 1 {
		return 1
	}

	primeFactorCount := 0
	for i := int64(2); i <= int64(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			// Check if i^2 divides n
			if (n/(i*i))*i*i == n {
				return 0
			}
			primeFactorCount++
			n /= i
		}
	}

	// If n is still greater than 1, it is a prime factor
	if n > 1 {
		primeFactorCount++
	}

	// Return (-1)^primeFactorCount
	if primeFactorCount%2 == 0 {
		return 1
	}
	return -1
}
