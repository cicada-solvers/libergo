package main

import (
	"fmt"
	"liberdatabase"
	"math"
	"math/big"
	"numeric"
	"runtime"
	"sequences"
	"sync"

	"gorm.io/gorm"
)

var conns []*gorm.DB

// main is the entry point of the program. It manages database setup, worker initialization, and Goldbach processing workflow.
func main() {
	// Initialize database
	_, _ = liberdatabase.InitTables()

	// Create a channel for numbers to be processed
	numberChannel := make(chan int64)

	// Determine the number of workers (CPU count × 2)
	numWorkers := runtime.NumCPU()
	conns = make([]*gorm.DB, numWorkers)
	fmt.Printf("Using %d worker goroutines\n", numWorkers)

	// Use WaitGroup to wait for all workers to finish
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		conns[i], _ = liberdatabase.InitConnection()
		go func(workerID int) {
			for num := range numberChannel {
				//fmt.Printf("Processing number %d\n", num)
				// Adding the initial number to the database.
				liberdatabase.AddGoldbachNumber(conns[workerID], num, numeric.IsNumberEven(num), sequences.IsPrime(big.NewInt(num)))

				// Add factors to the database
				gbs := numeric.NewGoldbachSets()
				gbs.Solve(num)
				if len(gbs.GoldBachSets) == 0 {
					continue
				}

				fmt.Printf("Number %d has %d Goldbach sets\n", num, len(gbs.GoldBachSets))
				var dbAddends []liberdatabase.GoldbachAddend
				for counter, addend := range gbs.GoldBachSets {
					GoldbachAddend := liberdatabase.GoldbachAddend{
						GoldbachNumber: num,
						AddendOne:      addend.AddendOne,
						AddendTwo:      addend.AddendTwo,
						AddendThree:    addend.AddendThree,
						SetNumber:      counter,
					}

					dbAddends = append(dbAddends, GoldbachAddend)

					if len(dbAddends) >= 500 {
						liberdatabase.AddGoldbachAddends(conns[workerID], dbAddends)
						dbAddends = []liberdatabase.GoldbachAddend{}
					}
				}

				liberdatabase.AddGoldbachAddends(conns[workerID], dbAddends)
			}

			wg.Done()
		}(i)
	}

	go func() {
		largestNumber := int64(math.MaxInt32)
		for i := int64(6); i <= largestNumber; i++ {
			if !numeric.IsNumberEven(i) {
				numberChannel <- i
			}
		}
		// Close the channel when done generating
		close(numberChannel)
		fmt.Println("Number generation completed")
	}()

	// Wait for all workers to finish
	wg.Wait()
	fmt.Println("All workers have completed")
}
