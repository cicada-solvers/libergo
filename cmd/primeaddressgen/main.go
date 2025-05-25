package main

import (
	"flag"
	"fmt"
	"math/big"
	"net"
	"sequences"
	"strings"
)

func main() {
	ranges := map[string][2]string{
		"A": {"1.0.0.0", "127.0.0.0"},
		"B": {"128.0.0.0", "191.255.255.255"},
		"C": {"192.0.0.0", "223.255.255.255"},
		"D": {"224.0.0.0", "239.255.255.255"},
		"E": {"240.0.0.0", "255.255.255.255"},
	}

	classPtr := flag.String("class", "", "IPv4 class range (A, B, C, D, E)")
	flag.Parse()

	class := strings.ToUpper(*classPtr)

	if _, exists := ranges[class]; !exists {
		fmt.Println("Invalid class. Please enter A, B, C, D, or E.")
		return
	}

	startIP := net.ParseIP(ranges[class][0]).To4()
	endIP := net.ParseIP(ranges[class][1]).To4()

	start := ipToInt(startIP)
	end := ipToInt(endIP)

	var primePortsToCheck []int64
	for i := 0; i < 65535; i++ {
		val := big.NewInt(int64(i))
		if sequences.IsPrime(val) {
			primePortsToCheck = append(primePortsToCheck, int64(i))
		}
	}

	for i := start; i <= end; i++ {
		ip := intToIP(i)
		ipString := ip.String()
		if checkIpString(ipString) {
			for _, port := range primePortsToCheck {
				// Write the address to the file
				fmt.Printf("%s:%d\n", ipString, port)
			}
		}
	}

	fmt.Println("Processed range successfully.")
}

func checkIpString(ip string) bool {
	// Check each octet is prime
	for _, octet := range strings.Split(ip, ".") {
		val, _ := big.NewInt(0).SetString(octet, 10)
		if !sequences.IsPrime(val) {
			return false
		}
	}
	return true
}

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
