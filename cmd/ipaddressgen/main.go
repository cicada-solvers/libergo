package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Define IPv4 ranges for classes A, B, C, and D
	ranges := map[string][2]string{
		"A": {"1.0.0.0", "126.255.255.255"},
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
		fileName := fmt.Sprintf("range_part_%d.txt", i+1)
		file, createError := os.Create(fileName)
		if createError != nil {
			fmt.Println("Error creating file:", createError)
			return
		}

		// Write IPs to the file
		start := ipToInt(startIP) + int64(i)*ipsPerFile
		end := start + ipsPerFile - 1
		if i == 7 { // Ensure the last file includes any remaining IPs
			end = ipToInt(endIP)
		}
		for ip := start; ip <= end; ip++ {
			_, err := file.WriteString(intToIP(ip).String() + "\n")
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		}

		closeError := file.Close()
		if closeError != nil {
			fmt.Println("Error closing file:", closeError)
		}

		fmt.Printf("File %s written successfully.\n", fileName)
	}
}

// Convert an IP address to an integer
func ipToInt(ip net.IP) int64 {
	var result int64
	for _, b := range ip {
		result = result<<8 + int64(b)
	}
	return result
}

// Convert an integer to an IP address
func intToIP(ipInt int64) net.IP {
	return net.IPv4(byte(ipInt>>24), byte(ipInt>>16), byte(ipInt>>8), byte(ipInt))
}
