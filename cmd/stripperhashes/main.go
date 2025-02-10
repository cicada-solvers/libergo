package main

import (
	"config"
	"fmt"
	"liberdatabase"
	"sync"
)

func main() {
	fmt.Println("Starting the program...")

	configuration, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}
	fmt.Println("Configuration loaded successfully")

	db, err := liberdatabase.InitConnection()
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		return
	}
	fmt.Println("Database connection initialized")

	// Run removeProcessedRows at the beginning
	fmt.Println("Removing processed rows...")
	liberdatabase.RemoveProcessedRows(db)
	fmt.Println("Processed rows removed")

	for {
		fmt.Println("Fetching byte array ranges...")
		ranges, err := liberdatabase.GetByteArrayRanges(db)
		if err != nil {
			fmt.Printf("Error getting byte array ranges: %v\n", err)
			return
		}

		if ranges == nil || len(ranges) == 0 {
			fmt.Println("No more rows to process")
			break
		}
		fmt.Printf("Fetched %d byte array ranges\n", len(ranges))

		rowCount := liberdatabase.GetCountOfPermutations(db)
		fmt.Printf("Total number of permutations: %d\n", rowCount)

		program := NewProgram()

		var wg sync.WaitGroup
		numWorkers := configuration.NumWorkers
		wg.Add(numWorkers)

		done := make(chan struct{})
		var once sync.Once

		fmt.Printf("Starting %d workers...\n", numWorkers)
		for j := 0; j < numWorkers; j++ {
			go processTasks(program.tasks, &wg, configuration.ExistingHash, done, &once, &rowCount)
		}

		for _, r := range ranges {
			startArray := r.StartArray

			// Since startArray and stopArray are the same, we can send it directly to tasks
			program.tasks <- startArray

			select {
			case <-done:
				fmt.Println("Processing done signal received")
			default:
			}

			liberdatabase.RemoveItem(db, r.ID)
		}

		close(program.tasks)
		fmt.Println("Waiting for all workers to finish...")
		wg.Wait()

		// Run removeProcessedRows at the end of each batch
		fmt.Println("Removing processed rows...")
		liberdatabase.RemoveProcessedRows(db)
		fmt.Println("Processed rows removed")
	}

	fmt.Println("Program finished")
}
