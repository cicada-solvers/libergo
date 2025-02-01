package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config represents the configuration
type Config struct {
	NumWorkers   int    `json:"num_workers"`
	ExistingHash string `json:"existing_hash"`
}

// loadConfig loads the configuration from a file
func loadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("error closing config file: %v", err)
		}
	}(file)

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("error decoding config file: %v", err)
	}

	return &config, nil
}
