package main

import (
	"fmt"
	"math/big"
	"os"
	"sync"
)

// main is the entry point of the program
func main() {
	existingHash := "36367763ab73783c7af284446c59466b4cd653239a311cb7116d4618dee09a8425893dc7500b464fdaf1672d7bef5e891c6e2274568926a49fb4f45132c2a8b4"
	fmt.Printf("Existing hash: %s\n", existingHash)

	var packNumber int
	var numWorkersEntry int

	// Prompt for pack number
	fmt.Print("Enter the pack number: ")
	_, err := fmt.Scan(&packNumber)
	if err != nil {
		fmt.Printf("Error reading pack number: %v\n", err)
		return
	}

	// Prompt for number of processors
	fmt.Print("Enter the number of processors you want to use: ")
	_, err = fmt.Scan(&numWorkersEntry)
	if err != nil {
		fmt.Printf("Error reading number of processors: %v\n", err)
		return
	}

	// Check to see if output directory exists
	var outputDir string
	outputDirCheck := fmt.Sprintf("PACK_%d", packNumber)
	if _, err := os.Stat(outputDirCheck); os.IsNotExist(err) {
		// Download the pack
		outputDirDownload, downloadErr := downloadAndExtractPack(packNumber)
		if downloadErr != nil {
			fmt.Printf("Error downloading pack: %v\n", downloadErr)
			return
		}

		outputDir = outputDirDownload
	} else {
		outputDir = outputDirCheck
	}

	// List all CSV files in the output directory
	csvFiles, err := listCSVFiles(outputDir)
	if err != nil {
		fmt.Printf("Error listing CSV files: %v\n", err)
		return
	}

	fmt.Printf("Found %d CSV files in %s\n", len(csvFiles), outputDir)
	// Process each CSV file
	for _, csvFile := range csvFiles {
		fmt.Printf("Processing file: %s\n", csvFile)

		lines, readErr := readCSVToPackFiles(csvFile)
		if readErr != nil {
			fmt.Printf("Error reading CSV file: %v\n", readErr)
			return
		}
		for _, line := range lines {
			totalPermutations := big.NewInt(line.NumberOfPermutations)
			startArray, stopArray := line.StartArray, line.EndArray
			fmt.Printf("Processing: %v - %v\n", startArray, stopArray)

			program := NewArrayGen()

			var wg sync.WaitGroup
			numWorkers := numWorkersEntry
			wg.Add(numWorkers)

			done := make(chan struct{})
			var once sync.Once

			var mu sync.Mutex

			for j := 0; j < numWorkers; j++ {
				go processTasks(program.segments, &wg, existingHash, done, &once, totalPermutations, &mu)
			}

			program.generateAllByteArrays(line.ArrayLength, startArray, stopArray)

			wg.Wait()

			select {
			case <-done:
			default:
			}

			deleteError := deleteLineById(csvFile, line.Id)
			if deleteError != nil {
				fmt.Printf("Error deleting line with ID %s: %v\n", line.Id, deleteError)
			} else {
				fmt.Printf("Deleted line with ID: %s\n", line.Id)
			}
		}

		fmt.Printf("Finished processing file: %s\n", csvFile)
		// Delete the CSV file after processing
		deleteErr := os.Remove(csvFile)
		if deleteErr != nil {
			fmt.Printf("Error deleting file %s: %v\n", csvFile, deleteErr)
		} else {
			fmt.Printf("Deleted file: %s\n", csvFile)
		}
	}
}
