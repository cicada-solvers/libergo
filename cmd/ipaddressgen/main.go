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

var processedCounter = big.NewInt(0)
var rateCounter = big.NewInt(0)

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

	// Define our processed file name
	processedFileName := "processed.txt"

	// Try to read the last processed IP from file
	lastProcessedIP, err := readProcessedIPFromFile(processedFileName)
	if err == nil {
		// If we have a valid IP and it's in range, start from the next IP
		if lastProcessedIP >= start && lastProcessedIP < end {
			fmt.Printf("Resuming from IP: %s\n", intToIP(lastProcessedIP).String())
			start = lastProcessedIP + 1
		} else {
			fmt.Println("Saved IP is outside the selected range, starting from the beginning")
		}
	} else {
		fmt.Println("Starting from the beginning of the range")
	}

	// We are going to put timer to see how many we have processed.
	processedTicker := time.NewTicker(time.Minute)
	defer processedTicker.Stop()

	// Create a variable to track the current IP position
	currentIP := start

	go func() {
		for range processedTicker.C {
			fmt.Printf("\rRate: %s/min - Processed %s items                                                    ",
				rateCounter.String(), processedCounter.String())

			// Write counter and current IP position to file
			if err := writeCounterToFile(processedCounter, currentIP, processedFileName); err != nil {
				fmt.Printf("\nError writing to file: %v\n", err)
			}

			rateCounter.SetInt64(int64(0))
		}
	}()

	checkIPs(start, end, &currentIP)

	removeErr := os.Remove(processedFileName)
	if removeErr != nil {
		fmt.Printf("Error removing processed file: %v\n", removeErr)
	}
	fmt.Println("Processed range successfully.")
}

func checkIPs(start, end int64, currentIP *int64) {
	totalIps := end - start + 1
	fmt.Printf("Processing %d IPs...\n", totalIps)

	schemes := getSchemes()

	// Create a channel to send lines
	addressChan := make(chan string, 1024)

	go func() {
		for ip := start; ip <= end; ip++ {
			processedCounter.Add(processedCounter, big.NewInt(1))
			*currentIP = ip
			ipString := intToIP(ip).String()
			addressChan <- ipString
		}
	}()

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Get the number of processors
	numProcessors := runtime.NumCPU() * 2 // Use double the number of processors for more concurrency

	// Start worker goroutines
	for i := 0; i < numProcessors; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for address := range addressChan {
				line := fmt.Sprintf("%s", address)
				checkLine(line)

				for _, scheme := range schemes {
					line = fmt.Sprintf("%s://%s", scheme, address)
					checkLine(line)

					line = fmt.Sprintf("%s://%s/", scheme, address)
					checkLine(line)
				}

				//for port := 1; port <= 65535; port++ {
				//	line = fmt.Sprintf("%s:%d", address, port)
				//	checkLine(line)

				//for _, scheme := range schemes {
				//	line = fmt.Sprintf("%s://%s:%d", scheme, address, port)
				//	checkLine(line)
				//
				//	line = fmt.Sprintf("%s://%s:%d/", scheme, address, port)
				//	checkLine(line)
				//}
				//}
			}
		}()
	}

	// Wait for all workers to finish
	wg.Wait()

	fmt.Printf("Processing: 100.00%% complete\n") // Ensure 100% is printed at the end
}

func checkLine(line string) {
	existingHash := "36367763ab73783c7af284446c59466b4cd653239a311cb7116d4618dee09a8425893dc7500b464fdaf1672d7bef5e891c6e2274568926a49fb4f45132c2a8b4"
	one := big.NewInt(1)

	// Convert the string to a byte array
	byteArray := []byte(line)

	hashes := generateHashes(byteArray)

	for hashName, hash := range hashes {
		if hash == existingHash {
			output := fmt.Sprintf("Found matching hash:%s - %s:%s\n", line, hashName, hash)
			fmt.Printf(output)
			writeToOutputFile("output.txt", []byte(output))
		}

		rateCounter.Add(rateCounter, one)
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

func readProcessedIPFromFile(filename string) (int64, error) {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return 0, fmt.Errorf("file does not exist")
	}

	// Read file contents
	data, err := os.ReadFile(filename)
	if err != nil {
		return 0, fmt.Errorf("error reading file: %v", err)
	}

	// Parse the line containing the processed count
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Last processed IP: ") {
			ipStr := strings.TrimPrefix(line, "Last processed IP: ")
			ipStr = strings.TrimSpace(ipStr)

			// Convert IP string to int64
			ip := net.ParseIP(ipStr).To4()
			if ip == nil {
				return 0, fmt.Errorf("invalid IP format in file: %s", ipStr)
			}

			return ipToInt(ip), nil
		}
	}

	return 0, fmt.Errorf("no valid IP found in file")
}

func writeCounterToFile(counter *big.Int, lastIP int64, filename string) error {
	ipStr := intToIP(lastIP).String()
	data := []byte(fmt.Sprintf("Processed count: %s\nLast processed IP: %s\nTimestamp: %s\n",
		counter.String(), ipStr, time.Now().Format(time.RFC3339)))

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening counter file: %v", err)
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("error writing to counter file: %v", err)
	}

	return nil
}
