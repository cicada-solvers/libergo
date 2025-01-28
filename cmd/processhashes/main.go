package main

import (
	"fmt"
	"math/big"
	"sync"
	"time"
)

// main is the entry point of the program
func main() {
	config, err := loadConfig("appsettings.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	db, err := initDatabase()
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		return
	}
	defer db.Close()

	// Run removeProcessedRows at the beginning
	if err := removeProcessedRows(db); err != nil {
		fmt.Printf("Error removing processed rows: %v\n", err)
		return
	}

	unprocessedCount, err := countUnprocessedRows(db)
	if err != nil {
		fmt.Printf("Error counting unprocessed rows: %v\n", err)
		return
	}
	fmt.Printf("Number of unprocessed rows: %d\n", unprocessedCount)

	ranges, err := getByteArrayRanges(db)
	if err != nil {
		fmt.Printf("Error getting byte array ranges: %v\n", err)
		return
	}

	for _, r := range ranges {
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

		startTime := time.Now()

		program.generateAllByteArrays(r.ArrayLength, startArray, stopArray)

		wg.Wait()

		select {
		case <-done:
		default:
		}

		duration := time.Since(startTime)
		fmt.Printf("Time taken to process range: %v\n", duration)

		if err := markAsProcessed(db, r.ID); err != nil {
			fmt.Printf("Error marking row as processed: %v\n", err)
		}

		unprocessedCount, err := countUnprocessedRows(db)
		if err != nil {
			fmt.Printf("Error counting unprocessed rows: %v\n", err)
			return
		}
		fmt.Printf("Number of unprocessed rows: %d\n", unprocessedCount)
	}

	// Run removeProcessedRows at the end
	if err := removeProcessedRows(db); err != nil {
		fmt.Printf("Error removing processed rows: %v\n", err)
		return
	}
}
