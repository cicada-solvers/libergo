package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// AppConfig represents the application configuration.
type AppConfig struct {
	NumWorkers              int    `json:"num_workers"`
	ExistingHash            string `json:"existing_hash"`
	AdminConnectionString   string `json:"admin_connection_string"`
	GeneralConnectionString string `json:"general_connection_string"`
	MaxPermutationsPerLine  int64  `json:"max_permutations_per_line"`
	MaxRangesPerSegment     int64  `json:"max_ranges_per_segment"`
	MaxSegmentsPerPackage   int64  `json:"max_segments_per_package"`
}

// getConfigFilePath returns the path to the configuration file.
func getConfigFilePath() (string, error) {
	var configDir string
	if runtime.GOOS == "windows" {
		configDir = filepath.Join(os.TempDir(), ".libergo")
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(homeDir, ".libergo")
	}
	return filepath.Join(configDir, "appsettings.json"), nil
}

// LoadConfig reads the configuration file and returns the AppConfig struct.
func LoadConfig() (*AppConfig, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	var config AppConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// CreateDefaultConfig creates a default configuration file.
func CreateDefaultConfig() error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	// Remove the existing configuration file if it exists
	if _, err := os.Stat(configFilePath); err == nil {
		err = os.Remove(configFilePath)
		if err != nil {
			return err
		}
	}

	configDir := filepath.Dir(configFilePath)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			return err
		}
	}

	defaultConfig := AppConfig{
		NumWorkers:              10,
		ExistingHash:            "36367763ab73783c7af284446c59466b4cd653239a311cb7116d4618dee09a8425893dc7500b464fdaf1672d7bef5e891c6e2274568926a49fb4f45132c2a8b4",
		AdminConnectionString:   "postgres://postgres:lppasswd@localhost:5432/postgres",
		GeneralConnectionString: "postgres://postgres:lppasswd@localhost:5432/libergodb",
		MaxPermutationsPerLine:  500000000,
		MaxRangesPerSegment:     250,
		MaxSegmentsPerPackage:   250,
	}

	data, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(configFilePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// UpdateConfig updates a specific field in the configuration file and saves it.
func UpdateConfig(key string, value interface{}) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	var config AppConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	switch key {
	case "NumWorkers":
		config.NumWorkers = value.(int)
	case "ExistingHash":
		config.ExistingHash = value.(string)
	case "AdminConnectionString":
		config.AdminConnectionString = value.(string)
	case "GeneralConnectionString":
		config.GeneralConnectionString = value.(string)
	case "MaxPermutationsPerLine":
		config.MaxPermutationsPerLine = value.(int64)
	case "MaxRangesPerSegment":
		config.MaxRangesPerSegment = value.(int64)
	case "MaxSegmentsPerPackage":
		config.MaxSegmentsPerPackage = value.(int64)
	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}

	data, err = json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(configFilePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
