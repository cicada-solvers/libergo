package main

import (
	"fmt"
	"gorm.io/gorm"
	"liberdatabase"
	"math"
	"math/big"
	"runtime"
	"sequences"
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

	// Start the workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			// Each worker processes numbers from the channel
			for num := range numberChannel {
				// Convert int64 to big.Int for compatibility with IsPrime
				bigNum := big.NewInt(num)
				if sequences.IsPrime(bigNum) {
					liberdatabase.AddPrimeValue(conn, num)
				}
			}
			fmt.Printf("Worker %d completed\n", workerID)
		}(i)
	}

	// Start the number generator in a separate goroutine
	go func() {
		// Generate numbers from 2 to MaxInt32
		for i := int64(2); i <= int64(math.MaxInt32); i++ {
			numberChannel <- i
		}
		// Close the channel when done generating
		close(numberChannel)
		fmt.Println("Number generation completed")
	}()

	// Wait for all workers to finish
	wg.Wait()
	fmt.Println("All workers have completed")
}
