package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"titler"
)

func main() {
	titler.PrintTitle("Integer Generation")

	if len(os.Args) < 2 {
		fmt.Println("Please provide the number of bits as an argument.")
		return
	}

	bits, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Invalid number of bits: %v\n", err)
		return
	}

	// Calculate the bit length for each prime
	primeBits := bits / 2

	// Generate two random prime numbers with the specified number of bits
	p, err := generateRandomPrime(primeBits)
	if err != nil {
		fmt.Printf("Error generating first prime: %v\n", err)
		return
	}

	q, err := generateRandomPrime(primeBits)
	if err != nil {
		fmt.Printf("Error generating second prime: %v\n", err)
		return
	}

	// Calculate the product of the two primes
	product := new(big.Int).Mul(p, q)

	fmt.Printf("Generated %d-bit prime p: %s\n", primeBits, p.String())
	fmt.Printf("Generated %d-bit prime q: %s\n", primeBits, q.String())
	fmt.Printf("Product of the two primes: %s\n", product.String())
	fmt.Printf("Number of bits in the product: %d\n", numberOfBits(product))
}

func generateRandomPrime(bits int) (*big.Int, error) {
	maxNum := new(big.Int).Lsh(big.NewInt(1), uint(bits))
	for {
		n, err := rand.Int(rand.Reader, maxNum)
		if err != nil {
			return nil, fmt.Errorf("error generating random integer: %v", err)
		}

		// Check if the number is prime
		if n.ProbablyPrime(20) {
			return n, nil
		}
	}
}

// numberOfBits returns the number of bits in a big.Int value
func numberOfBits(n *big.Int) int {
	return n.BitLen()
}
