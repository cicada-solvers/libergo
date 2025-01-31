package main

import (
	"database/sql"
	"fmt"
	"math/big"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

var insertCountMutex sync.Mutex
var insertCounter = big.NewInt(0)
var primes = []int{
	2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109,
}

const batchSize = 500 // Number of inserts per batch

func calculatePermutationRanges(length int) {
	totalPermutations := big.NewInt(1)
	for i := 0; i < length; i++ {
		totalPermutations.Mul(totalPermutations, big.NewInt(int64(len(primes))))
	}

	fmt.Printf("Total permutations are: %s\n", totalPermutations.String())

	db, err := initDatabase()
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		return
	}
	defer db.Close()

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
		go worker(fileChan, &wg, length, totalPermutations, i+1, db)
	}

	wg.Wait()

	// Compact the database to reclaim unused space
	fmt.Println("Compacting database...")
	_ = compactDatabase(db)
}

func worker(fileChan chan int64, wg *sync.WaitGroup, length int, totalPermutations *big.Int, workerIndex int, db *sql.DB) {
	defer wg.Done()

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	nextPrintThreshold := big.NewInt(random.Int63n(100000-1000) + 1000)

	var builder strings.Builder

	// Start the transaction
	//builder.WriteString("BEGIN TRANSACTION;\n")

	// Start the insert statement
	builder.WriteString("INSERT INTO permutations (id, startArray, endArray, packageName, permName, reportedToAPI, processed, arrayLength, numberOfPermutations) VALUES \n")

	firstTime := true
	insertCount := 0

	for i := range fileChan {
		start := big.NewInt(i)
		startArray := indexToArray(start, length)

		id := uuid.New().String()
		packageFileName := fmt.Sprintf("%d", 1)
		permFileName := fmt.Sprintf("%d", 1)
		reportedToAPI := false
		processed := false

		sqlStatement := ""

		if !firstTime {
			sqlStatement = fmt.Sprintf("\n,('%s', '%s', '%s', '%s', '%s', %t, %t, %d, 1)",
				id, arrayToString(startArray), arrayToString(startArray), packageFileName, permFileName, reportedToAPI, processed, length)
		} else {
			sqlStatement = fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', %t, %t, %d, 1)",
				id, arrayToString(startArray), arrayToString(startArray), packageFileName, permFileName, reportedToAPI, processed, length)
			firstTime = false
		}

		builder.WriteString(sqlStatement)

		if insertCount >= batchSize {
			// Commit the current transaction before closing the file
			//builder.WriteString(";\nCOMMIT;\n")

			// Need to write to the database here...
			err := insertWithRetry(db, builder.String())
			if err != nil {
				fmt.Printf("Error inserting into database: %v - %s\n", err, builder.String())
			}

			// Clear the builder
			builder.Reset()

			// Start a new transaction
			//builder.WriteString("BEGIN TRANSACTION;\n")

			// Start the insert statement
			builder.WriteString("INSERT INTO permutations (id, startArray, endArray, packageName, permName, reportedToAPI, processed, arrayLength, numberOfPermutations) VALUES \n")

			firstTime = true
			insertCount = 0
		}

		insertCount++
		insertCountMutex.Lock()
		insertCounter.Add(insertCounter, big.NewInt(1))
		if insertCounter.Cmp(nextPrintThreshold) >= 0 {
			fmt.Printf("%s permutations of %s written to the file.\n", insertCounter.String(), totalPermutations.String())
			nextPrintThreshold = nextPrintThreshold.Add(nextPrintThreshold, big.NewInt(random.Int63n(1.5e9-1e8)+1e8))
		}
		insertCountMutex.Unlock()

		if start.Cmp(totalPermutations) == 0 {
			break
		}
	}

	// Commit the transaction at the end of the file
	//builder.WriteString(";\nCOMMIT;\n")

	// Need to write to the database here...
	err := insertWithRetry(db, builder.String())
	if err != nil {
		fmt.Printf("Error inserting into database: %v - %s\n", err, builder.String())
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
