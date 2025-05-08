package main

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/jzelinskie/whirlpool"
	"golang.org/x/crypto/blake2b"
	"hashinglib"
	"net"
	"os"
	"strings"
	"sync"
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

	var wg sync.WaitGroup

	wg.Add(4) // Add the number of goroutines to the WaitGroup

	go func() {
		defer wg.Done()
		checkIPs(class, "ips", start, end, false, false, 1)
	}()

	go func() {
		defer wg.Done()
		checkIPs(class, "ipswport", start, end, true, false, 2)
	}()

	go func() {
		defer wg.Done()
		checkIPs(class, "ipswscheme", start, end, false, true, 3)
	}()

	go func() {
		defer wg.Done()
		checkIPs(class, "ipswportwscheme", start, end, true, true, 4)
	}()

	wg.Wait() // Wait for all goroutines to finish

	fmt.Println("Files written successfully.")
}

func checkIPs(class, portion string, start, end int64, includePorts, includeSchemes bool, index int) {
	totalIPs := end - start + 1

	threadID := index // Use the index to identify the thread
	for ip := start; ip <= end; ip++ {
		// Calculate and display progress
		progress := float64(ip-start+1) / float64(totalIPs) * 100
		fmt.Printf("\rThread %d - Processing %s (%s): %.2f%% complete", threadID, class, portion, progress)

		if includePorts {
			for port := 1; port <= 65535; port++ {
				line := fmt.Sprintf("%s:%d\n", intToIP(ip).String(), port)
				if includeSchemes {
					for _, scheme := range getSchemes() {
						line = fmt.Sprintf("%s://%s:%d\n", scheme, intToIP(ip).String(), port)
						checkLine(line)
					}
				} else {
					checkLine(line)
				}
			}
		} else {
			line := intToIP(ip).String() + "\n"
			if includeSchemes {
				for _, scheme := range getSchemes() {
					line = fmt.Sprintf("%s://%s\n", scheme, intToIP(ip).String())
					checkLine(line)
				}
			} else {
				checkLine(line)
			}
		}
	}
	fmt.Printf("Thread %d - Processing %s (%s): 100.00%% complete\n", threadID, class, portion) // Ensure 100% is printed at the end
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
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	}
}
