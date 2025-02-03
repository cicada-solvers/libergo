package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"liberdatabase"
	"math/big"
	"sequences"
)

// findCombos finds prime combos for a given number.
func findCombos(db *pgx.Conn, mainId string, n *big.Int) bool {
	number := new(big.Int).Set(n)
	zero := big.NewInt(0)
	seqNumber := int64(0)
	loopCounter := int64(0)

	for prime := range sequences.YieldPrimesDesc(n) {
		loopCounter++

		if loopCounter == 1000000 {
			fmt.Printf("Current prime at loop %d: %s\n", loopCounter, prime.String())
			loopCounter = 0 // Reset loopCounter
		}

		if new(big.Int).Mod(number, prime).Cmp(zero) == 0 {
			q := new(big.Int).Div(number, prime)

			if q.ProbablyPrime(20) {
				seqNumber++

				// Insert the prime combo into the database
				combo := liberdatabase.PrimeCombo{
					ID:        uuid.New().String(),
					ValueP:    prime.String(),
					ValueQ:    q.String(),
					MainId:    mainId,
					SeqNumber: seqNumber,
				}

				fmt.Println("Found combo: ", combo.ValueP, combo.ValueQ)

				err := liberdatabase.InsertPrimeCombo(db, combo)
				if err != nil {
					fmt.Printf("Error inserting factor: %v\n", err)
					return false
				}
			}
		}
	}

	return true
}
