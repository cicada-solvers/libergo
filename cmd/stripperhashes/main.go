package main

import (
	"config"
	"fmt"
	"liberdatabase"
	"sync"
)

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

	// Run removeProcessedRows at the beginning
	liberdatabase.RemoveProcessedRows(db)

	for {
		ranges, err := liberdatabase.GetByteArrayRanges(db)
		if err != nil {
			fmt.Printf("Error getting byte array ranges: %v\n", err)
			return
		}

		if ranges == nil || len(ranges) == 0 {
			fmt.Println("No more rows to process")
			break
		}

		rowCount := liberdatabase.GetCountOfPermutations(db)
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

			liberdatabase.RemoveItem(db, r.ID)
		}

		close(program.tasks)
		wg.Wait()

		// Run removeProcessedRows at the end of each batch
		liberdatabase.RemoveProcessedRows(db)
	}
}
