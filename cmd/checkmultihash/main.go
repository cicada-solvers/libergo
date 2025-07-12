package main

import (
	"bufio"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/jzelinskie/whirlpool"
	"golang.org/x/crypto/blake2b"
	"hashinglib"
	"os"
	"runtime"
	"sync"
)

// main reads a file line by line and processes each line concurrently to match against a predefined hash value.
func main() {
	existingHash := "36367763ab73783c7af284446c59466b4cd653239a311cb7116d4618dee09a8425893dc7500b464fdaf1672d7bef5e891c6e2274568926a49fb4f45132c2a8b4"

	// Check if the user provided an argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: checkmultihash <file-name>")
		return
	}

	// Get the file name from the arguments
	fileName := os.Args[1]

	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			fmt.Printf("Error closing file: %v\n", closeErr)
		}
	}(file)

	// Create a channel to send lines
	linesChan := make(chan string)

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Get the number of processors
	numProcessors := runtime.NumCPU()

	// Start worker goroutines
	for i := 0; i < numProcessors; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range linesChan {
				processLine(line, existingHash)
			}
		}()
	}

	// Read lines from the file and send them to the channel
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		linesChan <- scanner.Text()
	}
	close(linesChan)

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// Wait for all workers to finish
	wg.Wait()
}

// processLine processes an input string by generating its hash values and compares them with an existing hash.
// If a match is found, it outputs the matching hash type and value along with the input string.
func processLine(inputString, existingHash string) {
	// Convert the string to a byte array
	byteArray := []byte(inputString)

	hashes := generateHashes(byteArray)

	for hashName, hash := range hashes {
		if hash == existingHash {
			fmt.Printf("Found matching hash:%s - %s:%s\n", inputString, hashName, hash)
		}
	}
}

// generateHashes computes hash values for the input data using SHA-512, Whirlpool-512, Blake2b-512, and Blake-512 algorithms.
// It returns a map where keys are hash algorithm names and values are corresponding hash strings.
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
