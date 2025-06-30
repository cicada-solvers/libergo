package main

import (
	"fmt"
	"gorm.io/gorm"
	"liberdatabase"
	"numeric"
	"runtime"
	"sync"
)

func main() {
	// Initialize database
	_, _ = liberdatabase.InitTables()
	conn, _ := liberdatabase.InitConnection()
	defer func(db *gorm.DB) {
		err := liberdatabase.CloseConnection(db)
		if err != nil {
			fmt.Println("Error closing database connection:", err)
		}
	}(conn)

	// Create a channel for numbers to be processed
	numberChannel := make(chan int64)

	// Determine the number of workers (CPU count × 2)
	numWorkers := runtime.NumCPU() * 2
	fmt.Printf("Using %d worker goroutines\n", numWorkers)

	// Use WaitGroup to wait for all workers to finish
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			for num := range numberChannel {
				gbp := numeric.NewGoldbachPairs()
				primeList := liberdatabase.GetPrimeListLessThanValue(conn, num)

				solveError := gbp.SolveForNumber(num, &primeList)
				if solveError != nil {
					fmt.Println("Error solving for number:", num, solveError)
					fmt.Println("----------------------------------------")
				} else {
					pairCount := len(gbp.GetGoldbachPairs())
					fmt.Println("Number:", num, pairCount)
					fmt.Println("----------------------------------------")

					// Now we need to add them to the database
					goldbachNumber := liberdatabase.AddGoldbachNumber(conn, num, true)
					for _, pair := range gbp.GetGoldbachPairs() {
						liberdatabase.AddGoldbachAddend(conn, goldbachNumber.Id, pair.AddendOne, pair.AddendTwo)
					}
				}
			}
		}(i)
	}

	go func() {
		largestNumber := int64(2147483647) + int64(2147483647)
		for i := int64(4); i <= largestNumber; i++ {
			if numeric.IsNumberEven(i) {
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
