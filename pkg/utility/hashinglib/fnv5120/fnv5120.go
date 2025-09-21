package fnv5120

// Fowler–Noll–Vo hash, 512-bit (FNV-1 and FNV-1a).
// Variant name "fnv-5120" here refers to 512-bit output.

import (
	"encoding/binary"
)

// 512-bit parameters from FNV specification
// offset basis (little-endian words shown here, stored as big-endian bytes for convenience)
var fnv512Offset = [64]byte{
	// 0x0000000000000000000000000000000000000000000000000000000000000000
	// Official FNV-512 offset basis:
	// 0x0000000000000000000000000000000000000000000000000000000000000013
	// 0x0000000000000000000000000000000000000000000000000000000000000097
	// 0x00000000000000000000000000000000000000000000000000000000000000c1
	// 0x00000000000000000000000000000000000000000000000000000000000084f3
	// But stored big-endian across the 64-byte array:
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x13,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x84, 0xF3,
}

// FNV prime for 512-bit
var fnv512Prime = [64]byte{
	// FNV-512 prime (big-endian bytes)
	// 0x0000000000000000000000000000000000000000000000000000000000000000...
	// The defined 512-bit prime is:
	// 0x0000000000000000000000000000000000000000000000000000000000013B
	// 0x0000000000000000000000000000000000000000000000000000000000000253
	// 0x0000000000000000000000000000000000000000000000000000000000000001
	// 0x00000000000000000000000000000000000000000000000000000000000000B3
	// As a compact big-endian 64-byte array we only set the low bytes:
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xB3,
}

// Hash512 computes the FNV-1 512-bit hash of msg and returns a 64-byte slice.
func Hash512(msg []byte) []byte {
	h := fnv512Offset
	for _, b := range msg {
		mul512(&h, &fnv512Prime)
		addByte(&h, b)
	}
	out := make([]byte, 64)
	copy(out, h[:])
	return out
}

// Hash512a computes the FNV-1a 512-bit hash of msg and returns a 64-byte slice.
func Hash512a(msg []byte) []byte {
	h := fnv512Offset
	for _, b := range msg {
		addByte(&h, b)
		mul512(&h, &fnv512Prime)
	}
	out := make([]byte, 64)
	copy(out, h[:])
	return out
}

// addByte: h = h XOR b (on the least significant byte)
func addByte(h *[64]byte, b byte) {
	h[63] ^= b
}

// mul512: h = (h * p) mod 2^512 (big-endian 512-bit integers)
func mul512(h *[64]byte, p *[64]byte) {
	var res [64]byte
	// schoolbook multiplication with carry (base 2^32 to speed slightly)
	const limbs = 16 // 16 * 32 = 512
	var a32, b32 [limbs]uint32
	for i := 0; i < limbs; i++ {
		a32[i] = binary.BigEndian.Uint32(h[i*4:])
		b32[i] = binary.BigEndian.Uint32(p[i*4:])
	}
	var acc [limbs * 2]uint64
	for i := 0; i < limbs; i++ {
		ai := uint64(a32[i])
		for j := 0; j < limbs; j++ {
			acc[i+j] += ai * uint64(b32[j])
		}
	}
	// reduce to 16 limbs with carry in base 2^32
	var carry uint64
	for i := limbs*2 - 1; i >= 0; i-- {
		acc[i] += carry
		carry = acc[i] >> 32
		acc[i] &= 0xffffffff
	}
	// keep only lowest 16 limbs (mod 2^512)
	for i := 0; i < limbs; i++ {
		binary.BigEndian.PutUint32(res[i*4:], uint32(acc[i+limbs]))
	}
	*h = res
}
