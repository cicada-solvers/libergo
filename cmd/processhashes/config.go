package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	NumWorkers   int    `json:"num_workers"`
	ExistingHash string `json:"existing_hash"`
}

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
