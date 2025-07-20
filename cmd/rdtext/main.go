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

var processedCounter = big.NewInt(0)
var rateCounter = big.NewInt(0)
var charRepo runelib.CharacterRepo
var mapMutex sync.Mutex
var sentenceMap map[int][]liberdatabase.SentenceRecord
var sentenceChan chan Sentence

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
	sentenceChan = make(chan Sentence, 16384) // Increased buffer size
	sentenceMap = make(map[int][]liberdatabase.SentenceRecord)
	charRepo = *runelib.NewCharacterRepo()
	sheetName := "Worksheet"

	// We are going to put timer to see how many we have processed.
	processedTicker := time.NewTicker(time.Minute)
	defer processedTicker.Stop()

	go func() {
		for range processedTicker.C {
			fmt.Printf("Rate: %s/min - %s items remaining\n", rateCounter.String(), processedCounter.String())
			rateCounter.SetInt64(int64(0))
		}
	}()

	// Define command-line flags
	inputFile := flag.String("input", "", "Path to the input Excel file")
	reverseWords := flag.Bool("reverse", false, "Reverse the words in the sentence")

	// Parse the flags
	flag.Parse()

	// Making sure the tables are created.
	_, createErr := liberdatabase.InitTables()
	if createErr != nil {
		return
	}

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

	sheetInfo, colInfo, excelErr := getColInformation(f, sheetName, *inputFile)
	if excelErr != nil {
		fmt.Printf("Failed to get column info: %v", excelErr)
		return
	}

	db, _ := liberdatabase.InitConnection()
	_ = liberdatabase.AddSheetInformation(db, sheetInfo)
	_ = liberdatabase.CloseConnection(db)

	// Print the column information
	fmt.Printf("%v\n", colInfo)
	n := big.NewInt(1)
	for i := 0; i < len(colInfo); i++ {
		n.Mul(n, big.NewInt(int64(colInfo[i].RowCounts)))
	}
	fmt.Printf("Total combinations: %s\n", n.String())
	processedCounter.Set(n)

	// Initialize a string builder
	var builder strings.Builder

	// Threading for sentence processing.
	var wg sync.WaitGroup

	numWorkers := runtime.NumCPU() // Adjusted number of workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go insertSentenceToDB(i, &wg, *reverseWords)
	}

	// Call permuteCols with the provided output file name
	go func() {
		permuteErr := permuteCols(f, sheetName, colInfo, builder, 0, filepath.Base(*inputFile))
		if permuteErr != nil {
			fmt.Printf("Failed to permute cols: %v", permuteErr)
		}
		close(sentenceChan)
	}()

	wg.Wait()

	closeErr := f.Close()
	if closeErr != nil {
		log.Fatalf("Failed to close the Excel file: %v", closeErr)
	}
}

// permuteCols permutes the columns in the Excel file and writes the sentences to the output file.
func permuteCols(f *excelize.File, sheetName string, cols []ColInformation,
	builder strings.Builder, currentColIdx int, filename string) (err error) {
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
	}

	return nil
}

// insertSentenceToDB processes sentences from a channel, computes their properties, and inserts them into the database.
// workerId identifies the worker handling this operation.
// wg is the WaitGroup for synchronizing completion of concurrent workers.
// the reverseWords parameter indicates whether the sentences should be reversed during processing.
func insertSentenceToDB(workerId int, wg *sync.WaitGroup, reverseWords bool) {
	conn, connErr := liberdatabase.InitConnection()
	if connErr != nil {
		fmt.Printf("error initializing MySQL connection: %v", connErr)
	}

	for sentence := range sentenceChan {
		runeglish := runer.PrepLatinToRune(sentence.Content)
		runes := runer.TransposeLatinToRune(runeglish, reverseWords)
		gemSum := charRepo.CalculateGemSum(runes)
		gemSumBig := big.NewInt(int64(gemSum))

		if !sequences.IsPrime(gemSumBig) {
			incrementCounters()
			continue
		}

		mapMutex.Lock()
		record := liberdatabase.SentenceRecord{
			FileName:     sentence.FileName,
			DictSentence: sentence.Content,
			GemValue:     int64(gemSum),
			IsPrime:      true,
		}

		// Insert the record into the database
		if len(sentenceMap[workerId]) < 500 {
			sentenceMap[workerId] = append(sentenceMap[workerId], record)
		} else {
			err := liberdatabase.AddSentenceRecord(conn, sentenceMap[workerId])
			if err != nil {
				fmt.Printf("error inserting sentence record: %v", err)
			} else {
				sentenceMap[workerId] = []liberdatabase.SentenceRecord{}
			}
		}
		incrementCounters()
		mapMutex.Unlock()
	}

	mapMutex.Lock()
	if len(sentenceMap[workerId]) > 0 {
		err := liberdatabase.AddSentenceRecord(conn, sentenceMap[workerId])
		if err != nil {
			fmt.Printf("error inserting sentence record: %v", err)
		} else {
			sentenceMap[workerId] = []liberdatabase.SentenceRecord{}
		}
		incrementCountersByValue(len(sentenceMap[workerId]))
	}
	mapMutex.Unlock()

	closeError := liberdatabase.CloseConnection(conn)
	if closeError != nil {
		fmt.Printf("error closing MySQL connection: %v", closeError)
	}

	wg.Done()
}

// incrementCounters decrements the processedCounter by 1 and increments the rateCounter by 1.
func incrementCounters() {
	processedCounter.Sub(processedCounter, big.NewInt(1))
	rateCounter.Add(rateCounter, big.NewInt(1))
}

// incrementCountersByValue adjusts `processedCounter` by decrementing and `rateCounter` by incrementing the given value.
func incrementCountersByValue(value int) {
	processedCounter.Sub(processedCounter, big.NewInt(int64(value)))
	rateCounter.Add(rateCounter, big.NewInt(int64(value)))
}

// cloneStringBuilder clones the given strings.Builder and returns a new instance with the same content.
func cloneStringBuilder(sb *strings.Builder) *strings.Builder {
	// Create a new string builder
	newSb := &strings.Builder{}
	// Write the content of the original builder to the new builder
	newSb.WriteString(sb.String())
	return newSb
}

// getColInformation gets the column information from the given Excel file and sheet name.
func getColInformation(f *excelize.File, sheetName, fileName string) (sheet liberdatabase.SheetInformation, cols []ColInformation, err error) {
	sheetCol := intToExcelColumn(1)
	sheetRow := 1
	sheetCellValue, _ := f.GetCellValue(sheetName, fmt.Sprintf("%s%d", sheetCol, sheetRow))

	sheetInformation := liberdatabase.SheetInformation{
		FileName: fileName,
		Text:     sheetCellValue,
	}

	// Get the column information
	cols = make([]ColInformation, 0)

	// Get all the columns in the sheet
	columns, err := f.GetCols(sheetName)
	if err != nil {
		return sheetInformation, nil, fmt.Errorf("Failed to get columns: %v\n", err)
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
				return sheetInformation, nil, fmt.Errorf("Failed to get cell value: %v\n", cellError)
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

	return sheetInformation, cols, nil
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
