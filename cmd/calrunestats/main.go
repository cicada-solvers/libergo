package main

import (
	"bufio"
	runelib "characterrepo"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type LineStatistics struct {
	KeyWord      string
	Text         string
	DoubletCount int
	Statistics   []RuneStatistic
}

type RuneStatistic struct {
	Rune       string
	Percentage float64
}

func main() {
	// Initialize the character repository
	charRepo := runelib.NewCharacterRepo()

	// Define the flag for the directory path
	dirPath := flag.String("dir", "./your-directory-path", "Path to the directory containing text files")
	flag.Parse()

	// Walk through the directory and its subdirectories
	err := filepath.WalkDir(*dirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Open the file
		f, openErr := os.Open(path)
		if openErr != nil {
			fmt.Printf("Failed to open file %s: %v\n", path, openErr)
			return nil
		}
		defer func(f *os.File) {
			closeErr := f.Close()
			if closeErr != nil {
				fmt.Printf("Error closing file %s: %v\n", path, err)
			} else {
				fmt.Printf("Closed file %s successfully\n", path)
			}

			// Delete the file after processing
			//removeErr := os.Remove(path)
			//if removeErr != nil {
			//fmt.Printf("Error deleting file %s: %v\n", path, removeErr)
			//} else {
			//fmt.Printf("Deleted file %s successfully\n", path)
			//}
		}(f)

		// Read the file line by line
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			var lineStatistics []RuneStatistic
			line := scanner.Text()
			parts := strings.Split(line, " : ")
			if len(parts) > 1 {
				partTwoArray := strings.Split(parts[1], "")
				totalCount := getTotalRunes(partTwoArray)
				if totalCount > 0 {
					for _, runeChar := range charRepo.GetGematriaRunes() {
						count := getCountOfParticularRune(partTwoArray, runeChar)
						percentage := (float64(count) / float64(totalCount)) * 100
						statistic := RuneStatistic{
							Rune:       runeChar,
							Percentage: percentage,
						}

						lineStatistics = append(lineStatistics, statistic)
					}

					lineStatisticsAll := LineStatistics{
						KeyWord:      parts[0],
						Text:         parts[1],
						DoubletCount: charRepo.GetDoubletCount(parts[1], charRepo.GetGematriaRunes()),
						Statistics:   sortLineStatistics(lineStatistics),
					}

					writeErr := writeResultsToFile(len(charRepo.GetGematriaRunes()), lineStatisticsAll, path+".csv")
					if writeErr != nil {
						fmt.Printf("Error writing results to file %s: %v\n", path+".csv", writeErr)
					} else {
						fmt.Printf("Results written to %s successfully\n", path+".csv")
					}
				}
			}
		}

		// Check for errors during scanning
		if scanErr := scanner.Err(); scanErr != nil {
			fmt.Printf("Error reading file %s: %v\n", path, scanErr)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the directory: %v\n", err)
	}
}

func getTotalRunes(line []string) int {
	charRepo := runelib.NewCharacterRepo()
	total := 0
	for _, str := range line {
		if charRepo.IsRune(str, false) {
			total += len(str)
		}
	}
	return total
}

func getCountOfParticularRune(line []string, rune string) int {
	charRepo := runelib.NewCharacterRepo()
	count := 0
	for _, str := range line {
		if charRepo.IsRune(str, false) && str == rune {
			count++
		}
	}
	return count
}

// SortLineStatistics sorts the lineStatistics map by percentage in descending order.
func sortLineStatistics(statistics []RuneStatistic) []RuneStatistic {
	// Convert the map to a slice of RuneStatistic
	var sortedStats []RuneStatistic
	for _, percentage := range statistics {
		sortedStats = append(sortedStats, percentage)
	}

	// Sort the slice by percentage in descending order
	sort.Slice(sortedStats, func(i, j int) bool {
		return sortedStats[i].Percentage > sortedStats[j].Percentage
	})

	return sortedStats
}

func writeResultsToFile(alphabetCount int, results LineStatistics, outputPath string) error {
	// Open or create the result CSV file
	resultFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Error opening/creating result.csv: %v\n", err)
		return err
	}
	defer func(resultFile *os.File) {
		closeErr := resultFile.Close()
		if closeErr != nil {
			fmt.Printf("Error closing result.csv: %v\n", closeErr)
		} else {
			fmt.Println("Closed result.csv successfully")
		}
	}(resultFile)

	// Create a CSV writer
	csvWriter := csv.NewWriter(resultFile)
	defer csvWriter.Flush()

	// Write the header if the file is empty
	fileInfo, _ := resultFile.Stat()
	if fileInfo.Size() == 0 {
		header := []string{"KeyWord", "DoubletCount"}

		for i := 0; i < alphabetCount; i++ {
			header = append(header, fmt.Sprintf("Rune%d", i+1))
		}

		if headerErr := csvWriter.Write(header); headerErr != nil {
			fmt.Printf("Error writing header to CSV: %v\n", headerErr)
			return headerErr
		}
	}

	var record []string
	record = append(record, results.KeyWord)
	record = append(record, fmt.Sprintf("%d", results.DoubletCount))

	// Write the data
	for _, stat := range results.Statistics {
		rText := stat.Rune
		pText := fmt.Sprintf("%.2f", stat.Percentage)
		aText := fmt.Sprintf("%s|%s", rText, pText)
		record = append(record, aText)
	}

	if csvErr := csvWriter.Write(record); csvErr != nil {
		fmt.Printf("Error writing record to CSV: %v\n", csvErr)
		return err
	}

	return nil
}
