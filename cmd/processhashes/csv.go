package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type PackFile struct {
	Id                   string
	StartArray           []byte
	EndArray             []byte
	PackageName          string
	PermName             string
	ReportedToApi        bool
	Processed            bool
	ArrayLength          int
	NumberOfPermutations int64
}

// deleteLineById removes a line from the CSV file based on the given ID
func deleteLineById(filePath, id string) error {
	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Read all records
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV file: %v", err)
	}

	// Filter out the record with the matching ID
	var updatedRecords [][]string
	for _, record := range records {
		if len(record) > 0 && record[0] == id {
			continue // Skip the record with the matching ID
		}
		updatedRecords = append(updatedRecords, record)
	}

	// Write the updated records back to the file
	file.Close() // Close the file before overwriting
	file, err = os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	err = writer.WriteAll(updatedRecords)
	if err != nil {
		return fmt.Errorf("error writing to CSV file: %v", err)
	}

	return nil
}

// readCSVToPackFiles reads a CSV file and parses its rows into an array of PackFile structs
func readCSVToPackFiles(filePath string) ([]PackFile, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			fmt.Printf("Error closing file: %v\n", closeErr)
		}
	}(file)

	reader := csv.NewReader(file)

	// Skip the first row (header)
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading header row: %v", err)
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV file: %v", err)
	}

	var packFiles []PackFile
	for _, record := range records {
		if len(record) < 8 {
			continue // Skip rows with insufficient columns
		}

		arrayLength, _ := strconv.Atoi(record[7])
		numPermutations, _ := strconv.ParseInt(record[8], 10, 64)
		var startArray []byte
		var endArray []byte
		for _, b := range strings.Split(record[1], ",") {
			value, _ := strconv.Atoi(b)
			startArray = append(startArray, byte(value))
		}
		for _, b := range strings.Split(record[2], ",") {
			value, _ := strconv.Atoi(b)
			endArray = append(endArray, byte(value))
		}

		packFile := PackFile{
			Id:                   record[0],
			StartArray:           startArray,
			EndArray:             endArray,
			PackageName:          record[3],
			PermName:             record[4],
			ReportedToApi:        record[5] == "true",
			Processed:            record[6] == "true",
			ArrayLength:          arrayLength,
			NumberOfPermutations: numPermutations,
		}
		packFiles = append(packFiles, packFile)
	}

	return packFiles, nil
}

// downloadPack downloads a .7z file corresponding to the given pack number and extracts it
func downloadAndExtractPack(packNumber int) (string, error) {
	url := fmt.Sprintf("https://cmbsolver.com/downloads/packs/PACK_%d.7z", packNumber)
	fileName := fmt.Sprintf("PACK_%d.7z", packNumber)
	outputDir := fmt.Sprintf("PACK_%d", packNumber)

	downloadErr := downloadLargeFile(url, fileName)
	if downloadErr != nil {
		return "", fmt.Errorf("error downloading file: %v", downloadErr)
	}

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("error creating output directory: %v", err)
	}

	// Run the 7z command to extract the file
	cmd := exec.Command("7z", "x", fileName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error decompressing file: %v", err)
	}

	fmt.Printf("Pack %d downloaded and decompressed successfully to %s\n", packNumber, outputDir)
	return outputDir, nil
}

func downloadLargeFile(url, fileName string) error {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Check if the file already exists and get its size
	var startByte int64
	if fileInfo, err := os.Stat(fileName); err == nil {
		startByte = fileInfo.Size()
	}

	// Create request with Range header for resuming
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	if startByte > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", startByte))
	}

	// Set a User-Agent header to mimic a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	// Setting the header to avoid 503 errors
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}(resp.Body)

	// Check for valid response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Open file for appending
	out, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}(out)

	// Write the response body to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}

// listCSVFiles returns a slice of .csv file names in the specified directory
func listCSVFiles(dir string) ([]string, error) {
	var csvFiles []string

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".csv" {
			csvFiles = append(csvFiles, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return csvFiles, nil
}
