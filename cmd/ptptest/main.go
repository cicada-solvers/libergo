package main

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
	liteConn, err := liberdatabase.InitSQLiteConnection()

	nonPrimeCount := big.NewInt(int64(0))

	for i := big.NewInt(1); i.Cmp(big.NewInt(math.MaxInt64)) <= 0; i.Add(i, big.NewInt(1)) {
		if sequences.IsPrime(i) {
			fmt.Printf("%s is prime\n", i.String())

			record := liberdatabase.PrimeNumRecord{
				Number:                 i.String(),
				NumberCountBeforePrime: nonPrimeCount.String(),
				NumberIsPrime:          true,
				NumberFactorSize:       int64(2),
			}

			addErr := liberdatabase.AddPrimeNumRecord(conn, record)
			if addErr != nil {
				// Handle the error
				fmt.Printf("Error adding prime number: %v\n", addErr)
			}

			nonPrimeCount.SetInt64(int64(0))
		} else {
			fmt.Printf("%s is not prime\n", i.String())
			n := new(big.Int).Set(i)
			factorsize := factorize(liteConn, uuid.New().String(), n, 0)

			record := liberdatabase.PrimeNumRecord{
				Number:                 i.String(),
				NumberCountBeforePrime: nonPrimeCount.String(),
				NumberIsPrime:          false,
				NumberFactorSize:       factorsize,
			}

			addErr := liberdatabase.AddPrimeNumRecord(conn, record)
			if addErr != nil {
				// Handle the error
				fmt.Printf("Error adding prime number: %v\n", addErr)
			}

			nonPrimeCount.Add(nonPrimeCount, big.NewInt(1))
		}
	}
}

func factorize(db *gorm.DB, mainId string, n *big.Int, lastSeq int64) int64 {
	counter := big.NewInt(2)
	zero := big.NewInt(0)
	number := new(big.Int).Set(n)

	if lastSeq > 0 {
		lastRecord := liberdatabase.GetMaxSeqNumberByMainID(db, mainId)
		liberdatabase.RemoveFactorByID(db, lastRecord.ID)
	}

	// Check if n is divisible by x
	for counter.Cmp(number) <= 0 {
		if new(big.Int).Mod(number, counter).Cmp(zero) == 0 {
			number = n.Div(number, counter)

			// Insert the counter factor into the database
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
	counter = big.NewInt(0)
	for {
		counter.Add(counter, big.NewInt(1))
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
		return counter.Int64()
	} else {
		return factorize(db, mainId, number, lastSeq)
	}
}
