package main

import (
	"config"
	"fmt"
	"liberdatabase"
	"sync"
	"titler"
)

func main() {
	titler.PrintTitle("Single Hash Processor")

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

		rowCount := len(ranges)

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

		// Send tasks individually to the channel
		for _, perm := range ranges {
			program.tasks <- perm
		}
		close(program.tasks)

		fmt.Println("Waiting for all workers to finish...")
		wg.Wait()
	}

	fmt.Println("Program finished")
}
