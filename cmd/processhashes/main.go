package main

import (
	"archive/zip"
	"bufio"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/sha3"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Config struct {
	NumWorkers   int    `json:"num_workers"`
	ExistingHash string `json:"existing_hash"`
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

type Program struct {
	tasks chan []byte
}

func NewProgram() *Program {
	return &Program{
		tasks: make(chan []byte, 10000), // Increase buffer size
	}
}

func (p *Program) GenerateAllByteArrays(maxArrayLength int, startArray, stopArray []byte) {
	currentArray := make([]byte, len(startArray))
	copy(currentArray, startArray)
	p.GenerateByteArrays(maxArrayLength, 1, currentArray, stopArray)
	close(p.tasks)
}

func (p *Program) GenerateByteArrays(maxArrayLength, currentArrayLevel int, passedArray, stopArray []byte) bool {
	startForValue := int(passedArray[currentArrayLevel-1])

	if currentArrayLevel == maxArrayLength {
		currentArray := make([]byte, maxArrayLength)

		if passedArray != nil {
			copy(currentArray, passedArray)
		}

		for i := startForValue; i < 256; i++ {
			currentArray[currentArrayLevel-1] = byte(i)
			p.tasks <- append([]byte(nil), currentArray...) // Send a copy to avoid data race
			if compareArrays(currentArray, stopArray) == 0 {
				fmt.Printf("Stop Array Was: %v\n", stopArray)
				fmt.Printf("Finished processing: %v\n", currentArray)
				return false
			}
		}
	} else {
		currentArray := make([]byte, maxArrayLength)
		if passedArray != nil {
			copy(currentArray, passedArray)
		}
		for i := startForValue; i < 256; i++ {
			currentArray[currentArrayLevel-1] = byte(i)
			shouldContinue := p.GenerateByteArrays(maxArrayLength, currentArrayLevel+1, currentArray, stopArray)

			if shouldContinue == false {
				return false
			}

			// This resets that byte to zero of the next one up.
			currentArray[currentArrayLevel] = 0
		}
	}

	return true
}

func compareArrays(a, b []byte) int {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}
	if len(a) < len(b) {
		return -1
	} else if len(a) > len(b) {
		return 1
	}
	return 0
}

func processTasks(tasks chan []byte, wg *sync.WaitGroup, existingHash string, done chan struct{}, once *sync.Once) {
	defer wg.Done()

	// Open the file in append mode
	file, err := os.OpenFile("found_hashes.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}(file)

	buffer := make([]byte, 0, 4096) // Buffer for batching writes

	hashCount := 0
	taskLen := 0
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	go func() {
		colors := []string{"\033[31m", "\033[32m", "\033[33m", "\033[34m", "\033[35m", "\033[36m"} // Red, Green, Yellow, Blue, Magenta, Cyan
		colorIndex := 0
		for range ticker.C {

			fmt.Printf("%sHashes per minute: %d, Array size: %d\033[0m\n", colors[colorIndex], hashCount, taskLen)

			hashCount = 0
			colorIndex = (colorIndex + 1) % len(colors)
		}
	}()

	for {
		select {
		case task, ok := <-tasks:
			if !ok {
				// Write any remaining data in the buffer
				if len(buffer) > 0 {
					if _, err := file.Write(buffer); err != nil {
						fmt.Printf("Error writing to file: %v\n", err)
					}
				}
				return
			}
			taskLen = len(task)
			hashes := generateHashes(task)
			for hashName, hash := range hashes {
				hashCount++
				if hash == existingHash {
					// Convert byte array to comma-separated string
					var taskStr string
					for i, b := range task {
						if i > 0 {
							taskStr += ","
						}
						taskStr += fmt.Sprintf("%d", b)
					}

					output := fmt.Sprintf("Match found: %s, Hash Name: %s, Byte Array: %s\n", taskStr, hashName, hex.EncodeToString(task))
					fmt.Print(output)
					buffer = append(buffer, output...)
					if len(buffer) >= 4096 {
						if _, err := file.Write(buffer); err != nil {
							fmt.Printf("Error writing to file: %v\n", err)
						}
						buffer = buffer[:0]
					}
					once.Do(func() { close(done) }) // Signal all goroutines to stop
					return
				}
			}
		case <-done:
			return
		}
	}
}

func generateHashes(data []byte) map[string]string {
	hashes := make(map[string]string)

	// SHA-512
	sha512Hash := sha512.Sum512(data)
	hashes["SHA-512"] = hex.EncodeToString(sha512Hash[:])

	// SHA3-512
	sha3Hash := sha3.Sum512(data)
	hashes["SHA3-512"] = hex.EncodeToString(sha3Hash[:])

	// Blake2b-512
	blake2bHash := blake2b.Sum512(data)
	hashes["Blake2b-512"] = hex.EncodeToString(blake2bHash[:])

	return hashes
}

func stringToByteArray(s string) ([]byte, []byte, error) {
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("invalid format: expected one hyphen separating start and end arrays")
	}

	startArray, err := convertToByteArray(parts[0])
	if err != nil {
		return nil, nil, fmt.Errorf("error converting start array: %v", err)
	}

	stopArray, err := convertToByteArray(parts[1])
	if err != nil {
		return nil, nil, fmt.Errorf("error converting stop array: %v", err)
	}

	return startArray, stopArray, nil
}

func convertToByteArray(part string) ([]byte, error) {
	subParts := strings.Split(part, ",")
	var array []byte
	for _, subPart := range subParts {
		val, err := strconv.Atoi(subPart)
		if err != nil {
			return nil, fmt.Errorf("error converting string to byte: %v", err)
		}
		array = append(array, byte(val))
	}
	return array, nil
}

// GetAllZipFiles returns a list of all zip files in the specified directory and its subdirectories.
func GetAllZipFiles(rootDir string) ([]string, error) {
	var zipFiles []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".zip") {
			zipFiles = append(zipFiles, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return zipFiles, nil
}

// ExtractZip extracts a zip file and returns the list of extracted files.
func ExtractZip(src string) ([]string, error) {
	var extractedFiles []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return nil, fmt.Errorf("error opening zip file: %v", err)
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := f.Name

		if f.FileInfo().IsDir() {
			// Create directory
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return nil, fmt.Errorf("error creating directory: %v", err)
			}
			continue
		}

		// Create file
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return nil, fmt.Errorf("error creating directory: %v", err)
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return nil, fmt.Errorf("error creating file: %v", err)
		}

		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("error opening file in zip: %v", err)
		}

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return nil, fmt.Errorf("error copying file content: %v", err)
		}

		outFile.Close()
		rc.Close()

		extractedFiles = append(extractedFiles, fpath)
	}

	return extractedFiles, nil
}

func processTextFile(fileName string, config *Config) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	var startArray, stopArray []byte
	var remainingLines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if startArray == nil && stopArray == nil {
			startArray, stopArray, err = stringToByteArray(line)
			if err != nil {
				fmt.Printf("Error processing line: %v\n", err)
				return
			}
		} else {
			remainingLines = append(remainingLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	if startArray == nil || stopArray == nil {
		fmt.Printf("Invalid range in file: %v\n", fileName)
		return
	}

	fmt.Printf("Processing: %v - %v\n", startArray, stopArray)

	program := NewProgram()

	var wg sync.WaitGroup
	numWorkers := config.NumWorkers
	wg.Add(numWorkers)

	done := make(chan struct{})
	var once sync.Once

	for i := 0; i < numWorkers; i++ {
		go processTasks(program.tasks, &wg, config.ExistingHash, done, &once)
	}

	program.GenerateAllByteArrays(len(startArray), startArray, stopArray)
	wg.Wait()

	if len(remainingLines) > 0 {
		err = os.WriteFile(fileName, []byte(strings.Join(remainingLines, "\n")), 0644)
		if err != nil {
			fmt.Printf("Error writing remaining lines to file: %v\n", err)
		}
	} else {
		err = os.Remove(fileName)
		if err != nil {
			fmt.Printf("Error deleting file: %v\n", err)
		}
	}
}

// GetAllPermutationFiles returns a list of all permutation text files in the specified directory and its subdirectories.
func GetAllPermutationFiles(rootDir string) ([]string, error) {
	var permFiles []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.Contains(info.Name(), "permutation") && strings.HasSuffix(info.Name(), ".txt") {
			permFiles = append(permFiles, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return permFiles, nil
}

func main() {
	config, err := loadConfig("appsettings.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	// Process existing permutation text files first
	permFiles, err := GetAllPermutationFiles(".")
	if err != nil {
		fmt.Printf("Error getting permutation files: %v\n", err)
		return
	}

	for _, permFile := range permFiles {
		fmt.Printf("Processing permutation file: %v\n", permFile)
		processTextFile(permFile, config)
	}

	// Process zip files
	zipFiles, err := GetAllZipFiles(".")
	if err != nil {
		fmt.Printf("Error getting zip files: %v\n", err)
		return
	}

	for _, zipFile := range zipFiles {
		fmt.Printf("Processing zip file: %v\n", zipFile)
		extractedFiles, err := ExtractZip(zipFile)
		if err != nil {
			fmt.Printf("Error extracting zip file: %v\n", err)
			continue
		}

		for _, extractedFile := range extractedFiles {
			if strings.HasSuffix(extractedFile, ".txt") {
				fmt.Printf("Processing permutation file: %v\n", extractedFile)
				processTextFile(extractedFile, config)
			}
		}

		// Remove the zip file after processing
		err = os.Remove(zipFile)
		if err != nil {
			fmt.Printf("Error removing zip file: %v\n", err)
		} else {
			fmt.Printf("Removed zip file: %v\n", zipFile)
		}
	}
}
