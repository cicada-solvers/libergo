package main

import (
	"fmt"
	"liberdatabase"
	"math"
	"math/big"
	"sequences"
	"time"
	"titler"
)

// This program initializes a MySQL database connection, creates necessary tables, and populates them with prime number records.
func main() {
	titler.PrintTitle("PTP Test")

	// Initialize the database connection
	tableErr := liberdatabase.InitMySqlTables()
	if tableErr != nil {
		fmt.Printf("Error initializing MySQL tables: %v\n", tableErr)
		return
	}

	conn, connErr := liberdatabase.InitMySQLConnection()
	if connErr != nil {
		fmt.Printf("Error initializing MySQL connection: %v\n", connErr)
		return
	}

	for i := 0; i < math.MaxInt32; i++ {
		number := new(big.Int).SetInt64(int64(i))
		number2 := new(big.Int).SetInt64(int64(i))

		// Measure time for IsPrime
		start := time.Now()
		isPrime := sequences.IsPrime(number)
		durationIsPrime := time.Since(start)

		// Measure time for IsPrimeUsing6k
		start = time.Now()
		isPtpPrime := IsPrimeUsing6k(number2)
		durationIsPtpPrime := time.Since(start)

		if isPrime || isPtpPrime {
			record := liberdatabase.PrimeNumRecord{
				Number:             number.Int64(),
				IsPrime:            isPrime,
				IsPrimeDuration:    durationIsPrime.Microseconds(),
				IsPtpPrime:         isPtpPrime,
				IsPtpPrimeDuration: durationIsPtpPrime.Microseconds(),
			}

			addError := liberdatabase.AddPrimeNumRecord(conn, record)
			if addError != nil {
				fmt.Printf("Error adding prime number record: %v\n", addError)
				return
			}
		}
	}
}

func IsPrimeUsing6k(number *big.Int) bool {
	if number.Cmp(big.NewInt(10)) >= 0 {
		lastChar := number.String()[len(number.String())-1]
		if lastChar == '0' || lastChar == '2' || lastChar == '4' || lastChar == '5' || lastChar == '6' || lastChar == '8' {
			return false
		}
	}

	if number.Cmp(big.NewInt(2)) < 0 {
		return false
	}
	if number.Cmp(big.NewInt(2)) == 0 || number.Cmp(big.NewInt(3)) == 0 {
		return true
	}
	if new(big.Int).Mod(number, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 ||
		new(big.Int).Mod(number, big.NewInt(3)).Cmp(big.NewInt(0)) == 0 {
		return false
	}

	// Start checking with 6k ± 1
	sqrt := new(big.Int).Sqrt(number)
	k := big.NewInt(5)
	for k.Cmp(sqrt) <= 0 {
		if new(big.Int).Mod(number, k).Cmp(big.NewInt(0)) == 0 ||
			new(big.Int).Mod(number, new(big.Int).Add(k, big.NewInt(2))).Cmp(big.NewInt(0)) == 0 {
			return false
		}
		k.Add(k, big.NewInt(6))
	}

	return true
}
