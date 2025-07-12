package main

import (
	"fmt"
	"gorm.io/gorm"
	"liberdatabase"
	"math/big"
	"runtime"
	"sequences"
	"sync"
)

// main orchestrates the initialization of the system, worker creation, number generation, and prime number processing.
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
	numberChannel := make(chan string)

	// Determine the number of workers (CPU count × 2)
	numWorkers := runtime.NumCPU() * 2
	fmt.Printf("Using %d worker goroutines\n", numWorkers)

	// Use WaitGroup to wait for all workers to finish
	var wg sync.WaitGroup

	// Start the workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			// Each worker processes numbers from the channel
			for num := range numberChannel {
				// Convert int64 to big.Int for compatibility with IsPrime
				bigNum, _ := big.NewInt(0).SetString(num, 10)
				if sequences.IsPrime(bigNum) {
					liberdatabase.AddPrimeValue(conn, num, bigNum.BitLen())
				}
			}
			fmt.Printf("Worker %d completed\n", workerID)
		}(i)
	}

	// Start the number generator in a separate goroutine
	go func() {
		// Generate numbers from 2 to MaxInt32
		for i := big.NewInt(2); i.BitLen() <= 2048; i.Add(i, big.NewInt(1)) {
			numberChannel <- i.String()
		}
		// Close the channel when done generating
		close(numberChannel)
		fmt.Println("Number generation completed")
	}()

	// Wait for all workers to finish
	wg.Wait()
	fmt.Println("All workers have completed")
}
