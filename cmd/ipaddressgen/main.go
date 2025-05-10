package main

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/jzelinskie/whirlpool"
	"golang.org/x/crypto/blake2b"
	"hashinglib"
	"math/big"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

func main() {
	ranges := map[string][2]string{
		"A": {"1.0.0.0", "127.0.0.0"},
		"B": {"128.0.0.0", "191.255.255.255"},
		"C": {"192.0.0.0", "223.255.255.255"},
		"D": {"224.0.0.0", "239.255.255.255"},
		"E": {"240.0.0.0", "255.255.255.255"},
	}

	var class string
	fmt.Println("Enter the IPv4 class range (A, B, C, D, E):")
	_, scanErr := fmt.Scanln(&class)
	if scanErr != nil {
		return
	}
	class = strings.ToUpper(class)

	if _, exists := ranges[class]; !exists {
		fmt.Println("Invalid class. Please enter A, B, C, D, or E.")
		return
	}

	startIP := net.ParseIP(ranges[class][0]).To4()
	endIP := net.ParseIP(ranges[class][1]).To4()

	start := ipToInt(startIP)
	end := ipToInt(endIP)

	checkIPs(start, end)
	fmt.Println("Files written successfully.")
}

func checkIPs(start, end int64) {
	var processedCounter = big.NewInt(0)
	var rateCounter = big.NewInt(0)
	one := big.NewInt(1)
	totalIps := end - start + 1
	fmt.Printf("Processing %d IPs...\n", totalIps)

	// We are going to put timer to see how many we have processed.
	processedTicker := time.NewTicker(time.Minute)
	defer processedTicker.Stop()

	go func() {
		for range processedTicker.C {
			fmt.Printf("\rRate: %s/min - Processed %s items                                                    ",
				rateCounter.String(), processedCounter.String())
			rateCounter.SetInt64(int64(0))
		}
	}()

	schemes := getSchemes()

	// Create a channel to send lines
	linesChan := make(chan string, 1024)
	addressChan := make(chan string, 1024)

	go func() {
		for line := range linesChan {
			rateCounter.Add(rateCounter, one)
			checkLine(line)
		}
	}()

	go func() {
		for ip := start; ip <= end; ip++ {
			ipString := intToIP(ip).String()
			addressChan <- ipString
		}
	}()

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Get the number of processors
	numProcessors := runtime.NumCPU() // Use double the number of processors for more concurrency

	// Start worker goroutines
	for i := 0; i < numProcessors; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for address := range addressChan {
				processedCounter.Add(processedCounter, one)
				line := fmt.Sprintf("%s", address)
				linesChan <- line

				for _, scheme := range schemes {
					line = fmt.Sprintf("%s://%s", scheme, address)
					linesChan <- line

					line = fmt.Sprintf("%s://%s/", scheme, address)
					linesChan <- line
				}

				for port := 1; port <= 65535; port++ {
					line = fmt.Sprintf("%s:%d", address, port)
					linesChan <- line

					for _, scheme := range schemes {
						line = fmt.Sprintf("%s://%s:%d", scheme, address, port)
						linesChan <- line

						line = fmt.Sprintf("%s://%s:%d/", scheme, address, port)
						linesChan <- line
					}
				}
			}
		}()
	}

	// Wait for all workers to finish
	wg.Wait()

	fmt.Printf("Processing: 100.00%% complete\n") // Ensure 100% is printed at the end
}

func checkLine(line string) {
	existingHash := "36367763ab73783c7af284446c59466b4cd653239a311cb7116d4618dee09a8425893dc7500b464fdaf1672d7bef5e891c6e2274568926a49fb4f45132c2a8b4"

	// Convert the string to a byte array
	byteArray := []byte(line)

	hashes := generateHashes(byteArray)

	for hashName, hash := range hashes {
		if hash == existingHash {
			output := fmt.Sprintf("Found matching hash:%s - %s:%s\n", line, hashName, hash)
			fmt.Printf(output)
			writeToOutputFile("output.txt", []byte(output))
		}
	}
}

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

func writeToOutputFile(filename string, data []byte) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

	if _, writeErr := file.Write(data); writeErr != nil {
		fmt.Printf("Error writing to file: %v\n", writeErr)
	}
}
