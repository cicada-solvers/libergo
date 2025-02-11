package main

import (
	"fmt"
	"liberdatabase"
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

func calculatePermutationRanges(length int) {
	totalPermutations := big.NewInt(1)
	for i := 0; i < length; i++ {
		totalPermutations.Mul(totalPermutations, big.NewInt(int64(len(primes))))
	}

	fmt.Printf("Total permutations are: %s\n", totalPermutations.String())

	var wg sync.WaitGroup
	fileChan := make(chan int64, 16384) // Increased buffer size

	go func() {
		for i := int64(0); i < totalPermutations.Int64(); i++ {
			fileChan <- i
		}
		close(fileChan)
	}()

	numWorkers := runtime.NumCPU() * 2 // Adjusted number of workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(fileChan, &wg, length, totalPermutations)
	}

	wg.Wait()
}

func worker(fileChan chan int64, wg *sync.WaitGroup, length int, totalPermutations *big.Int) {
	defer wg.Done()

	db, err := liberdatabase.InitConnection()
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		return
	}

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	nextPrintThreshold := big.NewInt(random.Int63n(100000-1000) + 1000)

	var batch []liberdatabase.Permutation
	batchSize := 2000 // Batch size for inserts

	for i := range fileChan {
		start := big.NewInt(i)
		startArray := indexToArray(start, length)

		perm := liberdatabase.Permutation{
			ID:                   uuid.New().String(),
			StartArray:           arrayToString(startArray),
			EndArray:             arrayToString(startArray),
			PackageName:          fmt.Sprintf("%d", 1),
			PermName:             fmt.Sprintf("%d", 1),
			ReportedToAPI:        false,
			Processed:            false,
			ArrayLength:          length,
			NumberOfPermutations: 1,
		}

		batch = append(batch, perm)

		if len(batch) >= batchSize {
			liberdatabase.InsertBatch(db, batch)
			batch = batch[:0] // Clear the batch
		}

		insertCountMutex.Lock()
		insertCounter.Add(insertCounter, big.NewInt(1))
		if insertCounter.Cmp(nextPrintThreshold) >= 0 {
			fmt.Printf("%s permutations of %s written to the database.\n", insertCounter.String(), totalPermutations.String())
			nextPrintThreshold = nextPrintThreshold.Add(nextPrintThreshold, big.NewInt(random.Int63n(1.5e9-1e8)+1e8))
		}
		insertCountMutex.Unlock()

		if start.Cmp(totalPermutations) == 0 {
			break
		}
	}

	// Insert any remaining records in the batch
	if len(batch) > 0 {
		liberdatabase.InsertBatch(db, batch)
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
