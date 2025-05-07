package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Define IPv4 ranges for classes A, B, C, D, and E
	ranges := map[string][2]string{
		"A": {"1.0.0.0", "127.0.0.0"},
		"B": {"128.0.0.0", "191.255.255.255"},
		"C": {"192.0.0.0", "223.255.255.255"},
		"D": {"224.0.0.0", "239.255.255.255"},
		"E": {"240.0.0.0", "255.255.255.255"},
	}

	// Prompt user for input
	var class string
	fmt.Println("Enter the IPv4 class range (A, B, C, D, E):")
	_, scanErr := fmt.Scanln(&class)
	if scanErr != nil {
		return
	}
	class = strings.ToUpper(class)

	// Validate input
	if _, exists := ranges[class]; !exists {
		fmt.Println("Invalid class. Please enter A, B, C, D, or E.")
		return
	}

	// Get the range for the selected class
	startIP := net.ParseIP(ranges[class][0]).To4()
	endIP := net.ParseIP(ranges[class][1]).To4()

	// Calculate the total number of IPs in the range
	totalIPs := ipToInt(endIP) - ipToInt(startIP) + 1
	ipsPerFile := totalIPs / 8

	// Split the range and write to 8 files
	for i := 0; i < 8; i++ {
		file, failed := createFile(class, "ips", i)
		if failed {
			return
		}

		// Write IPs to the file
		start := ipToInt(startIP) + int64(i)*ipsPerFile
		end := start + ipsPerFile - 1
		if i == 7 { // Ensure the last file includes any remaining IPs
			end = ipToInt(endIP)
		}

		// Just IPs in the range
		for ip := start; ip <= end; ip++ {
			_, err := file.WriteString(intToIP(ip).String() + "\n")
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		}

		closeFile(file)
		file, failed = createFile(class, "ipswport", i)
		if failed {
			return
		}

		// IPs with port numbers
		for ip := start; ip <= end; ip++ {
			for port := 1; port <= 65535; port++ {
				ipWithPort := fmt.Sprintf("%s:%d\n", intToIP(ip).String(), port)
				_, err := file.WriteString(ipWithPort)
				if err != nil {
					fmt.Println("Error writing to file:", err)
					return
				}
			}
		}

		closeFile(file)
		file, failed = createFile(class, "ipswscheme", i)
		if failed {
			return
		}

		// Just IPs in the range with schemes
		for ip := start; ip <= end; ip++ {
			for _, scheme := range getSchemes() {
				ipWithScheme := fmt.Sprintf("%s://%s\n", scheme, intToIP(ip).String())
				_, err := file.WriteString(ipWithScheme)
				if err != nil {
					fmt.Println("Error writing to file:", err)
					return
				}

				ipWithSchemeEnd := fmt.Sprintf("%s://%s/\n", scheme, intToIP(ip).String())
				_, err = file.WriteString(ipWithSchemeEnd)
				if err != nil {
					fmt.Println("Error writing to file:", err)
					return
				}
			}
		}

		closeFile(file)
		file, failed = createFile(class, "ipswportwscheme", i)
		if failed {
			return
		}

		// IPs with port numbers with schemes
		for ip := start; ip <= end; ip++ {
			for _, scheme := range getSchemes() {
				for port := 1; port <= 65535; port++ {
					ipWithPortWithScheme := fmt.Sprintf("%s://%s:%d\n", scheme, intToIP(ip).String(), port)
					_, err := file.WriteString(ipWithPortWithScheme)
					if err != nil {
						fmt.Println("Error writing to file:", err)
						return
					}

					ipWithPortWithSchemeEnd := fmt.Sprintf("%s://%s:%d/\n", scheme, intToIP(ip).String(), port)
					_, err = file.WriteString(ipWithPortWithSchemeEnd)
					if err != nil {
						fmt.Println("Error writing to file:", err)
						return
					}
				}
			}
		}

		closeFile(file)

		fmt.Print("Files written successfully.\n")
	}
}

func closeFile(file *os.File) {
	closeError := file.Close()
	if closeError != nil {
		fmt.Println("Error closing file:", closeError)
	}
}

func createFile(class, portion string, i int) (*os.File, bool) {
	fileName := fmt.Sprintf("%s_%s_range_part_%d.txt", class, portion, i+1)
	file, createError := os.Create(fileName)
	if createError != nil {
		fmt.Println("Error creating file:", createError)
		return nil, true
	}
	return file, false
}
