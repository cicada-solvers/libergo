package main

import (
	"bufio"
	"fmt"
	"gorm.io/gorm"
	"liberdatabase"
	"math/big"
	"os"
	"sort"
	"strings"
)

// main is the entry point of the application, responsible for file parsing, database interaction, and data processing.
func main() {
	// Check if the filename is provided as a command line argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filename>")
		os.Exit(1)
	}

	// Get the filename from command-line arguments
	filename := os.Args[1]

	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			fmt.Printf("Error closing file: %v\n", closeErr)
		}
	}(file)

	// Create db connection
	_, _ = liberdatabase.InitTables()
	db, openErr := liberdatabase.InitConnection()
	if openErr != nil {
		fmt.Printf("Error opening database connection: %v\n", openErr)
		os.Exit(1)
	}
	defer func(db *gorm.DB) {
		closeErr := liberdatabase.CloseConnection(db)
		if closeErr != nil {
			fmt.Printf("Error closing database connection: %v\n", closeErr)
		}
	}(db)

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	lineNum := 0

	// Process each line
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Split the line by colon
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			fmt.Printf("Line %d: Invalid format, expected 'number: factor1,factor2,...'\n", lineNum)
			continue
		}

		// Parse the number
		numberStr := strings.TrimSpace(parts[0])
		number := new(big.Int)
		_, success := number.SetString(numberStr, 10)
		if !success {
			fmt.Printf("Line %d: Invalid number '%s'\n", lineNum, numberStr)
			continue
		}

		// Get the square root of the number
		sqrt := new(big.Int)
		sqrt.Sqrt(number)

		// Parse the factors
		factorsStr := strings.TrimSpace(parts[1])
		factorsSlice := strings.Split(factorsStr, ",")

		factors := make([]*big.Int, 0, len(factorsSlice))
		validFactors := true

		for _, factorStr := range factorsSlice {
			factorStr = strings.TrimSpace(factorStr)
			factor := new(big.Int)
			_, setSuccess := factor.SetString(factorStr, 10)
			if !setSuccess {
				fmt.Printf("Line %d: Invalid factor '%s'\n", lineNum, factorStr)
				validFactors = false
				break
			}
			factors = append(factors, factor)
		}

		factors = sortFactorsDesc(factors)

		if !validFactors {
			continue
		}

		// Verify that the factors are correct
		if verifyFactors(number, factors) {
			fmt.Printf("Line %d: %s has correct factors\n", lineNum, number.String())
			fmt.Printf("Line %d: factors: %v\n", lineNum, factors)
		} else {
			fmt.Printf("Line %d: %s has INCORRECT factors\n", lineNum, number.String())
			continue
		}

		// Add verified numbers to the database
		numberInformation := liberdatabase.AddAdvancedNumberInformation(db, number.Int64(), sqrt.Int64())
		for counter, factor := range factors {
			percentFromSquareRoot := getDistancePercentage(number, factor, sqrt)
			percentFromTwo := getDistancePercentage(number, factor, big.NewInt(2))
			percentFromMiddle := getDistancePercentage(number, factor, big.NewInt(0).Set(number).Div(number, big.NewInt(2)))
			percentFromNumber := getDistancePercentage(number, factor, big.NewInt(0).Set(number))

			liberdatabase.AddAdvancedNumberFactors(db, numberInformation.Id, factor.Int64(), counter, percentFromSquareRoot,
				percentFromNumber, percentFromTwo, percentFromMiddle)
		}
	}

	// Check for scanner errors
	if scanError := scanner.Err(); scanError != nil {
		fmt.Printf("Error reading file: %v\n", scanError)
	}
}

// verifyFactors checks if the given factors are correct for the number
func verifyFactors(number *big.Int, factors []*big.Int) bool {
	// Calculate product of all factors
	product := big.NewInt(1)
	for _, factor := range factors {
		product.Mul(product, factor)
	}

	// Check if product equals the number
	return product.Cmp(number) == 0
}

// sortFactorsDesc sorts a slice of *big.Int in descending order based on their values.
func sortFactorsDesc(factors []*big.Int) []*big.Int {
	sort.Slice(factors, func(i, j int) bool {
		return factors[i].Cmp(factors[j]) > 0
	})
	return factors
}

// getDistancePercentage calculates the percentage distance between two numbers based on a given reference number.
// It utilizes arbitrary-precision integers for the calculation to handle large values accurately.
func getDistancePercentage(number *big.Int, numberOne *big.Int, numberTwo *big.Int) float64 {
	// Calculate the absolute difference between numberOne and numberTwo
	diff := new(big.Int).Sub(numberOne, numberTwo)
	if diff.Sign() < 0 {
		diff.Neg(diff) // Get absolute value
	}

	// Convert to float64 for percentage calculation
	diffFloat := new(big.Float).SetInt(diff)
	numberFloat := new(big.Float).SetInt(number)

	// Calculate percentage (diff/number * 100)
	result := new(big.Float).Quo(diffFloat, numberFloat)
	result.Mul(result, big.NewFloat(100))

	// Convert to float64 for return
	percentage, _ := result.Float64()
	return percentage
}
