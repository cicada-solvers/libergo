package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

const maxFileSize = 5 << 30 // 1GB in bytes

func main() {
	ranges := map[string][2]string{
		"A": {"1.0.0.0", "127.0.0.0"},
		"B": {"128.0.0.0", "191.255.255.255"},
		"C": {"192.0.0.0", "223.255.255.255"},
		"D": {"224.0.0.0", "239.255.255.255"},
		"E": {"240.0.0.0", "255.255.255.255"},
	}

	var class string
	fmt.Println("Enter the IPv4 class range (A, B, C, D, E):")
	_, scanErr := fmt.Scanln(&class)
	if scanErr != nil {
		return
	}
	class = strings.ToUpper(class)

	if _, exists := ranges[class]; !exists {
		fmt.Println("Invalid class. Please enter A, B, C, D, or E.")
		return
	}

	startIP := net.ParseIP(ranges[class][0]).To4()
	endIP := net.ParseIP(ranges[class][1]).To4()

	start := ipToInt(startIP)
	end := ipToInt(endIP)

	var wg sync.WaitGroup

	wg.Add(4) // Add the number of goroutines to the WaitGroup

	go func() {
		defer wg.Done()
		writeIPsToFile(class, "ips", start, end, false, false, 0)
	}()

	go func() {
		defer wg.Done()
		writeIPsToFile(class, "ipswport", start, end, true, false, 0)
	}()

	go func() {
		defer wg.Done()
		writeIPsToFile(class, "ipswscheme", start, end, false, true, 0)
	}()

	go func() {
		defer wg.Done()
		writeIPsToFile(class, "ipswportwscheme", start, end, true, true, 0)
	}()

	wg.Wait() // Wait for all goroutines to finish

	fmt.Println("Files written successfully.")
}

func writeIPsToFile(class, portion string, start, end int64, includePorts, includeSchemes bool, index int) {
	fileIndex := 1
	file, size, failed := createFile(class, portion, index, fileIndex)
	if failed {
		return
	}

	for ip := start; ip <= end; ip++ {
		if includePorts {
			for port := 1; port <= 65535; port++ {
				line := fmt.Sprintf("%s:%d\n", intToIP(ip).String(), port)
				if includeSchemes {
					for _, scheme := range getSchemes() {
						line = fmt.Sprintf("%s://%s:%d\n", scheme, intToIP(ip).String(), port)
						size = writeLineToFile(file, line, size)
						if size > maxFileSize {
							closeFile(file)
							fileIndex++
							file, size, failed = createFile(class, portion, index, fileIndex)
							if failed {
								return
							}
						}
					}
				} else {
					size = writeLineToFile(file, line, size)
					if size > maxFileSize {
						closeFile(file)
						fileIndex++
						file, size, failed = createFile(class, portion, index, fileIndex)
						if failed {
							return
						}
					}
				}
			}
		} else {
			line := intToIP(ip).String() + "\n"
			if includeSchemes {
				for _, scheme := range getSchemes() {
					line = fmt.Sprintf("%s://%s\n", scheme, intToIP(ip).String())
					size = writeLineToFile(file, line, size)
					if size > maxFileSize {
						closeFile(file)
						fileIndex++
						file, size, failed = createFile(class, portion, index, fileIndex)
						if failed {
							return
						}
					}
				}
			} else {
				size = writeLineToFile(file, line, size)
				if size > maxFileSize {
					closeFile(file)
					fileIndex++
					file, size, failed = createFile(class, portion, index, fileIndex)
					if failed {
						return
					}
				}
			}
		}
	}
	closeFile(file)
}

func writeLineToFile(file *os.File, line string, size int64) int64 {
	lineSize := int64(len(line))
	_, err := file.WriteString(line)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return size
	}
	return size + lineSize
}

func createFile(class, portion string, index, fileIndex int) (*os.File, int64, bool) {
	fileName := fmt.Sprintf("%s_%s_part_%d_%d.txt", class, portion, index, fileIndex)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return nil, 0, true
	}
	return file, 0, false
}

func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		fmt.Println("Error closing file:", err)
	}
}
