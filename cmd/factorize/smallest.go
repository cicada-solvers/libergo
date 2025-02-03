package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"liberdatabase"
	"math/big"
)

// factorize returns the prime factors of a given big integer.
func factorize(db *pgx.Conn, mainId string, n *big.Int, lastSeq int64) bool {
	counter := big.NewInt(2)
	zero := big.NewInt(0)
	number := new(big.Int).Set(n)

	if lastSeq > 0 {
		lastRecord, dbErr := liberdatabase.GetMaxSeqNumberByMainID(db, mainId)
		if dbErr != nil {
			fmt.Printf("Error getting max sequence number: %v\n", dbErr)
		}

		err := liberdatabase.RemoveFactorByID(db, lastRecord.ID)
		if err != nil {
			fmt.Printf("Error getting max sequence number: %v\n", dbErr)
		}
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

			err := liberdatabase.InsertFactor(db, counterFactor)
			if err != nil {
				fmt.Printf("Error inserting factor: %v\n", err)
			}

			// Insert the number factor into the database
			lastSeq++
			numberFactor := liberdatabase.Factor{
				ID:        uuid.New().String(),
				Factor:    number.String(),
				MainId:    mainId,
				SeqNumber: lastSeq,
			}

			numberErr := liberdatabase.InsertFactor(db, numberFactor)
			if numberErr != nil {
				fmt.Printf("Error inserting factor: %v\n", numberErr)
			}

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
		factor, err := liberdatabase.GetFactorsByMainID(db, mainId, lastSeqNumber)
		if err != nil {
			fmt.Printf("Error getting factors: %v\n", err)
		}
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
