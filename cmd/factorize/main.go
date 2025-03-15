package main

import (
	"flag"
	"fmt"
	"github.com/google/uuid"
	"liberdatabase"
	"math/big"
	"os"
	"strings"
	"time"
	"titler"
)

func main() {
	titler.PrintTitle("Factorize")
	startTime := time.Now() // Record the start time

	// Parse command-line flags
	flag.Parse()

	// Check if the number is provided as an argument
	if len(flag.Args()) < 1 {
		fmt.Println("Please provide a number to be factorized as an argument.")
		os.Exit(1)
	}

	// Read input number
	numberStr := flag.Arg(0)

	// Convert input to bigint
	number := new(big.Int)
	number, ok := number.SetString(numberStr, 10)
	if !ok {
		fmt.Println("Invalid number format.")
		os.Exit(1)
	}

	if number.Cmp(big.NewInt(1)) == -1 || number.Cmp(big.NewInt(1)) == 0 {
		fmt.Printf("%s : 1,%s\n", numberStr, numberStr)
		return
	}

	if number.ProbablyPrime(20) {
		// You don't need to factorize a prime number
		fmt.Printf("%s : 1,%s\n", numberStr, numberStr)
		return
	}

	// Establish database connection
	cError := liberdatabase.InitSQLiteTables()
	if cError != nil {
		fmt.Printf("Error creating database tables: %v\n", cError)
		os.Exit(1)
	}

	db, connError := liberdatabase.InitSQLiteConnection()
	if connError != nil {
		fmt.Printf("Error initializing database connection: %v\n", connError)
		os.Exit(1)
	}

	// The mainId is the number being factorized
	mainId := uuid.New().String()

	fmt.Printf("Factorizing %s (%d bits)\n", number.String(), number.BitLen())

	// Perform factorization
	factorize(db, mainId, number, 0)

	// Output prime factors
	output := strings.Builder{}
	firstTime := true

	// Initialize the last sequence number
	var lastSeqNumber int64 = 0

	// Loop to get factors until nil is returned
	for {
		factor := liberdatabase.GetFactorsByMainID(db, mainId, lastSeqNumber)
		if factor == nil {
			break
		}

		// Update the last sequence number
		lastSeqNumber = factor.SeqNumber

		if !firstTime {
			output.WriteString(",")
		}

		// Append factor to output
		output.WriteString(factor.Factor)

		firstTime = false
	}

	fmt.Println(numberStr, ":", output.String())

	liberdatabase.RemoveFactorsByMainID(db, mainId)

	endTime := time.Now()                        // Record the end time
	duration := endTime.Sub(startTime)           // Calculate the duration
	fmt.Printf("Execution time: %v\n", duration) // Print the log message
}
