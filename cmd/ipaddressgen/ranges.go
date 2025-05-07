package main

import "net"

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
