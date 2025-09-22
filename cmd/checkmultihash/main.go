package main

import (
	"blake"
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"fnv5120"
	"fnv5121"
	"groestl"
	"io"
	jh2 "jh"
	"lsh"
	"md6"
	"os"
	"runtime"
	skein2 "skein"
	"strconv"
	"strings"
	"sync"
)

// main reads a file line by line and processes each line concurrently to match against a predefined hash value.
func main() {
	existingHash := "36367763ab73783c7af284446c59466b4cd653239a311cb7116d4618dee09a8425893dc7500b464fdaf1672d7bef5e891c6e2274568926a49fb4f45132c2a8b4"

	filename := flag.String("filename", "", "File to analyze")
	isByteFile := flag.Bool("bytefile", false, "File contains hashes")

	// Parse the flags
	flag.Parse()

	// Validate that text was provided
	if *filename == "" {
		fmt.Println("Error: No file provided for analysis")
		flag.Usage()
		os.Exit(1)
	}

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
				processLine(line, existingHash, *isByteFile)
			}
		}()
	}

	// Read lines from the file and send them to the channel
	lines, fileError := ReadFileLinesBytewise(*filename)
	if fileError != nil {
		fmt.Printf("Error reading file: %v\n", fileError)
		os.Exit(1)
	}

	for _, line := range lines {
		linesChan <- line
	}

	close(linesChan)

	// Wait for all workers to finish
	wg.Wait()
}

// ReadFileLinesBytewise reads the file at path one byte at a time and returns its lines.
// It treats '\n' as line separator and strips a preceding '\r' (handling CRLF).
func ReadFileLinesBytewise(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		closeError := f.Close()
		if closeError != nil {
			fmt.Printf("Error: %v\n", closeError)
		}
	}(f)

	r := bufio.NewReader(f)
	var (
		lines []string
		cur   []byte
	)

	for {
		b, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				if len(cur) > 0 {
					lines = append(lines, string(cur))
				}
				break
			}
			return nil, err
		}

		if b == '\n' {
			// Trim trailing '\r' if present (CRLF)
			if n := len(cur); n > 0 && cur[n-1] == '\r' {
				cur = cur[:n-1]
			}
			lines = append(lines, string(cur))
			cur = cur[:0]
			continue
		}
		cur = append(cur, b)
	}

	return lines, nil
}

// processLine processes an input string by generating its hash values and compares them with an existing hash.
// If a match is found, it outputs the matching hash type and value along with the input string.
func processLine(inputString, existingHash string, isByteFile bool) {
	// Convert the string to a byte array
	var byteArray []byte

	if isByteFile {
		convertedString := strings.ReplaceAll(inputString, " ", "")
		byteArray = convertByteCsvToByte(convertedString)
	} else {
		byteArray = convertByteCsvToByte(inputString)
	}

	hashes := generateHashes(byteArray)

	for hashName, hash := range hashes {
		if hash == existingHash {
			fmt.Printf("Found matching hash:%s - %s:%s\n", inputString, hashName, hash)
		}
	}
}

func convertByteCsvToByte(inputString string) []byte {
	byteArray := make([]byte, 0)
	byteStringArray := strings.Split(inputString, ",")
	for _, byteStr := range byteStringArray {
		s := strings.TrimSpace(byteStr)
		if s == "" {
			continue
		}
		v, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			// handle error as needed; here we skip invalid entries
			continue
		}
		byteArray = append(byteArray, byte(v))
	}
	return byteArray
}

// generateHashes computes hash values for the input data using SHA-512, Whirlpool-512, Blake2b-512, and Blake-512 algorithms.
// It returns a map where keys are hash algorithm names and values are corresponding hash strings.
func generateHashes(data []byte) map[string]string {
	hashes := make(map[string]string)

	blake512Hash := blake.Blake512Hash(data)
	hashes["Blake-512"] = hex.EncodeToString(blake512Hash)

	jh := jh2.NewJH512()
	jh.Write(data)
	jhHash := jh.Sum(nil)
	hashes["JH-512"] = hex.EncodeToString(jhHash)

	skein := skein2.NewSkein512()
	skein.Write(data)
	skeinHash := skein.Sum(nil)
	hashes["Skein-512"] = hex.EncodeToString(skeinHash)

	grostl := groestl.NewGroestl512()
	grostl.Write(data)
	grostlHash := grostl.Sum(nil)
	hashes["Groestl-512"] = hex.EncodeToString(grostlHash)

	fnv5120hash := fnv5120.Hash512(data)
	hashes["FNV5120-512"] = hex.EncodeToString(fnv5120hash)

	fnv5120ahash := fnv5120.Hash512a(data)
	hashes["FNV5120A-512"] = hex.EncodeToString(fnv5120ahash)

	fnv5121hash := fnv5121.Hash(data)
	hashes["FNV5121-512"] = hex.EncodeToString(fnv5121hash)

	fnv5121ahash := fnv5121.HashA(data)
	hashes["FNV5121A-512"] = hex.EncodeToString(fnv5121ahash)

	md6hash := md6.Sum512(data)
	hashes["MD6-512"] = hex.EncodeToString(md6hash)

	lshHash := lsh.Sum512(data)
	hashes["LSH-512"] = hex.EncodeToString(lshHash)

	return hashes
}
