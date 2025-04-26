package main

import (
	"fmt"
	"liberdatabase"
	"math"
	"math/big"
	"sequences"
)

func main() {
	err := liberdatabase.InitMySqlTables()
	if err != nil {
		fmt.Printf("Error initializing MySQL tables: %v\n", err)
		return
	}

	// Initialize the database connection
	conn, err := liberdatabase.InitMySQLConnection()

	nonPrimeCount := big.NewInt(int64(0))

	for i := big.NewInt(2); i.Cmp(big.NewInt(math.MaxInt64)) <= 0; i.Add(i, big.NewInt(1)) {
		if sequences.IsPrime(i) {
			fmt.Printf("%s is prime\n", i.String())

			record := liberdatabase.PrimeNumRecord{
				Number:                 i.Int64(),
				NumberCountBeforePrime: int(nonPrimeCount.Int64()),
			}

			addErr := liberdatabase.AddPrimeNumRecord(conn, record)
			if addErr != nil {
				// Handle the error
				fmt.Printf("Error adding prime number: %v\n", addErr)
			}

			nonPrimeCount.SetInt64(int64(0))
		} else {
			nonPrimeCount.Add(nonPrimeCount, big.NewInt(1))
		}
	}
}
