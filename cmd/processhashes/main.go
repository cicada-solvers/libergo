package main

import (
	"config"
	"fmt"
	"github.com/jackc/pgx/v5"
	"liberdatabase"
	"math/big"
	"sync"
)

// main is the entry point of the program
func main() {
	configuration, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	db, err := liberdatabase.InitConnection()
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		return
	}
	defer func(db *pgx.Conn) {
		err := liberdatabase.CloseConnection(db)
		if err != nil {
			fmt.Printf("Error closing database connection: %v\n", err)
		}
	}(db)

	// Run removeProcessedRows at the beginning
	if err := liberdatabase.RemoveProcessedRows(db); err != nil {
		fmt.Printf("Error removing processed rows: %v\n", err)
		return
	}

	for {
		r, err := liberdatabase.GetByteArrayRange(db)
		if err != nil {
			fmt.Printf("Error getting byte array range: %v\n", err)
			return
		}
		if r == nil {
			break // No more rows to process
		}

		totalPermutations := big.NewInt(r.NumberOfPermutations)
		startArray, stopArray := r.StartArray, r.EndArray
		fmt.Printf("Processing: %v - %v\n", startArray, stopArray)

		program := NewProgram()

		var wg sync.WaitGroup
		numWorkers := configuration.NumWorkers
		wg.Add(numWorkers)

		done := make(chan struct{})
		var once sync.Once

		var mu sync.Mutex

		for j := 0; j < numWorkers; j++ {
			go processTasks(program.tasks, &wg, configuration.ExistingHash, done, &once, totalPermutations, &mu)
		}

		program.generateAllByteArrays(r.ArrayLength, startArray, stopArray)

		wg.Wait()

		select {
		case <-done:
		default:
		}

		if err := liberdatabase.RemoveItem(db, r.ID); err != nil {
			fmt.Printf("Error marking row as processed: %v\n", err)
		}
	}

	// Run removeProcessedRows at the end
	if err := liberdatabase.RemoveProcessedRows(db); err != nil {
		fmt.Printf("Error removing processed rows: %v\n", err)
		return
	}
}
