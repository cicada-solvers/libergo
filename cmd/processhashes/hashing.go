package main

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jzelinskie/whirlpool"
	"golang.org/x/crypto/blake2b"
	"hashinglib"
)

// processTasks processes the tasks
func processTasks(tasks chan []byte, wg *sync.WaitGroup, existingHash string, done chan struct{}, once *sync.Once, totalPermutations *big.Int, mu *sync.Mutex) {
	defer wg.Done()

	file, err := os.OpenFile("found_hashes.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	buffer := make([]byte, 0, 4096)
	hashCount := 0
	taskLen := 0
	processedPermutations := big.NewInt(0)
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	go func() {
		colors := []string{"\033[31m", "\033[32m", "\033[33m", "\033[34m", "\033[35m", "\033[36m", "\033[37m", "\033[90m", "\033[91m", "\033[92m"}
		colorIndex := 0
		for range ticker.C {
			mu.Lock()
			remainingPermutations := new(big.Int).Sub(totalPermutations, processedPermutations)
			mu.Unlock()
			fmt.Printf("%sHashes per minute: %d, Array size: %d, Remaining hashes: %s\033[0m\n", colors[colorIndex], hashCount, taskLen, remainingPermutations.String())
			hashCount = 0
			colorIndex = (colorIndex + 1) % len(colors)
		}
	}()

	for {
		select {
		case task, ok := <-tasks:
			if !ok {
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
				mu.Lock()
				processedPermutations.Add(processedPermutations, big.NewInt(1))
				mu.Unlock()

				//fmt.Printf("Debug: Hash Name: %s, Byte Array: %s\n, Hash: %s", hashName, hex.EncodeToString(task), hash)

				if hash == existingHash {
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
					once.Do(func() { close(done) })
				}
			}
		case <-done:
			return
		}
	}
}

// generateHashes generates hashes for a given byte array
func generateHashes(data []byte) map[string]string {
	hashes := make(map[string]string)

	sha512Hash := sha512.Sum512(data)
	hashes["SHA-512"] = hex.EncodeToString(sha512Hash[:])

	whirlpoolHash := whirlpool.New()
	whirlpoolHash.Write(data)
	whirlHash := whirlpoolHash.Sum(nil)
	hashes["Whirlpool-512"] = hex.EncodeToString(whirlHash[:])

	blake2bHash := blake2b.Sum512(data)
	hashes["Blake2b-512"] = hex.EncodeToString(blake2bHash[:])

	blake512Hash := hashinglib.Blake512Hash(data)
	hashes["Blake-512"] = hex.EncodeToString(blake512Hash)

	return hashes
}

// convertToByteArray converts a string to a byte array
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
