package main

import (
	"config"
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"liberdatabase"
	"math/big"
	"sequences"
	"sync"
	"time"
)

// findCombos finds prime combos for a given number.
func findCombos(db *gorm.DB, mainId string, n *big.Int, pmax int) bool {
	number := new(big.Int).Set(n)
	seqNumber := int64(0)
	loopCounter := int64(0)

	// Get p values
	fmt.Println("Getting possible p values")
	getPValues(mainId, number, pmax)

	// Initialize the last sequence number
	var lastSeqNumber int64 = 0

	// Loop to get factors until nil is returned
	for {
		loopCounter++
		factor := liberdatabase.GetFactorsByMainID(db, mainId, lastSeqNumber)
		if factor == nil {
			break
		}

		// Update the last sequence number
		lastSeqNumber = factor.SeqNumber

		// Convert the factor to a big.Int
		prime := new(big.Int)
		if _, ok := prime.SetString(factor.Factor, 10); !ok {
			fmt.Printf("Error converting factor to big.Int: %s\n", factor.Factor)
			continue
		}

		if loopCounter == 1000000 {
			fmt.Printf("Current prime at loop %d: %s\n", loopCounter, factor.Factor)
			loopCounter = 0 // Reset loopCounter
		}

		fmt.Println("Processing factor: ", factor.Factor)

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

			fmt.Println("Found prime p,q factors: ", combo.ValueP, combo.ValueQ)

			liberdatabase.InsertPrimeCombo(db, combo)
		}

	}

	liberdatabase.RemoveFactorsByMainID(db, mainId)

	return true
}

// getPValues finds p values using multiple workers.
func getPValues(mainId string, n *big.Int, pmax int) {
	pcount := 0

	// Load worker count from config
	cfg, err := config.LoadConfig()
	workerCount := 4 // Default worker count
	if err != nil {
		fmt.Printf("Error loading config: %v\nUsing default worker count: %d\n", err, workerCount)
	} else {
		workerCount = cfg.NumWorkers / 2
	}

	// Create channels for distributing work and collecting results
	primeChan := make(chan *big.Int)
	resultChan := make(chan *big.Int)
	var wg sync.WaitGroup

	// Create a context to handle cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start worker goroutines
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for {
				if pcount >= pmax {
					cancel() // Cancel the context to stop the workers
					break
				}

				select {
				case <-ctx.Done():
					return
				case prime, ok := <-primeChan:
					if !ok {
						return
					}
					if new(big.Int).Mod(n, prime).Cmp(big.NewInt(0)) == 0 {
						select {
						case <-ctx.Done():
							return
						case resultChan <- prime:
						}
					}
				}
			}
		}(i)
	}

	// Start a goroutine to close the result channel once all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	primeCount := 0
	processedNumber := big.NewInt(0)
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	go func() {
		colors := []string{"\033[31m", "\033[32m", "\033[33m", "\033[34m", "\033[35m", "\033[36m", "\033[37m", "\033[90m", "\033[91m", "\033[92m"}
		colorIndex := 0
		for range ticker.C {
			aps := primeCount
			fmt.Printf("%s - %s Primes per minute: %d - Primes Processed: %s \033[0m\n", n.String(), colors[colorIndex], aps, processedNumber.String())
			primeCount = 0
			colorIndex = (colorIndex + 1) % len(colors)
		}
	}()

	// Start a goroutine to send primes to the workers
	go func() {
		for prime := range sequences.YieldPrimesAsc(n) {
			primeCount++
			processedNumber.Add(processedNumber, big.NewInt(1))
			if pcount >= pmax {
				cancel() // Cancel the context to stop the workers
				break
			}
			select {
			case <-ctx.Done():
				close(primeChan)
				return
			case primeChan <- prime:
			}
		}
		close(primeChan)
	}()

	seqValue := int64(0)

	// Collect results
	for prime := range resultChan {
		if pcount >= pmax {
			cancel() // Cancel the context to stop the workers
			break
		}
		seqValue++
		pcount++
		fmt.Printf("Found prime factor: %s\n", prime.String())
		// Insert the prime into the database or perform other actions
		factor := liberdatabase.Factor{
			ID:        uuid.New().String(),
			Factor:    prime.String(),
			MainId:    mainId,
			SeqNumber: seqValue,
		}

		// Create a new database connection for inserting the factor
		db, err := liberdatabase.InitConnection()
		if err != nil {
			fmt.Printf("Error initializing database connection: %v\n", err)
			continue
		}
		liberdatabase.InsertFactor(db, factor)
	}
}
