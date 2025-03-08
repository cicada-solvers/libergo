package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"liberdatabase"
	"math/big"
)

// factorize returns the prime factors of a given big integer.
func factorize(db *gorm.DB, mainId string, n *big.Int, lastSeq int64, p *tea.Program) bool {
	counter := big.NewInt(2)
	zero := big.NewInt(0)
	number := new(big.Int).Set(n)

	if lastSeq > 0 {
		lastRecord := liberdatabase.GetMaxSeqNumberByMainID(db, mainId)
		liberdatabase.RemoveFactorByID(db, lastRecord.ID)
	}

	for counter.Cmp(number) <= 0 {
		p.Send(counterMsg{counter: new(big.Int).Set(counter)}) // Send the current counter to the Bubble Tea model

		stringVal := counter.String()
		if len(stringVal) > 1 {
			if stringVal[len(stringVal)-1:] == "5" ||
				stringVal[len(stringVal)-1:] == "0" ||
				stringVal[len(stringVal)-1:] == "2" ||
				stringVal[len(stringVal)-1:] == "4" ||
				stringVal[len(stringVal)-1:] == "6" ||
				stringVal[len(stringVal)-1:] == "8" {
				counter.Add(counter, big.NewInt(1))
				continue
			}
		}

		if new(big.Int).Mod(number, counter).Cmp(zero) == 0 {
			number = n.Div(number, counter)

			lastSeq++
			counterFactor := liberdatabase.Factor{
				ID:        uuid.New().String(),
				Factor:    counter.String(),
				MainId:    mainId,
				SeqNumber: lastSeq,
			}

			liberdatabase.InsertFactor(db, counterFactor)

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

	areAllPrime := true
	lastSeqNumber := int64(0)

	for {
		factor := liberdatabase.GetFactorsByMainID(db, mainId, lastSeqNumber)
		if factor == nil {
			break
		}

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
		return factorize(db, mainId, number, lastSeq, p)
	}
}
