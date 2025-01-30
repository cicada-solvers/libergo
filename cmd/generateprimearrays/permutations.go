package main

import (
	"database/sql"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/google/uuid"
)

var primes = []int{
	2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109,
}

const (
	tableName = "permutations"
	batchSize = 100 // Number of inserts per batch
)

func calculatePermutationRanges(length int) {
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

	totalPermutations := big.NewInt(1)
	for i := 0; i < length; i++ {
		totalPermutations.Mul(totalPermutations, big.NewInt(int64(len(primes))))
	}

	var wg sync.WaitGroup
	fileChan := make(chan int64, config.NumWorkers*batchSize)

	go func() {
		for i := int64(0); i < totalPermutations.Int64(); i++ {
			fileChan <- i
		}
		close(fileChan)
	}()

	for i := 0; i < config.NumWorkers; i++ {
		wg.Add(1)
		go worker(db, fileChan, &wg, length, totalPermutations)
	}

	wg.Wait()

	// Compact the database to reclaim unused space
	fmt.Println("Compacting database...")
	_ = compactDatabase(db)
}

func worker(db *sql.DB, fileChan chan int64, wg *sync.WaitGroup, length int, totalPermutations *big.Int) {
	defer wg.Done()
	batch := make([]interface{}, 0, batchSize*9) // 9 columns per insert

	for i := range fileChan {
		start := big.NewInt(i)
		startArray := indexToArray(start, length)

		id := uuid.New().String()
		packageFileName := fmt.Sprintf("%d", 1)
		permFileName := fmt.Sprintf("%d", 1)
		reportedToAPI := false
		processed := false

		batch = append(batch, id, arrayToString(startArray), arrayToString(startArray), packageFileName, permFileName, reportedToAPI, processed, length, 1)

		if len(batch) >= batchSize*9 {
			err := insertBatch(db, batch)
			if err != nil {
				fmt.Printf("Error inserting batch into database: %v\n", err)
			}
			batch = batch[:0]
		}

		if start.Cmp(totalPermutations) == 0 {
			break
		}
	}

	if len(batch) > 0 {
		err := insertBatch(db, batch)
		if err != nil {
			fmt.Printf("Error inserting final batch into database: %v\n", err)
		}
	}
}

func insertBatch(db *sql.DB, batch []interface{}) error {
	query := fmt.Sprintf("INSERT INTO %s (id, startArray, endArray, packageName, permName, reportedToAPI, processed, arrayLength, numberOfPermutations) VALUES %s",
		tableName, strings.Repeat("(?, ?, ?, ?, ?, ?, ?, ?, ?),", len(batch)/9))
	query = query[:len(query)-1] // Remove the trailing comma

	return insertWithRetry(db, query, batch...)
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
