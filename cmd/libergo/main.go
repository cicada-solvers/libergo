package main

import (
	"flag"
	"fmt"
	"liberdatabase"
	"os"
	"os/exec"
	"runtime"

	"config"
)

func main() {
	initFlag := flag.Bool("init", false, "Initialize the default configuration")
	listFlag := flag.Bool("list", false, "List the current configuration")
	workersFlag := flag.Int("workers", 0, "Set the number of workers")
	initDBServerFlag := flag.Bool("initdbserver", false, "Initialize the podman database server")
	flag.Parse()

	if !*initFlag && !*listFlag && *workersFlag <= 0 && !*initDBServerFlag {
		fmt.Println("Usage:")
		flag.PrintDefaults()
	}

	if *initFlag {
		err := config.CreateDefaultConfig()
		if err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Error creating default config: %v\n", err)
			if err != nil {
				return
			}
			os.Exit(1)
		}
		fmt.Println("Default configuration created successfully.")
		os.Exit(0)
	}

	if *listFlag {
		cfg, err := config.LoadConfig()
		if err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			if err != nil {
				return
			}
			os.Exit(1)
		}
		fmt.Println("Current configuration:")
		fmt.Printf("  NumWorkers:              %d\n", cfg.NumWorkers)
		fmt.Printf("  ExistingHash:            %s\n", cfg.ExistingHash)
		fmt.Printf("  AdminConnectionString:   %s\n", cfg.AdminConnectionString)
		fmt.Printf("  GeneralConnectionString: %s\n", cfg.GeneralConnectionString)
		fmt.Printf("  MaxPermutationsPerLine:  %d\n", cfg.MaxPermutationsPerLine)
		fmt.Printf("  MaxRangesPerSegment:     %d\n", cfg.MaxRangesPerSegment)
		fmt.Printf("  MaxSegmentsPerPackage:   %d\n", cfg.MaxSegmentsPerPackage)
		os.Exit(0)
	}

	if *workersFlag != -1 {
		err := config.UpdateConfig("NumWorkers", *workersFlag)
		if err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Error updating NumWorkers: %v\n", err)
			if err != nil {
				return
			}
			os.Exit(1)
		}
		fmt.Println("NumWorkers updated successfully.")
		os.Exit(0)
	}

	if *initDBServerFlag {
		executeScript()
		os.Exit(0)
	}

	fmt.Println("This program is used to init or update the configuration for the toolset.")
}

func executeScript() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin", "linux":
		cmd = exec.Command("sh", ".libergo/create_podman_db.sh")
	default:
		fmt.Println("Unsupported operating system")
		return
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error executing script: %v\n", err)
	}

	_, dbError := liberdatabase.InitPermutationDatabase()
	if dbError != nil {
		return
	} else {
		fmt.Println("Database initialized successfully.")
	}
}
