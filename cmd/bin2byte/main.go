package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	inPath := flag.String("dir", "", "Path to the input files to read bytes from")
	outfileBase := "outfile"
	flag.Usage = usage
	flag.Parse()

	if *inPath == "" {
		usage()
		os.Exit(2)
	}

	filesToArray, _ := listFilesRecursive(*inPath)

	for counter, file := range filesToArray {
		fmt.Printf("Processing file: %s\n", file)
		data, _ := os.ReadFile(file)
		dataToWrite := bytesToCSV(data)
		writeBytesCSV(file, dataToWrite, fmt.Sprintf("%s_%d.txt", outfileBase, counter))
	}

	fmt.Printf("Wrote files\n")
}

func usage() {
	fmt.Printf("Usage: %s -in <input-file> [-out <output-file>]\n", filepath.Base(os.Args[0]))
	flag.PrintDefaults()
}

func listFilesRecursive(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			// Skip entries we can't access but continue walking
			_, _ = fmt.Fprintf(os.Stderr, "Skipping %s: %v\n", path, err)
			return nil
		}
		if d.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// bytesToCSV converts a byte slice to a comma-separated list of decimal numbers.
func bytesToCSV(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	// Pre-size roughly: up to 4 chars per byte plus commas
	out := make([]byte, 0, len(b)*4)
	first := true
	for _, v := range b {
		if first {
			first = false
		} else {
			out = append(out, ',')
		}
		out = strconv.AppendInt(out, int64(v), 10)
	}
	return string(out)
}

func writeBytesCSV(filename string, data string, outputPath string) {
	// Open output in append mode (create if missing)
	out, openError := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if openError != nil {
		fmt.Printf("open output: %w", openError)
	}
	defer func() {
		if closeError := out.Close(); closeError == nil && closeError != nil {
			fmt.Printf("close output: %w", closeError)
		}
	}()

	w := bufio.NewWriter(out)
	defer func(w *bufio.Writer) {
		flushErr := w.Flush()
		if flushErr != nil {
			flushErr = errors.New("flush writer: " + flushErr.Error())
		}
	}(w)

	// Write: filename,data\n as a CSV-like line
	dataToWrite := fmt.Sprintf("%s|%s\n\n", filename, data)
	if _, writeError := w.WriteString(dataToWrite); writeError != nil {
		fmt.Printf("write filename: %w", writeError)
	}
}
