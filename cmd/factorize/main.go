package main

import (
	"flag"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"liberdatabase"
	"math/big"
	"os"
	"sequences"
	"strings"
	"time"
	"titler"
)

var status strings.Builder

func main() {
	titler.PrintTitle("Factorize")
	startTime := time.Now() // Record the start time

	// Timer to write out the status every minute
	ticker := time.NewTicker(1 * time.Minute)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				if status.Len() > 0 {
					fmt.Println("Status update:", status.String())
				}
			case <-done:
				ticker.Stop()
				return
			}
		}
	}()
	// Make sure to stop the ticker when the program ends
	defer func() {
		done <- true
	}()

	// Parse command-line flags
	flag.Parse()

	// Check if the number is provided as an argument
	if len(flag.Args()) < 1 {
		fmt.Println("Please provide a number to be factorized as an argument.")
		os.Exit(1)
	}

	// Read the input number
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

	if sequences.IsPrime(number) {
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
	writeOutputToFile(fmt.Sprintf("%s : %s", numberStr, output.String()))

	liberdatabase.RemoveFactorsByMainID(db, mainId)

	endTime := time.Now()                        // Record the end time
	duration := endTime.Sub(startTime)           // Calculate the duration
	fmt.Printf("Execution time: %v\n", duration) // Print the log message
}

// factorize returns the prime factors of a given big integer.
func factorize(db *gorm.DB, mainId string, n *big.Int, lastSeq int64) bool {
	counter := big.NewInt(2)
	zero := big.NewInt(0)
	number := new(big.Int).Set(n)

	fmt.Printf("- Step - Factoring %s (%d bits)\n", number.String(), number.BitLen())

	if lastSeq > 0 {
		lastRecord := liberdatabase.GetMaxSeqNumberByMainID(db, mainId)
		liberdatabase.RemoveFactorByID(db, lastRecord.ID)
	}

	if number.ProbablyPrime(20) {
		fmt.Printf("-%s is prime\n", number.String())
	}

	// Check if n is divisible by x
	for counter.Cmp(number) <= 0 {
		status.Reset()
		status.WriteString(counter.String())

		if new(big.Int).Mod(number, counter).Cmp(zero) == 0 {
			fmt.Printf("- Factor %s found\n", counter.String())

			number = n.Div(number, counter)

			// Insert the count factor into the database
			lastSeq++
			counterFactor := liberdatabase.Factor{
				ID:        uuid.New().String(),
				Factor:    counter.String(),
				MainId:    mainId,
				SeqNumber: lastSeq,
			}

			liberdatabase.InsertFactor(db, counterFactor)

			// Insert the number factor into the database
			lastSeq++
			numberFactor := liberdatabase.Factor{
				ID:        uuid.New().String(),
				Factor:    number.String(),
				MainId:    mainId,
				SeqNumber: lastSeq,
			}

			liberdatabase.InsertFactor(db, numberFactor)
			break
		} else {
			counter.Add(counter, big.NewInt(1))
		}
	}

	// Check if all factors are prime
	areAllPrime := true
	lastSeqNumber := int64(0)

	// Loop to get factors until nil is returned
	for {
		factor := liberdatabase.GetFactorsByMainID(db, mainId, lastSeqNumber)
		if factor == nil {
			break
		}

		// Update the last sequence number
		lastSeqNumber = factor.SeqNumber

		factorValue := new(big.Int)
		factorValue, ok := factorValue.SetString(factor.Factor, 10)
		if !ok {
			fmt.Printf("Error converting factor to *big.Int: %v\n", factor.Factor)
		}

		if !factorValue.ProbablyPrime(20) {
			areAllPrime = false
			break
		}
	}

	if areAllPrime {
		return true
	} else {
		return factorize(db, mainId, number, lastSeq)
	}
}

func writeOutputToFile(output string) {
	file, err := os.OpenFile("factorize_output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer func(file *os.File) {
		closeError := file.Close()
		if err != nil {
			fmt.Printf("Error closing file: %v\n", closeError)
		}
	}(file)

	if _, err := file.WriteString(output + "\n"); err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	}
}
