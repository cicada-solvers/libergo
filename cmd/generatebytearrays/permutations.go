package main

import (
	"config"
	"encoding/csv"
	"fmt"
	"liberdatabase"
	"math/big"
	"os"
	"strings"

	"github.com/google/uuid"
)

// calculatePermutationRanges calculates the permutation ranges for the specified length
func calculatePermutationRanges(length int, maxPermutationsPerLine, maxPermutationsPerSegment int64, packageFileNumber *big.Int, config *config.AppConfig) {
	permutationChan := make(chan liberdatabase.Permutation, 2048) // Increased buffer size

	// Loop through the permutationChan in a single thread
	go func() {
		for perm := range permutationChan {
			filePath := fmt.Sprintf("permutations_%d_%d.csv", length, packageFileNumber.Int64())
			if err := writePermutationToCSV(filePath, perm); err != nil {
				fmt.Printf("Error writing permutation to CSV: %v\n", err)
			}
		}
	}()

	totalPermutations := big.NewInt(1)
	for i := 0; i < length; i++ {
		totalPermutations.Mul(totalPermutations, big.NewInt(256))
	}

	totalPackageFiles, err := calculateNumberOfPackageFiles(length, maxPermutationsPerLine, maxPermutationsPerSegment, config.MaxSegmentsPerPackage)
	if err != nil {
		fmt.Printf("Error calculating number of packages: %v\n", err)
		return
	}

	fmt.Printf("Current package: %s of %s\n", packageFileNumber.String(), totalPackageFiles.String())

	startFile := new(big.Int).Mul(new(big.Int).Sub(packageFileNumber, big.NewInt(1)), big.NewInt(config.MaxSegmentsPerPackage))
	endFile := new(big.Int).Add(startFile, big.NewInt(config.MaxSegmentsPerPackage))
	for i := new(big.Int).Set(startFile); i.Cmp(endFile) < 0; i.Add(i, big.NewInt(1)) {
		start := new(big.Int).Mul(i, big.NewInt(maxPermutationsPerLine*maxPermutationsPerSegment))
		if start.Cmp(totalPermutations) >= 0 {
			break
		}

		for j := int64(0); j < maxPermutationsPerSegment; j++ {
			lineStart := new(big.Int).Add(start, big.NewInt(j*maxPermutationsPerLine))
			lineEnd := new(big.Int).Add(lineStart, big.NewInt(maxPermutationsPerLine))
			if lineEnd.Cmp(totalPermutations) > 0 {
				lineEnd = totalPermutations
			}

			startArray := indexToArray(lineStart, length)
			endArray := indexToArray(new(big.Int).Sub(lineEnd, big.NewInt(1)), length)

			perm := liberdatabase.Permutation{
				ID:                   uuid.New().String(),
				StartArray:           arrayToString(startArray),
				EndArray:             arrayToString(endArray),
				PackageName:          fmt.Sprintf("package_%d", packageFileNumber),
				PermName:             fmt.Sprintf("permutation_seg_%d", i.Int64()),
				ReportedToAPI:        false,
				Processed:            false,
				ArrayLength:          length,
				NumberOfPermutations: config.MaxPermutationsPerLine,
			}

			permutationChan <- perm

			if lineEnd.Cmp(totalPermutations) == 0 {
				break
			}
		}
	}
	close(permutationChan) // Close the channel after all permutations are sent
}

// writePermutationToCSV writes a liberdatabase.Permutation to a CSV file
func writePermutationToCSV(filePath string, perm liberdatabase.Permutation) error {
	// Open the file for appending or create it if it doesn't exist
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header if the file is empty
	fileInfo, _ := file.Stat()
	if fileInfo.Size() == 0 {
		header := []string{
			"ID", "StartArray", "EndArray", "PackageName", "PermName",
			"ReportedToAPI", "Processed", "ArrayLength", "NumberOfPermutations",
		}
		if err := writer.Write(header); err != nil {
			return fmt.Errorf("error writing header: %v", err)
		}
	}

	// Write the permutation data
	record := []string{
		perm.ID,
		perm.StartArray,
		perm.EndArray,
		perm.PackageName,
		perm.PermName,
		fmt.Sprintf("%t", perm.ReportedToAPI),
		fmt.Sprintf("%t", perm.Processed),
		fmt.Sprintf("%d", perm.ArrayLength),
		fmt.Sprintf("%d", perm.NumberOfPermutations),
	}
	if err := writer.Write(record); err != nil {
		return fmt.Errorf("error writing record: %v", err)
	}

	return nil
}

// indexToArray converts an index to a byte array
func indexToArray(index *big.Int, length int) []byte {
	array := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		mod := new(big.Int)
		index.DivMod(index, big.NewInt(256), mod)
		array[i] = byte(mod.Int64())
	}
	return array
}

// arrayToString converts a byte array to a string
func arrayToString(array []byte) string {
	strArray := make([]string, len(array))
	for i, b := range array {
		strArray[i] = fmt.Sprintf("%d", b)
	}
	return strings.Join(strArray, ",")
}

// calculateNumberOfPackageFiles calculates the number of package files required to store all permutations
func calculateNumberOfPackageFiles(length int, maxPermutationsPerLine, maxPermutationsPerFile, maxFilesPerPackage int64) (*big.Int, error) {
	totalPermutations := big.NewInt(1)
	for i := 0; i < length; i++ {
		totalPermutations.Mul(totalPermutations, big.NewInt(256))
	}

	permutationsPerFile := big.NewInt(maxPermutationsPerLine * maxPermutationsPerFile)
	totalFiles := new(big.Int).Div(totalPermutations, permutationsPerFile)
	if new(big.Int).Mod(totalPermutations, permutationsPerFile).Cmp(big.NewInt(0)) != 0 {
		totalFiles.Add(totalFiles, big.NewInt(1))
	}

	totalPackageFiles := new(big.Int).Div(totalFiles, big.NewInt(maxFilesPerPackage))
	if new(big.Int).Mod(totalFiles, big.NewInt(maxFilesPerPackage)).Cmp(big.NewInt(0)) != 0 {
		totalPackageFiles.Add(totalPackageFiles, big.NewInt(1))
	}

	return totalPackageFiles, nil
}
