package main

import (
	"config"
	"fmt"
	"liberdatabase"
	"math/big"
	"sync"
	"titler"
)

// main is the entry point of the program
func main() {
	titler.PrintTitle("Process Hashes")

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

	// Run removeProcessedRows at the beginning
	liberdatabase.RemoveProcessedRows(db)

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

		liberdatabase.RemoveItem(db, r.ID)
	}

	// Run removeProcessedRows at the end
	liberdatabase.RemoveProcessedRows(db)
}
