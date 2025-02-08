package main

import (
	"config"
	"fmt"
	"liberdatabase"
	"math/big"
	"strings"

	"github.com/google/uuid"
)

// calculatePermutationRanges calculates the permutation ranges for the specified length
func calculatePermutationRanges(length int, maxPermutationsPerLine, maxPermutationsPerSegment int64, packageFileNumber *big.Int, config *config.AppConfig) {
	db, err := liberdatabase.InitConnection()
	if err != nil {
		fmt.Printf("Error initializing database connection: %v\n", err)
		return
	}

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
		fmt.Printf("Processing file index: %d\n", i.Int64())

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

			liberdatabase.InsertRecord(db, perm)

			if lineEnd.Cmp(totalPermutations) == 0 {
				break
			}
		}
	}
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
