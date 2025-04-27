package main

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"liberdatabase"
	"math"
	"math/big"
	"sequences"
	"strings"
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

	for i := big.NewInt(1); i.Cmp(big.NewInt(math.MaxInt32)) <= 0; i.Add(i, big.NewInt(1)) {
		if sequences.IsPrime(i) {
			modValue := new(big.Int).Set(i)

			record := liberdatabase.PrimeNumRecord{
				Number:                 i.Int64(),
				IsPrime:                true,
				NumberCountBeforePrime: nonPrimeCount.Int64(),
				PrimeFactorCount:       int64(2),
				PrimeFactors:           fmt.Sprintf("1,%s", i.String()),
				ModTwoTen:              modValue.Mod(modValue, big.NewInt(10)).Int64(),
			}

			addErr := liberdatabase.AddPrimeNumRecord(conn, record)
			if addErr != nil {
				// Handle the error
				fmt.Printf("Error adding prime number: %v\n", addErr)
			}

			nonPrimeCount.SetInt64(int64(0))
		} else {
			n := new(big.Int).Set(i)
			modValue := new(big.Int).Set(i)

			factors := factorize(liteConn, uuid.New().String(), n, 0)

			var factorStrings []string
			for _, factor := range factors {
				factorStrings = append(factorStrings, factor.String())
			}

			record := liberdatabase.PrimeNumRecord{
				Number:                 i.Int64(),
				IsPrime:                false,
				NumberCountBeforePrime: nonPrimeCount.Int64(),
				PrimeFactorCount:       int64(len(factors)),
				PrimeFactors:           strings.Join(factorStrings, ","),
				ModTwoTen:              modValue.Mod(modValue, big.NewInt(210)).Int64(),
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

// factorize returns the prime factors of a given big integer.
func factorize(db *gorm.DB, mainId string, n *big.Int, lastSeq int64) []big.Int {
	counter := big.NewInt(2)
	zero := big.NewInt(0)
	number := new(big.Int).Set(n)

	if lastSeq > 0 {
		lastRecord := liberdatabase.GetMaxSeqNumberByMainID(db, mainId)
		liberdatabase.RemoveFactorByID(db, lastRecord.ID)
	}

	if number.ProbablyPrime(20) {
		fmt.Printf("-%s is prime\n", number.String())
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
	var returnedFactors []big.Int
	for {
		factor := liberdatabase.GetFactorsByMainID(db, mainId, lastSeqNumber)
		if factor == nil {
			break
		}

		// Update the last sequence number
		lastSeqNumber = factor.SeqNumber

		factorValue := new(big.Int)
		factorValue, ok := factorValue.SetString(factor.Factor, 10)
		returnedFactors = append(returnedFactors, *factorValue)
		if !ok {
			fmt.Printf("Error converting factor to *big.Int: %v\n", factor.Factor)
		}

		if !factorValue.ProbablyPrime(20) {
			areAllPrime = false
			break
		}
	}

	if areAllPrime {
		return returnedFactors
	} else {
		return factorize(db, mainId, number, lastSeq)
	}
}
