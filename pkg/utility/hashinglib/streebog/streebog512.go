package streebog

// GOST R 34.11-2012 (Streebog) hash implementation.
// This file implements the 512-bit variant. For 256-bit output, hash and then take the lower 32 bytes.

import (
	"encoding/binary"
)

const (
	blockSize = 64 // 512-bit blocks
)

// iv512 is the initial vector for Streebog-512 (all zeros).
var iv512 = [64]byte{}

// tau permutation
var tau = [64]byte{
	0, 8, 16, 24, 32, 40, 48, 56,
	1, 9, 17, 25, 33, 41, 49, 57,
	2, 10, 18, 26, 34, 42, 50, 58,
	3, 11, 19, 27, 35, 43, 51, 59,
	4, 12, 20, 28, 36, 44, 52, 60,
	5, 13, 21, 29, 37, 45, 53, 61,
	6, 14, 22, 30, 38, 46, 54, 62,
	7, 15, 23, 31, 39, 47, 55, 63,
}

// pi S-box (from the standard)
var pi = [256]byte{
	0xFC, 0xEE, 0xDD, 0x11, 0xCF, 0x6E, 0x31, 0x16,
	0xFB, 0xC4, 0xFA, 0xDA, 0x23, 0xC5, 0x04, 0x4D,
	0xE9, 0x77, 0xF0, 0xDB, 0x93, 0x2E, 0x99, 0xBA,
	0x17, 0x36, 0xF1, 0xBB, 0x14, 0xCD, 0x5F, 0xC1,
	0xF9, 0x18, 0x65, 0x5A, 0xE2, 0x5C, 0xEF, 0x21,
	0x81, 0x1C, 0x3C, 0x42, 0x8B, 0x01, 0x8E, 0x4F,
	0x05, 0x84, 0x02, 0xAE, 0xE3, 0x6A, 0x8F, 0xA0,
	0x06, 0x0B, 0xED, 0x98, 0x7F, 0xD4, 0xD3, 0x1F,
	0xEB, 0x34, 0x2C, 0x51, 0xEA, 0xC8, 0x48, 0xAB,
	0xF2, 0x2A, 0x68, 0xA2, 0xFD, 0x3A, 0xCE, 0xCC,
	0xB5, 0x70, 0x0E, 0x56, 0x08, 0x0C, 0x76, 0x12,
	0xBF, 0x72, 0x13, 0x47, 0x9C, 0xB7, 0x5D, 0x87,
	0x15, 0xA1, 0x96, 0x29, 0x10, 0x7B, 0x9A, 0xC7,
	0xF3, 0x91, 0x78, 0x6F, 0x9D, 0x9E, 0xB2, 0xB1,
	0x32, 0x75, 0x19, 0x3D, 0xFF, 0x35, 0x8A, 0x7E,
	0x6D, 0x54, 0xC6, 0x80, 0xC3, 0xBD, 0x0D, 0x57,
	0xDF, 0xF5, 0x24, 0xA9, 0x3E, 0xA8, 0x43, 0xC9,
	0xD7, 0x79, 0xD6, 0xF6, 0x7C, 0x22, 0xB9, 0x03,
	0xE0, 0x0F, 0xEC, 0xDE, 0x7A, 0x94, 0xB0, 0xBC,
	0xDC, 0xE8, 0x28, 0x50, 0x4E, 0x33, 0x0A, 0x4A,
	0xA7, 0x97, 0x60, 0x73, 0x1E, 0x00, 0x62, 0x44,
	0x1A, 0xB8, 0x38, 0x82, 0x64, 0x9F, 0x26, 0x41,
	0xAD, 0x45, 0x46, 0x92, 0x27, 0x5E, 0x55, 0x2F,
	0x8C, 0xA3, 0xA5, 0x7D, 0x69, 0xD5, 0x95, 0x3B,
	0x07, 0x58, 0xB3, 0x40, 0x86, 0xAC, 0x1D, 0xF7,
	0x30, 0x37, 0x6B, 0xE4, 0x88, 0xD9, 0xE7, 0x89,
	0xE1, 0x1B, 0x83, 0x49, 0x4C, 0x3F, 0xF8, 0xFE,
	0x8D, 0x53, 0xAA, 0x90, 0xCA, 0xD8, 0x85, 0x61,
	0x20, 0x71, 0x67, 0xA4, 0x2D, 0x2B, 0x09, 0x5B,
	0xCB, 0x9B, 0x25, 0xD0, 0xBE, 0xE5, 0x6C, 0x52,
	0x59, 0xA6, 0x74, 0xD2, 0xE6, 0xF4, 0xB4, 0xC0,
	0xD1, 0x66, 0xAF, 0xC2, 0x39, 0x4B, 0x63, 0xB6,
}

// A is the binary matrix for L transformation as polynomials over GF(2)
var A = [64]uint64{
	0x8e20faa72ba0b470, 0x47107ddd9b505a38, 0xad08b0e0c3282d1c, 0xd8045870ef14980e,
	0x6c022c38f90a4c07, 0x3601161cf205268d, 0x1b8e0b0e798c13c8, 0x83478b07b2468764,
	0xa011d380818e8f40, 0x5086e740ce47c920, 0x2843fd2067adea10, 0x14aff010bdd87508,
	0x0ad97808d06cb404, 0x05e23c0468365a02, 0x8c711e02341b2d01, 0x46b60f011a83988e,
	0x90dab52a387ae76f, 0x486dd4151c3dfdb9, 0x24b86a840e90f0d2, 0x125c354207487869,
	0x092e94218d243cba, 0x8a174a9ec8121e5d, 0x4585254f64090fa0, 0xaccc9ca9328a8950,
	0x9d4df05d5f661451, 0xc0a878a0a1330aa6, 0x60543c50de970553, 0x302a1e286fc58ca7,
	0x18150f14b9ec46dd, 0x0c84890ad27623e0, 0x0642ca05693b9f70, 0x0321658cba93c138,
	0x86275df09ce8aaa8, 0x439da0784e745554, 0xafc0503c273aa42a, 0xd960281e9d1d5215,
	0xe230140fc0802984, 0x71180a8960409a42, 0xb60c05ca30204d21, 0x5b068c651810a89e,
	0x456c34887a3805b9, 0xac361a443d1c8cd2, 0x561b0d22900e4669, 0x2b838811480723ba,
	0x9bcf4486248d9f5d, 0xc3e9224312c8c1a0, 0xeffa11af0964ee50, 0xf97d86d98a327728,
	0xe4fa2054a80b329c, 0x727d102a548b194e, 0x39b008152acb8227, 0x9258048415eb419d,
	0x492c024284fbaec0, 0xaa16012142f35760, 0x550b8e9e21f7a530, 0xa48b474f9ef5dc18,
	0x70a6a56e2440598e, 0x3853dc371220a247, 0x1ca76e95091051ad, 0x0edd37c48a08a6d8,
	0x07e095624504536c, 0x8d70c431ac02a736, 0xc83862965601dd1b, 0x641c314b2b8ee083,
}

// C constants for 12 rounds
var C = [12][64]byte{}

func init() {
	// Precompute C[i] = L(P(S( (i+1) || zeros ) ))
	for i := 0; i < 12; i++ {
		var x [64]byte
		x[0] = byte(i + 1)
		S(&x)
		P(&x)
		L(&x)
		C[i] = x
	}
}

// S - byte substitution
func S(x *[64]byte) {
	for i := 0; i < 64; i++ {
		x[i] = pi[x[i]]
	}
}

// P - byte permutation
func P(x *[64]byte) {
	var t [64]byte
	for i := 0; i < 64; i++ {
		t[i] = x[tau[i]]
	}
	*x = t
}

// L - linear transformation
func L(x *[64]byte) {
	for i := 0; i < 8; i++ {
		var v uint64
		for j := 0; j < 8; j++ {
			b := x[i*8+j]
			for k := 0; k < 8; k++ {
				if (b>>(7-k))&1 == 1 {
					v ^= A[j*8+k]
				}
			}
		}
		binary.BigEndian.PutUint64(x[i*8:], v)
	}
}

func xor64(a *[64]byte, b *[64]byte) {
	for i := 0; i < 64; i++ {
		a[i] ^= b[i]
	}
}

// E - the encryption-like transformation with key k of input m
func E(k *[64]byte, m *[64]byte) [64]byte {
	state := *m
	key := *k
	for i := 0; i < 12; i++ {
		// state = S ∘ P ∘ L (state XOR key)
		xor64(&state, &key)
		S(&state)
		P(&state)
		L(&state)
		// key = S ∘ P ∘ L (key XOR C[i])
		tmp := key
		xor64(&tmp, &C[i])
		S(&tmp)
		P(&tmp)
		L(&tmp)
		key = tmp
	}
	xor64(&state, &key)
	return state
}

// gN - compression function
func gN(h, N, m *[64]byte) [64]byte {
	// K = L(P(S(h XOR N)))
	K := *h
	xor64(&K, N)
	S(&K)
	P(&K)
	L(&K)

	// t = E(K, m)
	t := E(&K, m)

	// return t XOR h XOR m
	xor := t
	tmp := *h
	xor64(&xor, &tmp)
	tmp = *m
	xor64(&xor, &tmp)
	return xor
}

// Hash512 computes Streebog-512 for the provided message and returns a 64-byte slice.
func Hash512(msg []byte) []byte {
	var h [64]byte = iv512
	var N [64]byte // processed bits mod 2^512
	var Sigma [64]byte

	processBlock := func(block []byte) {
		var m [64]byte
		copy(m[:], block)
		h = gN(&h, &N, &m)

		// N = (N + (len(block) * 8)) mod 2^512
		var t uint16 = uint16(len(block) * 8)
		for i := 63; i >= 0; i-- {
			t += uint16(N[i])
			N[i] = byte(t)
			t >>= 8
		}

		// Sigma = (Sigma + m) mod 2^512
		var c uint16
		for i := 63; i >= 0; i-- {
			c = uint16(Sigma[i]) + uint16(m[i]) + (c >> 8)
			Sigma[i] = byte(c)
		}
	}

	// process full blocks
	i := 0
	for ; i+blockSize <= len(msg); i += blockSize {
		processBlock(msg[i : i+blockSize])
	}

	// final block with padding
	var last [64]byte
	rem := msg[i:]
	padLen := 64 - len(rem)
	// right-justified: zeros then 1 then message at the end
	// copy remainder to the end
	copy(last[64-len(rem):], rem)
	// set single '1' bit before the message
	last[64-len(rem)-1] = 0x01
	_ = padLen // kept for clarity; padding already applied as above

	// process final padded message with N of remaining bits
	var nFinal [64]byte
	bits := uint64(len(rem) * 8)
	for j := 0; j < 8; j++ {
		nFinal[63-j] = byte(bits)
		bits >>= 8
	}
	h = gN(&h, &N, &last)

	// update N with remaining bits
	{
		var c uint16
		for i := 63; i >= 0; i-- {
			c = uint16(N[i]) + uint16(nFinal[i]) + (c >> 8)
			N[i] = byte(c)
		}
	}

	// update Sigma with last
	{
		var c uint16
		for i := 63; i >= 0; i-- {
			c = uint16(Sigma[i]) + uint16(last[i]) + (c >> 8)
			Sigma[i] = byte(c)
		}
	}

	// finalization: gN with N=0 and m = N, then with m = Sigma
	var zero [64]byte
	h = gN(&h, &zero, &N)
	h = gN(&h, &zero, &Sigma)

	// return as slice
	out := make([]byte, 64)
	copy(out, h[:])
	return out
}
