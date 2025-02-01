package main

import (
	"fmt"
	"math/big"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "init" {
		fmt.Println("Initializing database.")
		_, err := initDatabase()
		if err != nil {
			fmt.Printf("Error initializing database: %v\n", err)
			return
		}

		fmt.Println("Database initialized successfully.")
		return
	}

	var length int
	fmt.Print("Enter the array length: ")
	_, err := fmt.Scan(&length)
	if err != nil {
		fmt.Printf("Invalid input: %v\n", err)
		return
	}

	config, err := loadConfig("appsettings.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	maxPermutationsPerLine := config.MaxPermutationsPerLine
	maxPermutationsPerFile := config.MaxRangesPerSegment

	totalPackageFiles, err := calculateNumberOfPackageFiles(length, maxPermutationsPerLine, maxPermutationsPerFile, config.MaxSegmentsPerPackage)
	if err != nil {
		fmt.Printf("Error calculating number of packages: %v\n", err)
		return
	}

	fmt.Printf("Total number of packages: %s\n", totalPackageFiles.String())

	if totalPackageFiles.Cmp(big.NewInt(1)) == 0 {
		fmt.Println("Only one package to generate.")
		calculatePermutationRanges(length, maxPermutationsPerLine, maxPermutationsPerFile, big.NewInt(1), config)
		return
	}

	var choice string
	fmt.Print("Do you want to generate a single package or a range of packages? (single/range): ")
	_, err = fmt.Scan(&choice)
	if err != nil {
		fmt.Printf("Invalid input: %v\n", err)
		return
	}

	if choice == "single" {
		var packageFileNumberStr string
		fmt.Print("Enter the package number to generate: ")
		_, err = fmt.Scan(&packageFileNumberStr)
		if err != nil {
			fmt.Printf("Invalid input: %v\n", err)
			return
		}

		packageFileNumber := new(big.Int)
		packageFileNumber, ok := packageFileNumber.SetString(packageFileNumberStr, 10)
		if !ok || packageFileNumber.Cmp(big.NewInt(1)) < 0 || packageFileNumber.Cmp(totalPackageFiles) > 0 {
			fmt.Printf("Invalid package number: %v\n", err)
			return
		}

		calculatePermutationRanges(length, maxPermutationsPerLine, maxPermutationsPerFile, packageFileNumber, config)
	} else if choice == "range" {
		var startPackageFileNumberStr, endPackageFileNumberStr string
		fmt.Print("Enter the start package number to generate: ")
		_, err = fmt.Scan(&startPackageFileNumberStr)
		if err != nil {
			fmt.Printf("Invalid input: %v\n", err)
			return
		}

		fmt.Print("Enter the end package number to generate: ")
		_, err = fmt.Scan(&endPackageFileNumberStr)
		if err != nil {
			fmt.Printf("Invalid input: %v\n", err)
			return
		}

		startPackageFileNumber := new(big.Int)
		startPackageFileNumber, ok := startPackageFileNumber.SetString(startPackageFileNumberStr, 10)
		if !ok || startPackageFileNumber.Cmp(big.NewInt(1)) < 0 || startPackageFileNumber.Cmp(totalPackageFiles) > 0 {
			fmt.Printf("Invalid start package number: %v\n", err)
			return
		}

		endPackageFileNumber := new(big.Int)
		endPackageFileNumber, ok = endPackageFileNumber.SetString(endPackageFileNumberStr, 10)
		if !ok || endPackageFileNumber.Cmp(big.NewInt(1)) < 0 || endPackageFileNumber.Cmp(totalPackageFiles) > 0 || endPackageFileNumber.Cmp(startPackageFileNumber) < 0 {
			fmt.Printf("Invalid end package number: %v\n", err)
			return
		}

		for packageFileNumber := new(big.Int).Set(startPackageFileNumber); packageFileNumber.Cmp(endPackageFileNumber) <= 0; packageFileNumber.Add(packageFileNumber, big.NewInt(1)) {
			calculatePermutationRanges(length, maxPermutationsPerLine, maxPermutationsPerFile, packageFileNumber, config)
		}
	} else {
		fmt.Println("Invalid choice. Please enter 'single' or 'range'.")
	}
}
