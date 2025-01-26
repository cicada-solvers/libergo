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
	"strconv"
	"strings"
	"sync"
)

type Config struct {
	NumWorkers             int   `json:"num_workers"`
	MaxPermutationsPerLine int64 `json:"max_permutations_per_line"`
	MaxPermutationsPerFile int64 `json:"max_permutations_per_file"`
	MaxFilesPerZip         int64 `json:"max_files_per_zip"`
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

type ZipWriter struct {
	MaxFilesPerZip       int64
	currentZipFile       *zip.Writer
	currentZipFileHandle *os.File
	currentFileCount     int64
	zipFileIndex         int
	mu                   sync.Mutex
}

func NewZipWriter(maxFilesPerZip int64) *ZipWriter {
	return &ZipWriter{
		MaxFilesPerZip: maxFilesPerZip,
		zipFileIndex:   0,
	}
}

func (zw *ZipWriter) AddFileToZip(folder, filePath string) error {
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

	w, err := zw.currentZipFile.Create(filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("error creating zip entry: %v", err)
	}

	bufferedWriter := bufio.NewWriter(w)
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

func (zw *ZipWriter) createNewZip(folder string) error {
	zipFileName := fmt.Sprintf("package_%d.zip", zw.zipFileIndex)
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
	zw.zipFileIndex++
	return nil
}

func calculatePermutationRanges(length int, maxPermutationsPerLine, maxPermutationsPerFile int64) {
	config, err := loadConfig("appsettings.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

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
		for i := int64(0); ; i++ {
			start := new(big.Int).Mul(big.NewInt(i), big.NewInt(maxPermutationsPerLine*maxPermutationsPerFile))
			if start.Cmp(totalPermutations) >= 0 {
				break
			}
			fileChan <- i
		}
		close(fileChan)
	}()

	zipWriter := NewZipWriter(config.MaxFilesPerZip)

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
				filePath := filepath.Join(folder, fileName)
				file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Printf("Error opening file: %v\n", err)
					return
				}

				bufferedWriter := bufio.NewWriter(file)
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

				if err := zipWriter.AddFileToZip(folder, filePath); err != nil {
					fmt.Printf("Error adding file to zip: %v\n", err)
				}
			}
		}()
	}

	wg.Wait()

	if err := zipWriter.closeCurrentZip(); err != nil {
		fmt.Printf("Error closing zip file: %v\n", err)
	}
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

	config, err := loadConfig("appsettings.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	maxPermutationsPerLine := config.MaxPermutationsPerLine
	maxPermutationsPerFile := config.MaxPermutationsPerFile
	calculatePermutationRanges(length, maxPermutationsPerLine, maxPermutationsPerFile)
}
