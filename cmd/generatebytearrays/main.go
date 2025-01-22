package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type Config struct {
	NumWorkers int `json:"num_workers"`
}

func loadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("error decoding config file: %v", err)
	}

	return &config, nil
}

func calculatePermutationRanges(length int, maxPermutationsPerFile int64) {
	config, err := loadConfig("appsettings.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	totalPermutations := big.NewInt(1)
	for i := 0; i < length; i++ {
		totalPermutations.Mul(totalPermutations, big.NewInt(256))
	}

	numFiles := new(big.Int).Div(totalPermutations, big.NewInt(maxPermutationsPerFile))
	if new(big.Int).Mod(totalPermutations, big.NewInt(maxPermutationsPerFile)).Cmp(big.NewInt(0)) != 0 {
		numFiles.Add(numFiles, big.NewInt(1))
	}

	// Create the subdirectory
	folder := fmt.Sprintf("%010d", length)
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		fmt.Printf("Error creating folder: %v\n", err)
		return
	}

	var wg sync.WaitGroup
	fileChan := make(chan int64, numFiles.Int64())

	for i := int64(0); i < numFiles.Int64(); i++ {
		fileChan <- i
	}
	close(fileChan)

	for i := 0; i < config.NumWorkers; i++ { // Use the number of workers from the config
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range fileChan {
				start := new(big.Int).Mul(big.NewInt(i), big.NewInt(maxPermutationsPerFile))
				end := new(big.Int).Add(start, big.NewInt(maxPermutationsPerFile))
				if end.Cmp(totalPermutations) > 0 {
					end = totalPermutations
				}

				fileName := fmt.Sprintf("permutations_%d.txt", i)
				filePath := filepath.Join(folder, fileName)
				file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Printf("Error opening file: %v\n", err)
					return
				}

				startArray := indexToArray(start, length)
				endArray := indexToArray(new(big.Int).Sub(end, big.NewInt(1)), length)

				_, err = file.WriteString(fmt.Sprintf("%s\n%s\n", arrayToString(startArray), arrayToString(endArray)))
				if err != nil {
					fmt.Printf("Error writing to file: %v\n", err)
				}

				err = file.Close()
				if err != nil {
					fmt.Printf("Error closing file: %v\n", err)
				}
			}
		}()
	}

	wg.Wait()
}

func indexToArray(index *big.Int, length int) []byte {
	array := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		mod := new(big.Int)
		index.DivMod(index, big.NewInt(256), mod)
		array[i] = byte(mod.Int64())
	}
	return array
}

func arrayToString(array []byte) string {
	strArray := make([]string, len(array))
	for i, b := range array {
		strArray[i] = fmt.Sprintf("%d", b)
	}
	return strings.Join(strArray, ",")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./generatebytearrays <length>")
		return
	}

	length, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Invalid length: %v\n", err)
		return
	}

	maxPermutationsPerFile := int64(5000000000)
	calculatePermutationRanges(length, maxPermutationsPerFile)
}
