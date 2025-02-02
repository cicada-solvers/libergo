package main

import (
	"fmt"
	"math/big"
	"os"
	"sequences"
	"strings"
)

// factorize returns the prime factors of a given big integer.
func factorize(n *big.Int, existing []*big.Int) []*big.Int {
	counter := big.NewInt(2)
	zero := big.NewInt(0)
	number := new(big.Int).Set(n)

	if len(existing) > 0 {
		existing = existing[:len(existing)-1] // Remove the last item
	}

	// Check if n is divisible by x
	for counter.Cmp(number) <= 0 {
		if new(big.Int).Mod(number, counter).Cmp(zero) == 0 {
			number = n.Div(number, counter)
			existing = append(existing, counter)
			existing = append(existing, number)
			break
		} else {
			counter.Add(counter, big.NewInt(1))
		}
	}

	if areAllFactorsPrime(existing) {
		return existing
	} else {
		return factorize(n, existing)
	}
}

func areAllFactorsPrime(factors []*big.Int) bool {
	for _, factor := range factors {
		if !sequences.IsPrime(factor) {
			return false
		}
	}
	return true
}

func main() {
	// Check if the number is provided as an argument
	if len(os.Args) < 2 {
		fmt.Println("Please provide a number to be factorized as an argument.")
		os.Exit(1)
	}

	// Read input number
	numberStr := os.Args[1]

	// Convert input to bigint
	number := new(big.Int)
	_, ok := number.SetString(numberStr, 10)
	if !ok {
		fmt.Println("Invalid number format.")
		os.Exit(1)
	}

	if number.Cmp(big.NewInt(1)) == -1 || number.Cmp(big.NewInt(1)) == 0 {
		fmt.Printf("%s : 1,%s\n", numberStr, numberStr)
		return
	}

	if sequences.IsPrime(number) {
		// You don't need to factorize a prime number
		fmt.Printf("%s : 1,%s\n", numberStr, numberStr)
		return
	}

	// Perform factorization
	primeList := factorize(number, []*big.Int{})

	// Output prime factors
	output := strings.Builder{}
	for counter, factor := range primeList {
		if counter == 0 {
			output.WriteString(factor.String())
		} else {
			output.WriteString(",")
			output.WriteString(factor.String())
		}
	}

	fmt.Println(numberStr, ":", output.String())
}
