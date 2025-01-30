package main

import (
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

var primes = []int{
	2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109,
}

const (
	tableName   = "permutations"
	batchSize   = 500                    // Number of inserts per batch
	maxFileSize = 5 * 1024 * 1024 * 1024 // 5 GB
)

func createTableScript() {
	script := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	id TEXT PRIMARY KEY,
	startArray TEXT,
	endArray TEXT,
	packageName TEXT,
	permName TEXT,
	reportedToAPI BOOLEAN,
	processed BOOLEAN,
	arrayLength INTEGER,
	numberOfPermutations INTEGER
);`, tableName)

	file, err := os.Create("create_table.sql")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(script)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	fmt.Println("Database table creation script written to create_table.sql")
}

func calculatePermutationRanges(length int) {
	totalPermutations := big.NewInt(1)
	for i := 0; i < length; i++ {
		totalPermutations.Mul(totalPermutations, big.NewInt(int64(len(primes))))
	}

	fmt.Printf("Total permutations are: %s\n", totalPermutations.String())

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
		go worker(fileChan, &wg, length, totalPermutations, i+1)
	}

	wg.Wait()
}

var insertCountMutex sync.Mutex
var insertCount = big.NewInt(0)

func getFileSize(file *os.File) int64 {
	info, err := file.Stat()
	if err != nil {
		fmt.Printf("Error getting file size: %v\n", err)
		return 0
	}
	return info.Size()
}

func worker(fileChan chan int64, wg *sync.WaitGroup, length int, totalPermutations *big.Int, workerIndex int) {
	defer wg.Done()

	maxRetries := 10
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	nextPrintThreshold := big.NewInt(random.Int63n(100000-1000) + 1000)

	fileName := fmt.Sprintf("sql_statements_%d_%d.sql", length, workerIndex)
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Start the transaction
	_, err = file.WriteString("BEGIN TRANSACTION;\n")
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	for i := range fileChan {
		start := big.NewInt(i)
		startArray := indexToArray(start, length)

		id := uuid.New().String()
		packageFileName := fmt.Sprintf("%d", 1)
		permFileName := fmt.Sprintf("%d", 1)
		reportedToAPI := false
		processed := false

		sqlStatement := fmt.Sprintf("INSERT INTO %s (id, startArray, endArray, packageName, permName, reportedToAPI, processed, arrayLength, numberOfPermutations) VALUES ('%s', '%s', '%s', '%s', '%s', %t, %t, %d, 1);\n",
			tableName, id, arrayToString(startArray), arrayToString(startArray), packageFileName, permFileName, reportedToAPI, processed, length)

		if getFileSize(file) >= maxFileSize {
			// Commit the current transaction before closing the file
			_, err = file.WriteString("COMMIT;\n")
			if err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
				return
			}
			file.Close()
			fileName = fmt.Sprintf("sql_statements_%d_%d.sql", length, workerIndex)
			file, err = os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Printf("Error opening new file: %v\n", err)
				return
			}
			// Start a new transaction
			_, err = file.WriteString("BEGIN TRANSACTION;\n")
			if err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
				return
			}
		}

		_, err = file.WriteString(sqlStatement)
		if err != nil {
			if strings.Contains(err.Error(), "file is locked") {
				retryCount := 0
				for retryCount < maxRetries {
					fmt.Println("File is locked, retrying write...")
					time.Sleep(time.Millisecond * 100)
					_, err = file.WriteString(sqlStatement)
					if err == nil {
						break
					}
					retryCount++
				}
				if retryCount >= maxRetries {
					fmt.Printf("Max retries reached for record %s, skipping write\n", id)
					continue
				}
			} else {
				fmt.Printf("Error writing record to file: %v\n", err)
				return
			}
		}

		insertCountMutex.Lock()
		insertCount.Add(insertCount, big.NewInt(1))
		if insertCount.Cmp(nextPrintThreshold) >= 0 {
			fmt.Printf("%s permutations of %s written to the file.\n", insertCount.String(), totalPermutations.String())
			nextPrintThreshold = nextPrintThreshold.Add(nextPrintThreshold, big.NewInt(random.Int63n(1.5e9-1e8)+1e8))
		}
		insertCountMutex.Unlock()

		if start.Cmp(totalPermutations) == 0 {
			break
		}
	}

	// Commit the transaction at the end of the file
	_, err = file.WriteString("COMMIT;\n")
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
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
