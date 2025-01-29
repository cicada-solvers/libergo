package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config represents the configuration for the application
type Config struct {
	NumWorkers             int   `json:"num_workers"`
	MaxPermutationsPerLine int64 `json:"max_permutations_per_line"`
	MaxRangesPerSegment    int64 `json:"max_ranges_per_segment"`
	MaxSegmentsPerPackage  int64 `json:"max_segments_per_package"`
}

// loadConfig loads the configuration from the specified file
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
