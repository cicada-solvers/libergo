package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"
)

const base60Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// Base60ToBase10 converts a base60 string to a base10 *big.Int.
func Base60ToBase10(base60 string) (*big.Int, error) {
	result := big.NewInt(0)
	base := big.NewInt(60)
	for _, char := range base60 {
		index := strings.IndexRune(base60Chars, char)
		if index == -1 {
			return nil, fmt.Errorf("invalid character '%c' in base60 string", char)
		}
		result.Mul(result, base)
		result.Add(result, big.NewInt(int64(index)))
	}
	return result, nil
}

// Base10ToBase60 converts a base10 *big.Int to a base60 string.
func Base10ToBase60(base10 *big.Int) string {
	if base10.Cmp(big.NewInt(0)) == 0 {
		return "0"
	}

	result := ""
	base := big.NewInt(60)
	zero := big.NewInt(0)
	mod := new(big.Int)

	for base10.Cmp(zero) > 0 {
		base10.DivMod(base10, base, mod)
		result = string(base60Chars[mod.Int64()]) + result
	}

	return result
}

// handleError prints the error message and exits the program.
func handleError(err error) {
	fmt.Println("Error:", err)
	os.Exit(1)
}

func main() {
	base60ToBase10 := flag.String("base60ToBase10", "", "Convert base60 to base10")
	base10ToBase60 := flag.String("base10ToBase60", "", "Convert base10 to base60")
	flag.Parse()

	if *base60ToBase10 != "" {
		base60Numbers := strings.Split(*base60ToBase10, ",")
		var base10Results []string
		for _, base60Number := range base60Numbers {
			result, err := Base60ToBase10(strings.TrimSpace(base60Number))
			if err != nil {
				handleError(err)
			}
			base10Results = append(base10Results, result.String())
		}
		fmt.Println(strings.Join(base10Results, "\n"))
	} else if *base10ToBase60 != "" {
		base10Numbers := strings.Split(*base10ToBase60, ",")
		var base60Results []string
		for _, base10Number := range base10Numbers {
			base10 := new(big.Int)
			if _, ok := base10.SetString(strings.TrimSpace(base10Number), 10); !ok {
				handleError(fmt.Errorf("invalid base10 number"))
			}
			result := Base10ToBase60(base10)
			base60Results = append(base60Results, result)
		}
		fmt.Println(strings.Join(base60Results, "\n"))
	} else {
		fmt.Println("Please provide a flag to convert either from base60 to base10 or base10 to base60")
		flag.Usage()
		os.Exit(1)
	}
}
