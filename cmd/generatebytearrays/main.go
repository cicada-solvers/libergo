package main

import (
	"archive/zip"
	"bufio"
	"compress/flate"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Config represents the configuration for the application
type Config struct {
	NumWorkers             int   `json:"num_workers"`
	MaxPermutationsPerLine int64 `json:"max_permutations_per_line"`
	MaxPermutationsPerFile int64 `json:"max_permutations_per_file"`
	MaxFilesPerZip         int64 `json:"max_files_per_zip"`
}

// loadConfig loads the configuration from the specified file
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

// ZipWriter is a struct that writes files to a zip archive
type ZipWriter struct {
	MaxFilesPerZip       int64
	currentZipFile       *zip.Writer
	currentZipFileHandle *os.File
	currentFileCount     int64
	mu                   sync.Mutex
	totalZipFiles        *big.Int
	zipFileNumber        *big.Int
	ArrayLength          int
}

// NewZipWriter creates a new ZipWriter
func NewZipWriter(maxFilesPerZip int64) *ZipWriter {
	return &ZipWriter{
		MaxFilesPerZip: maxFilesPerZip,
	}
}

// addFileToZip adds a file to the zip archive
func (zw *ZipWriter) addFileToZip(folder, filePath string) error {
	zw.mu.Lock()
	defer zw.mu.Unlock()

	if zw.currentZipFile == nil || zw.currentFileCount >= zw.MaxFilesPerZip {
		if zw.currentZipFile != nil {
			if err := zw.closeCurrentZip(); err != nil {
				return err
			}
		}
		if err := zw.createNewZip(folder); err != nil {
			return err
		}
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Truncate the file name if it exceeds 255 characters
	fileName := filepath.Base(filePath)
	if len(fileName) > 255 {
		ext := filepath.Ext(fileName)
		name := fileName[:255-len(ext)]
		fileName = name + ext
	}

	w, err := zw.currentZipFile.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating zip entry: %v", err)
	}

	bufferedWriter := bufio.NewWriterSize(w, 64*1024) // Increase buffer size to 64KB
	if _, err := io.Copy(bufferedWriter, file); err != nil {
		return fmt.Errorf("error writing file to zip: %v", err)
	}
	bufferedWriter.Flush()

	zw.currentFileCount++

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("error deleting file: %v", err)
	}

	return nil
}

// closeCurrentZip closes the current zip archive
func (zw *ZipWriter) closeCurrentZip() error {
	if zw.currentZipFile == nil {
		return nil
	}
	fmt.Printf("Closing zip file\n")
	if err := zw.currentZipFile.Close(); err != nil {
		return fmt.Errorf("error closing zip file: %v", err)
	}
	if err := zw.currentZipFileHandle.Close(); err != nil {
		return fmt.Errorf("error closing zip file handle: %v", err)
	}
	zw.currentZipFile = nil
	zw.currentZipFileHandle = nil
	return nil
}

// createNewZip creates a new zip archive
func (zw *ZipWriter) createNewZip(folder string) error {
	zipFileName := fmt.Sprintf("package_l%d_%s_of_%s.zip", zw.ArrayLength, zw.zipFileNumber.String(), zw.totalZipFiles.String())
	zipFilePath := filepath.Join(folder, zipFileName)
	fmt.Printf("Creating new zip file: %s\n", zipFilePath)
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return fmt.Errorf("error creating zip file: %v", err)
	}

	zw.currentZipFile = zip.NewWriter(zipFile)
	zw.currentZipFile.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})
	zw.currentZipFileHandle = zipFile
	zw.currentFileCount = 0
	return nil
}

// calculatePermutationRanges calculates the permutation ranges for the specified length
func calculatePermutationRanges(length int, maxPermutationsPerLine, maxPermutationsPerFile int64, zipFileNumber *big.Int) {
	config, err := loadConfig("appsettings.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	totalZipFiles, err := calculateNumberOfZipFiles(length, maxPermutationsPerLine, maxPermutationsPerFile, config.MaxFilesPerZip)
	if err != nil {
		fmt.Printf("Error calculating number of zip files: %v\n", err)
		return
	}

	fmt.Printf("Total number of zip files: %s\n", totalZipFiles.String())

	totalPermutations := big.NewInt(1)
	for i := 0; i < length; i++ {
		totalPermutations.Mul(totalPermutations, big.NewInt(256))
	}

	folder := fmt.Sprintf("%010d", length)
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		fmt.Printf("Error creating folder: %v\n", err)
		return
	}

	var wg sync.WaitGroup
	fileChan := make(chan int64)

	go func() {
		startFile := new(big.Int).Mul(new(big.Int).Sub(zipFileNumber, big.NewInt(1)), big.NewInt(config.MaxFilesPerZip))
		endFile := new(big.Int).Add(startFile, big.NewInt(config.MaxFilesPerZip))
		for i := new(big.Int).Set(startFile); i.Cmp(endFile) < 0; i.Add(i, big.NewInt(1)) {
			start := new(big.Int).Mul(i, big.NewInt(maxPermutationsPerLine*maxPermutationsPerFile))
			if start.Cmp(totalPermutations) >= 0 {
				break
			}
			fileChan <- i.Int64()
		}
		close(fileChan)
	}()

	zipWriter := NewZipWriter(config.MaxFilesPerZip)
	zipWriter.totalZipFiles = totalZipFiles
	zipWriter.ArrayLength = length
	zipWriter.zipFileNumber = zipFileNumber

	for i := 0; i < config.NumWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range fileChan {
				start := new(big.Int).Mul(big.NewInt(i), big.NewInt(maxPermutationsPerLine*maxPermutationsPerFile))
				end := new(big.Int).Add(start, big.NewInt(maxPermutationsPerLine*maxPermutationsPerFile))
				if end.Cmp(totalPermutations) > 0 {
					end = totalPermutations
				}

				fileName := fmt.Sprintf("permutations_%d.txt", i)
				if len(fileName) > 255 {
					ext := filepath.Ext(fileName)
					name := fileName[:255-len(ext)]
					fileName = name + ext
				}
				filePath := filepath.Join(folder, fileName)
				file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Printf("Error opening file: %v\n", err)
					return
				}

				bufferedWriter := bufio.NewWriter(file)
				_, err = bufferedWriter.WriteString(fmt.Sprintf("%d\n", maxPermutationsPerLine))
				if err != nil {
					fmt.Printf("Error writing to file: %v\n", err)
				}

				for j := int64(0); j < maxPermutationsPerFile; j++ {
					lineStart := new(big.Int).Add(start, big.NewInt(j*maxPermutationsPerLine))
					lineEnd := new(big.Int).Add(lineStart, big.NewInt(maxPermutationsPerLine))
					if lineEnd.Cmp(totalPermutations) > 0 {
						lineEnd = totalPermutations
					}

					startArray := indexToArray(lineStart, length)
					endArray := indexToArray(new(big.Int).Sub(lineEnd, big.NewInt(1)), length)

					_, err = bufferedWriter.WriteString(fmt.Sprintf("%s-%s\n", arrayToString(startArray), arrayToString(endArray)))
					if err != nil {
						fmt.Printf("Error writing to file: %v\n", err)
					}

					if lineEnd.Cmp(totalPermutations) == 0 {
						break
					}
				}
				bufferedWriter.Flush()

				err = file.Close()
				if err != nil {
					fmt.Printf("Error closing file: %v\n", err)
				}

				if err := zipWriter.addFileToZip(folder, filePath); err != nil {
					fmt.Printf("Error adding file to zip: %v\n", err)
				}
			}
		}()
	}

	wg.Wait()

	if err := zipWriter.closeCurrentZip(); err != nil {
		fmt.Printf("Error closing zip file: %v\n", err)
	}

	// Check if the folder exists in processhashes
	destinationFolder := filepath.Join("..", "processhashes", folder)
	if _, err := os.Stat(destinationFolder); os.IsNotExist(err) {
		// Move the folder to the processhashes folder
		if err := os.Rename(folder, destinationFolder); err != nil {
			fmt.Printf("Error moving folder: %v\n", err)
		} else {
			fmt.Printf("Folder moved to %s\n", destinationFolder)
		}
	} else {
		// Move the zip file to the existing folder in processhashes
		zipFileName := fmt.Sprintf("package_l%d_%s_of_%s.zip", length, zipFileNumber.String(), totalZipFiles.String())
		zipFilePath := filepath.Join(folder, zipFileName)
		newZipFilePath := filepath.Join(destinationFolder, zipFileName)
		if err := os.Rename(zipFilePath, newZipFilePath); err != nil {
			fmt.Printf("Error moving zip file: %v\n", err)
		} else {
			fmt.Printf("Zip file moved to %s\n", newZipFilePath)
		}
		// Remove the folder under generatebytearrays
		if err := os.RemoveAll(folder); err != nil {
			fmt.Printf("Error removing folder: %v\n", err)
		} else {
			fmt.Printf("Folder %s removed\n", folder)
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

// calculateNumberOfZipFiles calculates the number of zip files required to store all permutations
func calculateNumberOfZipFiles(length int, maxPermutationsPerLine, maxPermutationsPerFile, maxFilesPerZip int64) (*big.Int, error) {
	totalPermutations := big.NewInt(1)
	for i := 0; i < length; i++ {
		totalPermutations.Mul(totalPermutations, big.NewInt(256))
	}

	permutationsPerFile := big.NewInt(maxPermutationsPerLine * maxPermutationsPerFile)
	totalFiles := new(big.Int).Div(totalPermutations, permutationsPerFile)
	if new(big.Int).Mod(totalPermutations, permutationsPerFile).Cmp(big.NewInt(0)) != 0 {
		totalFiles.Add(totalFiles, big.NewInt(1))
	}

	totalZipFiles := new(big.Int).Div(totalFiles, big.NewInt(maxFilesPerZip))
	if new(big.Int).Mod(totalFiles, big.NewInt(maxFilesPerZip)).Cmp(big.NewInt(0)) != 0 {
		totalZipFiles.Add(totalZipFiles, big.NewInt(1))
	}

	return totalZipFiles, nil
}

// main is the entry point for the application
func main() {
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
	maxPermutationsPerFile := config.MaxPermutationsPerFile

	totalZipFiles, err := calculateNumberOfZipFiles(length, maxPermutationsPerLine, maxPermutationsPerFile, config.MaxFilesPerZip)
	if err != nil {
		fmt.Printf("Error calculating number of zip files: %v\n", err)
		return
	}

	fmt.Printf("Total number of zip files: %s\n", totalZipFiles.String())

	var zipFileNumberStr string
	fmt.Print("Enter the zip file number to generate: ")
	_, err = fmt.Scan(&zipFileNumberStr)
	if err != nil {
		fmt.Printf("Invalid input: %v\n", err)
		return
	}

	zipFileNumber := new(big.Int)
	zipFileNumber, ok := zipFileNumber.SetString(zipFileNumberStr, 10)
	if !ok || zipFileNumber.Cmp(big.NewInt(1)) < 0 || zipFileNumber.Cmp(totalZipFiles) > 0 {
		fmt.Printf("Invalid zip file number: %v\n", err)
		return
	}

	calculatePermutationRanges(length, maxPermutationsPerLine, maxPermutationsPerFile, zipFileNumber)
}
