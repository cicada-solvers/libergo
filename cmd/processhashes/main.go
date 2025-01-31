package main

import (
	"fmt"
	"math/big"
	"sync"
)

// main is the entry point of the program
func main() {
	config, err := loadConfig("appsettings.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	db, err := initConnection()
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		return
	}

	// Run removeProcessedRows at the beginning
	if err := removeProcessedRows(db); err != nil {
		fmt.Printf("Error removing processed rows: %v\n", err)
		return
	}

	for {
		r, err := getByteArrayRange(db)
		if err != nil {
			fmt.Printf("Error getting byte array range: %v\n", err)
			return
		}
		if r == nil {
			break // No more rows to process
		}

		totalPermutations := big.NewInt(int64(r.NumberOfPermutations))
		startArray, stopArray := r.StartArray, r.EndArray
		fmt.Printf("Processing: %v - %v\n", startArray, stopArray)

		program := NewProgram()

		var wg sync.WaitGroup
		numWorkers := config.NumWorkers
		wg.Add(numWorkers)

		done := make(chan struct{})
		var once sync.Once

		var mu sync.Mutex

		for j := 0; j < numWorkers; j++ {
			go processTasks(program.tasks, &wg, config.ExistingHash, done, &once, totalPermutations, &mu)
		}

		program.generateAllByteArrays(r.ArrayLength, startArray, stopArray)

		wg.Wait()

		select {
		case <-done:
		default:
		}

		if err := removeItem(db, r.ID); err != nil {
			fmt.Printf("Error marking row as processed: %v\n", err)
		}
	}

	// Run removeProcessedRows at the end
	if err := removeProcessedRows(db); err != nil {
		fmt.Printf("Error removing processed rows: %v\n", err)
		return
	}
}
