package main

import (
	"database/sql"
	"fmt"
	"math/big"
	"runtime"
	"strings"
	"sync"

	"github.com/google/uuid"
)

var primes = []int{
	2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109,
}

const (
	tableName = "permutations"
	batchSize = 50 // Number of inserts per batch
)

func calculatePermutationRanges(length int) {
	db, err := initDatabase()
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		return
	}
	defer db.Close()

	totalPermutations := big.NewInt(1)
	for i := 0; i < length; i++ {
		totalPermutations.Mul(totalPermutations, big.NewInt(int64(len(primes))))
	}

	var wg sync.WaitGroup
	fileChan := make(chan int64, runtime.NumCPU()*batchSize)

	go func() {
		for i := int64(0); i < totalPermutations.Int64(); i++ {
			fileChan <- i
		}
		close(fileChan)
	}()

	numWorkers := runtime.NumCPU() // Get the number of CPU cores
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(db, fileChan, &wg, length, totalPermutations)
	}

	wg.Wait()

	// Compact the database to reclaim unused space
	fmt.Println("Compacting database...")
	_ = compactDatabase(db)
}

var dbMutex sync.Mutex

func worker(db *sql.DB, fileChan chan int64, wg *sync.WaitGroup, length int, totalPermutations *big.Int) {
	defer wg.Done()

	insertCount := 0
	maxRetries := 10

	for i := range fileChan {
		start := big.NewInt(i)
		startArray := indexToArray(start, length)

		id := uuid.New().String()
		packageFileName := fmt.Sprintf("%d", 1)
		permFileName := fmt.Sprintf("%d", 1)
		reportedToAPI := false
		processed := false

		retryCount := 0
		for {
			dbMutex.Lock()
			tx, err := db.Begin()
			if err != nil {
				dbMutex.Unlock()
				fmt.Printf("Error starting transaction: %v\n", err)
				return
			}

			_, err = tx.Exec(fmt.Sprintf("INSERT INTO %s (id, startArray, endArray, packageName, permName, reportedToAPI, processed, arrayLength, numberOfPermutations) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1)",
				tableName), id, arrayToString(startArray), arrayToString(startArray), packageFileName, permFileName, reportedToAPI, processed, length)
			if err != nil {
				if strings.Contains(err.Error(), "database is locked") {
					retryCount++
					fmt.Println("Database is locked, retrying insert...")
					if retryCount >= maxRetries {
						fmt.Printf("Max retries reached for record %s, skipping insert\n", id)
						dbMutex.Unlock()
						break
					}
					dbMutex.Unlock()
					continue
				} else {
					fmt.Printf("Error inserting record into database: %v\n", err)
					dbMutex.Unlock()
					return
				}
			}

			err = tx.Commit()
			dbMutex.Unlock()
			if err != nil {
				fmt.Printf("Error committing transaction: %v\n", err)
				return
			}
			break
		}

		if retryCount > 0 && retryCount < maxRetries {
			fmt.Printf("Insert successful after %d retries\n", retryCount)
		}

		insertCount++
		if insertCount >= 5000 {
			dbMutex.Lock()
			tx, err := db.Begin()
			if err != nil {
				dbMutex.Unlock()
				fmt.Printf("Error starting new transaction: %v\n", err)
				return
			}
			err = tx.Commit()
			if err != nil {
				fmt.Printf("Error committing transaction: %v\n", err)
				dbMutex.Unlock()
				return
			}
			dbMutex.Unlock()
			insertCount = 0
		}

		if start.Cmp(totalPermutations) == 0 {
			break
		}
	}

	if insertCount > 0 {
		dbMutex.Lock()
		tx, err := db.Begin()
		if err != nil {
			dbMutex.Unlock()
			fmt.Printf("Error starting final transaction: %v\n", err)
			return
		}
		err = tx.Commit()
		if err != nil {
			fmt.Printf("Error committing final transaction: %v\n", err)
		}
		dbMutex.Unlock()
	}
}

func indexToArray(index *big.Int, length int) []int {
	array := make([]int, length)
	primeLen := big.NewInt(int64(len(primes)))
	for i := length - 1; i >= 0; i-- {
		mod := new(big.Int)
		index.DivMod(index, primeLen, mod)
		array[i] = primes[mod.Int64()]
	}
	return array
}

func arrayToString(array []int) string {
	strArray := make([]string, len(array))
	for i, b := range array {
		strArray[i] = fmt.Sprintf("%d", b)
	}
	return strings.Join(strArray, ",")
}
