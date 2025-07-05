package main

import (
	runelib "characterrepo"
	"flag"
	"fmt"
	"liberdatabase"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"runer"
	"runtime"
	"sequences"
	"strings"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"
)

var fileMutex sync.Mutex
var processedCounter = big.NewInt(0)
var rateCounter = big.NewInt(0)
var charRepo runelib.CharacterRepo
var mapMutex sync.Mutex
var sentences map[int][]liberdatabase.SentenceRecord

// ColInformation represents the column information with its name and row counts.
type ColInformation struct {
	ColName   string
	RowCounts int
}

// Sentence represents a sentence with its content, output file name, and column index.
type Sentence struct {
	FileName    string
	Content     string
	Output      string
	ColumnIndex int
}

// main is the entry point of the program.
func main() {
	sentences = make(map[int][]liberdatabase.SentenceRecord)
	charRepo = *runelib.NewCharacterRepo()
	sheetName := "Worksheet"

	// We are going to put timer to see how many we have processed.
	processedTicker := time.NewTicker(time.Minute)
	defer processedTicker.Stop()

	go func() {
		for range processedTicker.C {
			fmt.Printf("Rate: %s/min - Processed %s items\n", rateCounter.String(), processedCounter.String())
			rateCounter.SetInt64(int64(0))
		}
	}()

	// Define command-line flags
	inputFile := flag.String("input", "", "Path to the input Excel file")
	outputFile := flag.String("output", "", "Path to the output file")
	isCreating := flag.Bool("create", false, "Create a new database")

	// Parse the flags
	flag.Parse()

	// Making sure the tables are created.
	_, createErr := liberdatabase.InitTables()
	if createErr != nil {
		return
	}

	if *isCreating {
		// Check if the input file is provided
		if *inputFile == "" {
			log.Fatalf("Input file is required")
		}

		// Process the input file
		infoFile, _ := os.Stat(*inputFile)
		fmt.Printf("Processing file: %s\n", infoFile.Name())

		// Open the Excel file
		f, fileErr := excelize.OpenFile(*inputFile)
		if fileErr != nil {
			log.Fatalf("Failed to open the Excel file: %v", fileErr)
		}

		colInfo, excelErr := getColInformation(f, sheetName)
		if excelErr != nil {
			fmt.Printf("Failed to get column info: %v", excelErr)
			return
		}

		// Print the column information
		fmt.Printf("%v\n", colInfo)
		n := big.NewInt(1)
		for i := 0; i < len(colInfo); i++ {
			n.Mul(n, big.NewInt(int64(colInfo[i].RowCounts)))
		}
		fmt.Printf("Total combinations: %s\n", n.String())

		// Initialize a strings.Builder
		var builder strings.Builder

		// Call permuteCols with the provided output file name
		permuteErr := permuteCols(f, sheetName, colInfo, builder, 0, filepath.Base(*inputFile))
		if permuteErr != nil {
			fmt.Printf("Failed to permute cols: %v", permuteErr)
		}

		closeErr := f.Close()
		if closeErr != nil {
			log.Fatalf("Failed to close the Excel file: %v", closeErr)
		}
	} else {
		// Check if the output file already exists
		if _, err := os.Stat(*outputFile); err == nil {
			log.Fatalf("Output file already exists")
		}

		// Check if the output file is provided
		if *outputFile == "" {
			log.Fatalf("Output file is required")
		}

		// Now we are going to remove the million records from the database.
		conn, connErr := liberdatabase.InitConnection()
		if connErr != nil {
			fmt.Printf("error initializing MySQL connection: %v", connErr)
		}

		// Get the top million sentence records
		records, getErr := liberdatabase.GetTopMillionSentenceRecords(conn, filepath.Base(*inputFile))
		if getErr != nil {
			fmt.Printf("error getting top million sentence records: %v", getErr)
		}

		var wg sync.WaitGroup
		sentenceChan := make(chan Sentence, 16384) // Increased buffer size

		go func() {
			for _, record := range records {
				// Create a new Sentence instance
				sentence := Sentence{
					FileName:    record.FileName,
					Content:     record.DictSentence,
					Output:      *outputFile,
					ColumnIndex: 0, // Set the column index as needed
				}
				sentenceChan <- sentence
			}
			close(sentenceChan)
		}()

		numWorkers := runtime.NumCPU() // Adjusted number of workers
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go calculateGemSumAndWriteToFile(sentenceChan, &wg)
		}

		wg.Wait()

		// Remove the million records from the database
		removeErr := liberdatabase.RemoveMillionSentenceRecords(conn, records)
		if removeErr != nil {
			fmt.Printf("error removing million sentence records: %v", removeErr)
		}
	}
}

// permuteCols permutes the columns in the Excel file and writes the sentences to the output file.
func permuteCols(f *excelize.File, sheetName string, cols []ColInformation, builder strings.Builder, currentColIdx int, filename string) (err error) {
	localBuilder := cloneStringBuilder(&builder)

	if currentColIdx < (len(cols) - 1) {
		for i := 3; i < cols[currentColIdx].RowCounts+3; i++ {
			localBuilder = cloneStringBuilder(&builder)
			cellValue, cellError := f.GetCellValue(sheetName, fmt.Sprintf("%s%d", cols[currentColIdx].ColName, i))
			if cellError != nil {
				return fmt.Errorf("Failed to get cell value: %v\n", cellError)
			}

			spacer := " "
			if currentColIdx == 0 {
				spacer = ""
			}

			localBuilder.WriteString(spacer + cellValue)

			permuteErr := permuteCols(f, sheetName, cols, *localBuilder, currentColIdx+1, filename)
			if permuteErr != nil {
				return fmt.Errorf("Failed to permute columns: %v\n", permuteErr)
			}
		}
	} else {
		var wg sync.WaitGroup
		sentenceChan := make(chan Sentence, 16384) // Increased buffer size

		go func() {
			for i := 3; i < cols[currentColIdx].RowCounts+3; i++ {
				cellValue, cellError := f.GetCellValue(sheetName, fmt.Sprintf("%s%d", cols[currentColIdx].ColName, i))
				if cellError != nil {
					fmt.Printf("Failed to get cell value: %v\n", cellError)
				}

				localBuilder = cloneStringBuilder(&builder)

				localBuilder.WriteString(" " + cellValue)

				sentence := Sentence{
					FileName:    filename,
					Content:     localBuilder.String(),
					ColumnIndex: currentColIdx,
				}

				sentenceChan <- sentence
			}
			close(sentenceChan)
		}()

		numWorkers := runtime.NumCPU() // Adjusted number of workers
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go insertSentenceToDB(i, sentenceChan, &wg)
		}

		wg.Wait()
	}

	return nil
}

func insertSentenceToDB(workerId int, sentChan chan Sentence, wg *sync.WaitGroup) {
	// Create a new SentenceRecord
	defer wg.Done()

	conn, connErr := liberdatabase.InitConnection()
	if connErr != nil {
		fmt.Printf("error initializing MySQL connection: %v", connErr)
	}

	for sentence := range sentChan {
		mapMutex.Lock()
		record := liberdatabase.SentenceRecord{
			FileName:     sentence.FileName,
			DictSentence: sentence.Content,
		}

		// Insert the record into the database
		if len(sentences[workerId]) < 500 {
			sentences[workerId] = append(sentences[workerId], record)
		} else {
			err := liberdatabase.AddSentenceRecord(conn, sentences[workerId])
			if err != nil {
				fmt.Printf("error inserting sentence record: %v", err)
			} else {
				sentences[workerId] = []liberdatabase.SentenceRecord{}
			}
		}

		processedCounter.Add(processedCounter, big.NewInt(1))
		rateCounter.Add(rateCounter, big.NewInt(1))
		mapMutex.Unlock()
	}

	mapMutex.Lock()
	if len(sentences[workerId]) > 0 {
		err := liberdatabase.AddSentenceRecord(conn, sentences[workerId])
		if err != nil {
			fmt.Printf("error inserting sentence record: %v", err)
		} else {
			sentences[workerId] = []liberdatabase.SentenceRecord{}
		}
	}
	mapMutex.Unlock()

	closeError := liberdatabase.CloseConnection(conn)
	if closeError != nil {
		fmt.Printf("error closing MySQL connection: %v", closeError)
	}
}

// calculateGemSumAndWriteToFile calculates the probability of a sentence being a valid English sentence and writes it to the output file.
func calculateGemSumAndWriteToFile(sentChan chan Sentence, wg *sync.WaitGroup) {
	one := big.NewInt(1)

	defer wg.Done()

	for sentence := range sentChan {
		runeglish := runer.PrepLatinToRune(sentence.Content)
		runes := runer.TransposeLatinToRune(runeglish)

		gemSum := charRepo.CalculateGemSum(runes)
		gemSumBig := big.NewInt(int64(gemSum))

		if sequences.IsPrime(gemSumBig) {
			fmt.Printf("Prime Sentence: %s\n", sentence.Content)

			// Write the content to the output file
			outputBytes := []byte(sentence.Content + "\n\n")

			WriteContentsToOutputFile(sentence.Output, outputBytes)
		}

		processedCounter.Add(processedCounter, one)
		rateCounter.Add(rateCounter, one)
	}
}

func WriteContentsToOutputFile(outputfile string, outputBytes []byte) {
	for {
		fileMutex.Lock()
		file, openError := os.OpenFile(outputfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if openError != nil {
			fmt.Printf("Failed to open file: %v\n", openError)
			fileMutex.Unlock()
			time.Sleep(100 * time.Millisecond) // Wait before retrying
			continue
		}

		_, writeErr := file.Write(outputBytes)
		if writeErr != nil {
			fmt.Printf("Failed to write to file: %v\n", writeErr)
			err := file.Close()
			if err != nil {
				fmt.Printf("Failed to close file: %v\n", err)
			}
			fileMutex.Unlock()
			time.Sleep(100 * time.Millisecond) // Wait before retrying
			continue
		}

		closeError := file.Close()
		if closeError != nil {
			fmt.Printf("Failed to close file: %v\n", closeError)
		}
		fileMutex.Unlock()
		break
	}
}

// cloneStringBuilder clones the given strings.Builder and returns a new instance with the same content.
func cloneStringBuilder(sb *strings.Builder) *strings.Builder {
	// Create a new strings.Builder
	newSb := &strings.Builder{}
	// Write the content of the original builder to the new builder
	newSb.WriteString(sb.String())
	return newSb
}

// getColInformation gets the column information from the given Excel file and sheet name.
func getColInformation(f *excelize.File, sheetName string) (cols []ColInformation, err error) {
	// Get the column information
	cols = make([]ColInformation, 0)

	// Get all the columns in the sheet
	columns, err := f.GetCols(sheetName)
	if err != nil {
		return nil, fmt.Errorf("Failed to get columns: %v\n", err)
	}

	for colId, colValues := range columns {
		// Get the row count
		rowCounts := len(colValues)
		colName := intToExcelColumn(colId + 1)
		colInfo := ColInformation{
			ColName:   colName,
			RowCounts: 0,
		}

		for i := 3; i < rowCounts+5; i++ {
			cellValue, cellError := f.GetCellValue(sheetName, fmt.Sprintf("%s%d", colName, i))
			if cellError != nil {
				return nil, fmt.Errorf("Failed to get cell value: %v\n", cellError)
			}

			if cellValue != "" {
				colInfo.RowCounts++
			} else {
				break
			}
		}

		if colInfo.RowCounts <= 0 {
			colInfo.RowCounts = 1
		}

		cols = append(cols, colInfo)
	}

	return cols, nil
}

// intToExcelColumn converts the given integer to an Excel column name.
func intToExcelColumn(n int) string {
	column := ""
	for n > 0 {
		n-- // Adjust for 1-based indexing
		column = string('A'+(n%26)) + column
		n /= 26
	}
	return column
}
