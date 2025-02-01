package main

import (
	"fmt"
	"sync"
)

// main is the entry point of the program
func main() {
	// Fancy print statement
	fmt.Println("****************************************")
	fmt.Println("*                                      *")
	fmt.Println("*   This program only dances for       *")
	fmt.Println("*             SINGLES!                 *")
	fmt.Println("*         Make is rain!!               *")
	fmt.Println("*                                      *")
	fmt.Println("****************************************")

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

	ranges, err := getByteArrayRanges(db)
	if err != nil {
		fmt.Printf("Error getting byte array ranges: %v\n", err)
		return
	}

	if len(ranges) == 0 {
		fmt.Println("No more rows to process")
		return
	}

	program := NewProgram()

	var wg sync.WaitGroup
	numWorkers := config.NumWorkers
	wg.Add(numWorkers)

	done := make(chan struct{})
	var once sync.Once

	for j := 0; j < numWorkers; j++ {
		go processTasks(program.tasks, &wg, config.ExistingHash, done, &once)
	}

	for _, r := range ranges {
		startArray := r.StartArray

		// Since startArray and stopArray are the same, we can send it directly to tasks
		program.tasks <- startArray

		select {
		case <-done:
		default:
		}

		if err := removeItem(db, r.ID); err != nil {
			fmt.Printf("Error marking row as processed: %v\n", err)
		}
	}

	close(program.tasks)
	wg.Wait()

	// Run removeProcessedRows at the end
	if err := removeProcessedRows(db); err != nil {
		fmt.Printf("Error removing processed rows: %v\n", err)
		return
	}
}
