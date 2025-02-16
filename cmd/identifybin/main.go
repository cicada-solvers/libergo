package main

import (
	"filetypeinterrogator"
	"flag"
	"fmt"
	"liberdatabase"
	"os"
)

func main() {
	// Define the filename flag
	filename := flag.String("filename", "", "The name of the file to identify")
	flag.Parse()

	// Check if the filename is provided
	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Get all file type definitions
	fileTypeInfos, err := liberdatabase.GetAllFileTypeInfo()
	if err != nil {
		fmt.Println("Error retrieving file type definitions:", err)
		os.Exit(1)
	}

	// Initialize the FileTypeInterrogator
	interrogator := filetypeinterrogator.NewFileTypeInterrogator(fileTypeInfos)

	// Open the file
	file, err := os.Open(*filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(file)

	// Detect the file type from the file stream
	fileTypeInfo, err := interrogator.DetectTypeFromStream(file)
	if err != nil {
		fmt.Println("Error detecting file type:", err)
		os.Exit(1)
	}

	// Print the detected file type information
	if fileTypeInfo != nil {
		fmt.Printf("Detected file type: %s\n", fileTypeInfo.FileType)
		fmt.Printf("MIME type: %s\n", fileTypeInfo.MimeType)
	} else {
		fmt.Println("File type could not be detected")
	}
}
