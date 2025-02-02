package main

import (
	"config"
	"fmt"
	"liberdatabase"
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
	defer func() {
		if err := liberdatabase.CloseConnection(db); err != nil {
			fmt.Printf("Error closing database connection: %v\n", err)
		}
	}()

	// Run removeProcessedRows at the beginning
	if err := liberdatabase.RemoveProcessedRows(db); err != nil {
		fmt.Printf("Error removing processed rows: %v\n", err)
		return
	}

	ranges, err := liberdatabase.GetByteArrayRanges(db)
	if err != nil {
		fmt.Printf("Error getting byte array ranges: %v\n", err)
		return
	}

	if len(ranges) == 0 {
		fmt.Println("No more rows to process")
		return
	}

	rowCount, _ := liberdatabase.GetCountOfPermutations()
	fmt.Printf("Total number of permutations: %d\n", rowCount)

	program := NewProgram()

	var wg sync.WaitGroup
	numWorkers := configuration.NumWorkers
	wg.Add(numWorkers)

	done := make(chan struct{})
	var once sync.Once

	for j := 0; j < numWorkers; j++ {
		go processTasks(program.tasks, &wg, configuration.ExistingHash, done, &once)
	}

	for _, r := range ranges {
		startArray := r.StartArray

		// Since startArray and stopArray are the same, we can send it directly to tasks
		program.tasks <- startArray

		select {
		case <-done:
		default:
		}

		if err := liberdatabase.RemoveItem(db, r.ID); err != nil {
			fmt.Printf("Error marking row as processed: %v\n", err)
		}
	}

	close(program.tasks)
	wg.Wait()

	// Run removeProcessedRows at the end
	if err := liberdatabase.RemoveProcessedRows(db); err != nil {
		fmt.Printf("Error removing processed rows: %v\n", err)
		return
	}
}
