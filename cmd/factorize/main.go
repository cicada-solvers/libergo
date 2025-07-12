package main

import (
	"flag"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"liberdatabase"
	"math/big"
	"os"
	"runtime"
	"sequences"
	"sort"
	"strings"
	"sync"
	"time"
	"titler"
)

// NumberToCheck represents information for verifying divisors of a number during factorization.
type NumberToCheck struct {
	Number  string
	Counter string
	IsBound bool
}

// status represents the current status of a computation or process, stored as a big integer for large number handling.
var status big.Int

// main is the entry point of the program. It handles input parsing, initializes resources, and executes the factorization process.
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
				fmt.Printf("Status update: %s (bits %d)\n", status.String(), status.BitLen())
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
	zero := big.NewInt(0)
	number := new(big.Int).Set(n)
	var modNumberArray []big.Int
	processedCounter := big.NewInt(0)

	fmt.Printf("- Step - Factoring %s (%d bits)\n", number.String(), number.BitLen())

	if lastSeq > 0 {
		lastRecord := liberdatabase.GetMaxSeqNumberByMainID(db, mainId)
		liberdatabase.RemoveFactorByID(db, lastRecord.ID)
	}

	if number.ProbablyPrime(20) {
		fmt.Printf("-%s is prime\n", number.String())
		modNumberArray = append(modNumberArray, *number)
	} else {

		// We're going to use threads to check this out
		checkChannel := make(chan NumberToCheck)
		var wg sync.WaitGroup
		numProcessors := runtime.NumCPU()
		waits := 0
		if number.Cmp(big.NewInt(int64(numProcessors*4))) > 0 {
			waits = numProcessors * 4
		} else {
			waits = numProcessors
		}

		// Start worker goroutines
		for i := 0; i < numProcessors*2; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for checkValue := range checkChannel {
					myBigNumber, _ := new(big.Int).SetString(checkValue.Number, 10)
					myBigCounter, _ := new(big.Int).SetString(checkValue.Counter, 10)

					if myBigNumber.Cmp(myBigCounter) <= 0 {
						fmt.Printf("Error: Number %s is less than or equal to counter %s\n", myBigNumber.String(), myBigCounter.String())
						continue
					}

					if new(big.Int).Mod(myBigNumber, myBigCounter).Cmp(zero) == 0 {
						modNumberArray = append(modNumberArray, *myBigCounter)

						if checkValue.IsBound {
							fmt.Printf("- %s Factor found in bound.\n", myBigCounter.String())
						}

						break
					}
				}
			}()
		}

		go func() {
			one := big.NewInt(1)
			two := big.NewInt(2)
			myCounter := new(big.Int).Set(two)
			var counterString []string

			// If the number is greater than 1000, then we are going to check the percent
			// bounds.  Most factors are going to be within 6% of the square root
			// or 2% of 2.
			if number.Cmp(big.NewInt(1000)) >= 0 {
				twoPercentLowerBound, twoPercentUpperBound, _ := getTwoPercentFromSquareRoot(number)
				myCounter = new(big.Int).Set(twoPercentLowerBound)

				for myCounter.Cmp(twoPercentUpperBound) <= 0 {
					counterString = strings.Split(myCounter.String(), "")
					if len(counterString) >= 2 {
						switch counterString[len(counterString)-1] {
						case "0", "2", "4", "5", "6", "8":
							myCounter.Add(myCounter, one)
							continue
						default:
							// nothing
						}
					}

					status.Set(myCounter)

					checkChannel <- NumberToCheck{
						Number:  number.String(),
						Counter: myCounter.String(),
						IsBound: true,
					}

					myCounter.Add(myCounter, one)

					if len(modNumberArray) > 0 && processedCounter.Cmp(big.NewInt(int64(waits))) > 0 {
						break
					}

					processedCounter.Add(processedCounter, one)
				}

				sixPercentLowerBound, sixPercentUpperBound, _ := getSixPercentFromSquareRoot(number)
				myCounter = new(big.Int).Set(sixPercentLowerBound)

				for myCounter.Cmp(sixPercentUpperBound) <= 0 {
					counterString = strings.Split(myCounter.String(), "")
					if len(counterString) >= 2 {
						switch counterString[len(counterString)-1] {
						case "0", "2", "4", "5", "6", "8":
							myCounter.Add(myCounter, one)
							continue
						default:
							// nothing
						}
					}

					status.Set(myCounter)

					checkChannel <- NumberToCheck{
						Number:  number.String(),
						Counter: myCounter.String(),
						IsBound: true,
					}

					myCounter.Add(myCounter, one)

					if len(modNumberArray) > 0 && processedCounter.Cmp(big.NewInt(int64(waits))) > 0 {
						break
					}

					processedCounter.Add(processedCounter, one)
				}
			}

			myCounter = new(big.Int).Set(two)
			for myCounter.Cmp(number) <= 0 {
				counterString = strings.Split(myCounter.String(), "")
				if len(counterString) >= 2 {
					switch counterString[len(counterString)-1] {
					case "0", "2", "4", "5", "6", "8":
						myCounter.Add(myCounter, one)
						continue
					default:
						// nothing
					}
				}

				status.Set(myCounter)

				checkChannel <- NumberToCheck{
					Number:  number.String(),
					Counter: myCounter.String(),
					IsBound: false,
				}

				myCounter.Add(myCounter, one)

				if len(modNumberArray) > 0 && processedCounter.Cmp(big.NewInt(int64(waits))) > 0 {
					break
				}

				processedCounter.Add(processedCounter, one)
			}
			close(checkChannel)
		}()

		wg.Wait()

		// grabbing the smallest factor
		modNumberArray = sortBigInts(modNumberArray)
	}

	modArrayFirst := modNumberArray[0]

	fmt.Printf("- Factor %s found\n", modArrayFirst.String())

	number = n.Div(number, &modArrayFirst)

	// Insert the count factor into the database
	lastSeq++
	counterFactor := liberdatabase.Factor{
		ID:        uuid.New().String(),
		Factor:    modArrayFirst.String(),
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

// sortBigInts sorts a slice of big.Int in ascending order and returns the sorted slice.
func sortBigInts(bigInts []big.Int) []big.Int {
	sort.Slice(bigInts, func(i, j int) bool {
		return bigInts[i].Cmp(&bigInts[j]) < 0
	})
	return bigInts
}

// writeOutputToFile appends the given output string to a file named "factorize_output.txt".
func writeOutputToFile(output string) {
	file, err := os.OpenFile("factorize_output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer func(file *os.File) {
		closeError := file.Close()
		if closeError != nil {
			fmt.Printf("Error closing file: %v\n", closeError)
		}
	}(file)

	if _, writeError := file.WriteString(output + "\n"); writeError != nil {
		fmt.Printf("Error writing to file: %v\n", writeError)
	}
}

// getSixPercentFromSquareRoot calculates square root of a number and determines bounds +/- 6% around the square root.
// It returns the lower bound, upper bound, and an error if any operation fails.
func getSixPercentFromSquareRoot(number *big.Int) (*big.Int, *big.Int, error) {
	sqrt := new(big.Int).Sqrt(number)

	// Convert to big.Int
	sixPercentInt := calculateSixPercent(sqrt)

	// Calculate bounds
	lowerBound := new(big.Int).Sub(sqrt, sixPercentInt)
	upperBound := new(big.Int).Add(sqrt, sixPercentInt)

	return lowerBound, upperBound, nil
}

// getTwoPercentFromSquareRoot calculates 2% of the input number and returns the lower and upper bounds around 2.
func getTwoPercentFromSquareRoot(number *big.Int) (*big.Int, *big.Int, error) {
	twoPercent := calculateTwoPercent(number)

	// Calculate the bounds
	lowerBound := big.NewInt(2)
	upperBound := new(big.Int).Add(lowerBound, twoPercent)

	return lowerBound, upperBound, nil
}

// calculateTwoPercent calculates 2% of the given value rounded to the nearest integer using high precision arithmetic.
// Returns 1 if the resulting value is less than 1 and the original value is greater than 0.
func calculateTwoPercent(value *big.Int) *big.Int {
	// Method 2: Using big.Float for more precision
	valueFloat := new(big.Float).SetInt(value)
	percentFloat := new(big.Float).Mul(valueFloat, big.NewFloat(0.02))

	// Round to the nearest integer
	percentFloat.Add(percentFloat, big.NewFloat(0.5))

	// Convert back to big.Int
	twoPercent := new(big.Int)
	percentFloat.Int(twoPercent)

	// Ensure the minimum value of 1 if the original value is non-zero
	if value.Sign() > 0 && twoPercent.Sign() == 0 {
		twoPercent.SetInt64(1)
	}

	return twoPercent
}

// calculateSixPercent calculates 6% of the given *big.Int value and returns the result as *big.Int, ensuring precision.
// If the original value is non-zero, the function ensures a minimum result of 1. It rounds to the nearest integer.
func calculateSixPercent(value *big.Int) *big.Int {
	// Method 2: Using big.Float for more precision
	valueFloat := new(big.Float).SetInt(value)
	percentFloat := new(big.Float).Mul(valueFloat, big.NewFloat(0.06))

	// Round to the nearest integer
	percentFloat.Add(percentFloat, big.NewFloat(0.5))

	// Convert back to big.Int
	twoPercent := new(big.Int)
	percentFloat.Int(twoPercent)

	// Ensure the minimum value of 1 if the original value is non-zero
	if value.Sign() > 0 && twoPercent.Sign() == 0 {
		twoPercent.SetInt64(1)
	}

	return twoPercent
}
