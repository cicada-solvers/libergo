package keccak3

import (
	"crypto/sha3"
)

// Keccak3_512 computes the Keccak3-512 (SHA3-512) hash of the input.
// It returns a 64-byte (512-bit) digest.
func Keccak3_512(data []byte) []byte {
	sum := sha3.Sum512(data)
	// Return a new slice to avoid exposing the array-backed value
	out := make([]byte, len(sum))
	copy(out, sum[:])
	return out
}

// Keccak3_512Hex computes the Keccak3-512 (SHA3-512) hash and returns
// the lowercase hexadecimal string encoding of the digest.
func Keccak3_512Hex(data []byte) string {
	sum := sha3.Sum512(data)
	// Manually hex-encode to avoid adding extra dependencies
	const hextable = "0123456789abcdef"
	dst := make([]byte, 128)
	for i, b := range sum {
		dst[i*2] = hextable[b>>4]
		dst[i*2+1] = hextable[b&0x0f]
	}
	return string(dst)
}
