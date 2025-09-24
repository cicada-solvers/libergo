package main

import (
	"fmt"
	"liberdatabase"
	"math"
	"math/big"
	"sequences"
	"strings"
)

func main() {
	_, _ = liberdatabase.InitTables()

	db, _ := liberdatabase.InitConnection()

	collatzPrimes := make([]liberdatabase.CollatzPrime, 0)

	for i := int64(2); i <= math.MaxInt32; i++ {
		fmt.Printf("\rProcessing %d", i)
		var sequence *sequences.NumericSequence

		sequence, _ = sequences.GetCollatzSequence(i, false)
		isPrime := sequences.IsPrime(big.NewInt(i))
		isCollatzLengthPrime := sequences.IsPrime(big.NewInt(int64(len(sequence.Sequence))))

		var sequenceStrings []string
		for _, num := range sequence.Sequence {
			sequenceStrings = append(sequenceStrings, num.String())
		}

		collatzPrime := liberdatabase.CollatzPrime{
			CollatzSequence:           strings.Join(sequenceStrings, ","),
			IsPrime:                   isPrime,
			CollatzPrimeLength:        int64(len(sequence.Sequence)),
			Number:                    i,
			CollatzPrimeLengthIsPrime: isCollatzLengthPrime,
		}

		collatzPrimes = append(collatzPrimes, collatzPrime)
		if len(collatzPrimes) == 500 {
			liberdatabase.AddCollatzPrimes(db, collatzPrimes)
			collatzPrimes = make([]liberdatabase.CollatzPrime, 0)
		}
	}

	if len(collatzPrimes) > 0 {
		liberdatabase.AddCollatzPrimes(db, collatzPrimes)
	}

	_ = liberdatabase.CloseConnection(db)

	fmt.Printf("\n")
}
