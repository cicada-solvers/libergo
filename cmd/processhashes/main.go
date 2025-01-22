package main

import (
	"bufio"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/sha3"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

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
				fmt.Printf("Stopped on: %v\n", currentArray)
				fmt.Printf("Stop Array Was: %v\n", stopArray)
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
			fmt.Printf("Hashing: %v\n", task)
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

func stringToByteArray(s string) []byte {
	parts := strings.Split(s, ",")
	array := make([]byte, len(parts))
	for i, part := range parts {
		val, err := strconv.Atoi(part)
		if err != nil {
			fmt.Printf("Error converting string to byte: %v\n", err)
			return nil
		}
		array[i] = byte(val)
	}
	return array
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./processhashes <filename>")
		return
	}

	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	var startArray, stopArray []byte
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		startArray = stringToByteArray(scanner.Text())
	}
	if scanner.Scan() {
		stopArray = stringToByteArray(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	fmt.Printf("Processing: %v - %v\n", startArray, stopArray)

	program := NewProgram()

	// Read the existing hash from file
	existingHashBytes, err := os.ReadFile("existinghash.txt")
	if err != nil {
		fmt.Printf("Error reading existing hash: %v\n", err)
		return
	}
	existingHash := string(existingHashBytes)

	var wg sync.WaitGroup
	numWorkers := 10
	wg.Add(numWorkers)

	done := make(chan struct{})
	var once sync.Once

	for i := 0; i < numWorkers; i++ {
		go processTasks(program.tasks, &wg, existingHash, done, &once)
	}

	program.GenerateAllByteArrays(len(startArray), startArray, stopArray)
	wg.Wait()
}
