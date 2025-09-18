package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	inPath := flag.String("in", "", "Path to the input file to read bytes from")
	outPath := flag.String("out", "", "Path to the output CSV file (defaults to <input>.bytes.txt)")
	flag.Usage = usage
	flag.Parse()

	if *inPath == "" {
		usage()
		os.Exit(2)
	}

	info, err := os.Stat(*inPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot access input file: %v\n", err)
		os.Exit(1)
	}
	if !info.Mode().IsRegular() {
		fmt.Fprintf(os.Stderr, "Error: input is not a regular file\n")
		os.Exit(1)
	}

	out := *outPath
	if out == "" {
		out = outputPathFor(*inPath)
	}

	if err := writeBytesCSV(*inPath, out); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Wrote byte CSV to: %s\n", out)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s -in <input-file> [-out <output-file>]\n", filepath.Base(os.Args[0]))
	flag.PrintDefaults()
}

func outputPathFor(inputPath string) string {
	dir := filepath.Dir(inputPath)
	base := filepath.Base(inputPath)
	return filepath.Join(dir, base+".bytes.txt")
}

func writeBytesCSV(inputPath, outputPath string) error {
	in, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("open input: %w", err)
	}
	defer in.Close()

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer func() {
		if cerr := out.Close(); err == nil && cerr != nil {
			err = cerr
		}
	}()

	reader := bufio.NewReader(in)
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	const bufSize = 64 * 1024
	buf := make([]byte, bufSize)

	first := true
	for {
		n, rerr := reader.Read(buf)
		if n > 0 {
			tmp := make([]byte, 0, n*4)
			for i := 0; i < n; i++ {
				if first {
					first = false
				} else {
					tmp = append(tmp, ',')
				}
				tmp = strconv.AppendInt(tmp, int64(buf[i]), 10)
			}
			if _, werr := writer.Write(tmp); werr != nil {
				return fmt.Errorf("write output: %w", werr)
			}
		}
		if errors.Is(rerr, io.EOF) {
			break
		}
		if rerr != nil {
			return fmt.Errorf("read input: %w", rerr)
		}
	}

	if _, werr := writer.WriteString("\n"); werr != nil {
		return fmt.Errorf("finalize output: %w", werr)
	}

	return nil
}
