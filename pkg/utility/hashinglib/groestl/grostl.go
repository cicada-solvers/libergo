package groestl

import (
	"encoding/binary"
	"hash"
)

// Groestl-512 parameters
const (
	gr512BlockSize = 64 // 512-bit message block
	gr512OutSize   = 64 // 512-bit output
	gr512Cols      = 8  // 8 columns of 64 bits (state is 8x8 bytes)
	gr512Rows      = 8  // 8 rows
	gr512Rounds    = 14 // number of rounds for Groestl-512
)

// Groestl512 implements hash.Hash for Groestl-512.
type Groestl512 struct {
	h      [gr512Rows][gr512Cols]byte // chaining value H (state matrix 8x8)
	buf    [gr512BlockSize]byte
	bufLen int
	len    uint64 // total message length in bytes
}

// NewGroestl512 returns a new Groestl-512 hasher.
func NewGroestl512() hash.Hash {
	var g Groestl512
	g.Reset()
	return &g
}

func (g *Groestl512) Size() int      { return gr512OutSize }
func (g *Groestl512) BlockSize() int { return gr512BlockSize }

func (g *Groestl512) Reset() {
	// IV is 8x8 with first row containing output size (in bits) in big-endian (per spec), rest zeros.
	// For 512-bit output, IV = 0.. except h[0] = 0x00..02 00 (i.e., 512) at last two bytes in big-endian position.
	for r := 0; r < gr512Rows; r++ {
		for c := 0; c < gr512Cols; c++ {
			g.h[r][c] = 0
		}
	}
	// Put 512 (bit length) as 64-bit big-endian into first row.
	var tmp [8]byte
	binary.BigEndian.PutUint64(tmp[:], 512)
	for c := 0; c < gr512Cols; c++ {
		g.h[0][c] = tmp[c]
	}
	g.bufLen = 0
	g.len = 0
}

func (g *Groestl512) Write(p []byte) (int, error) {
	n := len(p)
	g.len += uint64(n)
	if g.bufLen > 0 {
		r := gr512BlockSize - g.bufLen
		if r > n {
			r = n
		}
		copy(g.buf[g.bufLen:], p[:r])
		g.bufLen += r
		p = p[r:]
		if g.bufLen == gr512BlockSize {
			g.processBlock(g.buf[:])
			g.bufLen = 0
		}
	}
	for len(p) >= gr512BlockSize {
		g.processBlock(p[:gr512BlockSize])
		p = p[gr512BlockSize:]
	}
	if len(p) > 0 {
		copy(g.buf[:], p)
		g.bufLen = len(p)
	}
	return n, nil
}

func (g *Groestl512) Sum(in []byte) []byte {
	cp := *g
	cp.finalize()
	var out [gr512OutSize]byte
	// Output transformation: Ω(H) = P(H) ⊕ H, then take rightmost 64 bytes (whole state here)
	omega := cp.permutationP(cp.h)
	for r := 0; r < gr512Rows; r++ {
		for c := 0; c < gr512Cols; c++ {
			omega[r][c] ^= cp.h[r][c]
		}
	}
	// Serialize state by rows (rightmost 64 bytes means whole 8x8 matrix in row-major)
	idx := 0
	for r := 0; r < gr512Rows; r++ {
		for c := 0; c < gr512Cols; c++ {
			out[idx] = omega[r][c]
			idx++
		}
	}
	return append(in, out[:]...)
}

func (g *Groestl512) finalize() {
	// Pad: append 0x80 then zeros so that last 16 bytes can hold message length in bits (big-endian)
	var block [gr512BlockSize]byte
	copy(block[:], g.buf[:g.bufLen])
	block[g.bufLen] = 0x80

	if g.bufLen+1 > gr512BlockSize-16 {
		// not enough space for length, process this block zero-padded
		g.processBlock(block[:])
		for i := range block {
			block[i] = 0
		}
	}

	// total length in bits
	msgBits := g.len * 8
	binary.BigEndian.PutUint64(block[gr512BlockSize-8:], msgBits)
	// The preceding 8 bytes (from -16 to -8) are zero for Groestl-512.
	// Already zero due to zeroed block.
	g.processBlock(block[:])
}

// Compression function: H' = P(M ⊕ H) ⊕ Q(M) ⊕ H
func (g *Groestl512) processBlock(b []byte) {
	// Convert to state matrices
	M := bytesToState(b)
	H := g.h
	// M ⊕ H
	MxH := xorState(M, H)

	PH := g.permutationP(MxH)
	QM := g.permutationQ(M)

	// H' = PH ⊕ QM ⊕ H
	Hp := xorState(PH, QM)
	Hp = xorState(Hp, H)
	g.h = Hp
}

// State is 8x8 bytes. Rows x Cols.
type state [gr512Rows][gr512Cols]byte

func bytesToState(b []byte) state {
	var s state
	// Row-major: 8 rows, 8 cols; b[0:8] is row0, etc.
	for r := 0; r < gr512Rows; r++ {
		copy(s[r][:], b[r*8:(r+1)*8])
	}
	return s
}

func stateToBytes(s state) [64]byte {
	var out [64]byte
	for r := 0; r < gr512Rows; r++ {
		copy(out[r*8:(r+1)*8], s[r][:])
	}
	return out
}

// AES S-box
var sbox = [256]byte{
	0x63, 0x7c, 0x77, 0x7b, 0xf2, 0x6b, 0x6f, 0xc5, 0x30, 0x01, 0x67, 0x2b, 0xfe, 0xd7, 0xab, 0x76,
	0xca, 0x82, 0xc9, 0x7d, 0xfa, 0x59, 0x47, 0xf0, 0xad, 0xd4, 0xa2, 0xaf, 0x9c, 0xa4, 0x72, 0xc0,
	0xb7, 0xfd, 0x93, 0x26, 0x36, 0x3f, 0xf7, 0xcc, 0x34, 0xa5, 0xe5, 0xf1, 0x71, 0xd8, 0x31, 0x15,
	0x04, 0xc7, 0x23, 0xc3, 0x18, 0x96, 0x05, 0x9a, 0x07, 0x12, 0x80, 0xe2, 0xeb, 0x27, 0xb2, 0x75,
	0x09, 0x83, 0x2c, 0x1a, 0x1b, 0x6e, 0x5a, 0xa0, 0x52, 0x3b, 0xd6, 0xb3, 0x29, 0xe3, 0x2f, 0x84,
	0x53, 0xd1, 0x00, 0xed, 0x20, 0xfc, 0xb1, 0x5b, 0x6a, 0xcb, 0xbe, 0x39, 0x4a, 0x4c, 0x58, 0xcf,
	0xd0, 0xef, 0xaa, 0xfb, 0x43, 0x4d, 0x33, 0x85, 0x45, 0xf9, 0x02, 0x7f, 0x50, 0x3c, 0x9f, 0xa8,
	0x51, 0xa3, 0x40, 0x8f, 0x92, 0x9d, 0x38, 0xf5, 0xbc, 0xb6, 0xda, 0x21, 0x10, 0xff, 0xf3, 0xd2,
	0xcd, 0x0c, 0x13, 0xec, 0x5f, 0x97, 0x44, 0x17, 0xc4, 0xa7, 0x7e, 0x3d, 0x64, 0x5d, 0x19, 0x73,
	0x60, 0x81, 0x4f, 0xdc, 0x22, 0x2a, 0x90, 0x88, 0x46, 0xee, 0xb8, 0x14, 0xde, 0x5e, 0x0b, 0xdb,
	0xe0, 0x32, 0x3a, 0x0a, 0x49, 0x06, 0x24, 0x5c, 0xc2, 0xd3, 0xac, 0x62, 0x91, 0x95, 0xe4, 0x79,
	0xe7, 0xc8, 0x37, 0x6d, 0x8d, 0xd5, 0x4e, 0xa9, 0x6c, 0x56, 0xf4, 0xea, 0x65, 0x7a, 0xae, 0x08,
	0xba, 0x78, 0x25, 0x2e, 0x1c, 0xa6, 0xb4, 0xc6, 0xe8, 0xdd, 0x74, 0x1f, 0x4b, 0xbd, 0x8b, 0x8a,
	0x70, 0x3e, 0xb5, 0x66, 0x48, 0x03, 0xf6, 0x0e, 0x61, 0x35, 0x57, 0xb9, 0x86, 0xc1, 0x1d, 0x9e,
	0xe1, 0xf8, 0x98, 0x11, 0x69, 0xd9, 0x8e, 0x94, 0x9b, 0x1e, 0x87, 0xe9, 0xce, 0x55, 0x28, 0xdf,
	0x8c, 0xa1, 0x89, 0x0d, 0xbf, 0xe6, 0x42, 0x68, 0x41, 0x99, 0x2d, 0x0f, 0xb0, 0x54, 0xbb, 0x16,
}

// MixBytes uses the MDS matrix of Groestl (similar to AES MixColumns but different constants)
func mixBytes(s state) state {
	var out state
	for c := 0; c < gr512Cols; c++ {
		b0 := s[0][c]
		b1 := s[1][c]
		b2 := s[2][c]
		b3 := s[3][c]
		b4 := s[4][c]
		b5 := s[5][c]
		b6 := s[6][c]
		b7 := s[7][c]
		out[0][c] = gmul2(b0) ^ gmul2(b1) ^ gmul3(b2) ^ gmul4(b3) ^ gmul5(b4) ^ b5 ^ b6 ^ b7
		out[1][c] = b0 ^ gmul2(b1) ^ gmul2(b2) ^ gmul3(b3) ^ gmul4(b4) ^ gmul5(b5) ^ b6 ^ b7
		out[2][c] = b0 ^ b1 ^ gmul2(b2) ^ gmul2(b3) ^ gmul3(b4) ^ gmul4(b5) ^ gmul5(b6) ^ b7
		out[3][c] = b0 ^ b1 ^ b2 ^ gmul2(b3) ^ gmul2(b4) ^ gmul3(b5) ^ gmul4(b6) ^ gmul5(b7)
		out[4][c] = gmul5(b0) ^ b1 ^ b2 ^ b3 ^ gmul2(b4) ^ gmul2(b5) ^ gmul3(b6) ^ gmul4(b7)
		out[5][c] = gmul4(b0) ^ gmul5(b1) ^ b2 ^ b3 ^ b4 ^ gmul2(b5) ^ gmul2(b6) ^ gmul3(b7)
		out[6][c] = gmul3(b0) ^ gmul4(b1) ^ gmul5(b2) ^ b3 ^ b4 ^ b5 ^ gmul2(b6) ^ gmul2(b7)
		out[7][c] = gmul2(b0) ^ gmul3(b1) ^ gmul4(b2) ^ gmul5(b3) ^ b4 ^ b5 ^ b6 ^ gmul2(b7)
	}
	return out
}

// Galois field multiplications over GF(2^8) with modulus x^8 + x^4 + x^3 + x + 1 (0x11B)
func xtime(x byte) byte {
	v := uint16(x) << 1
	if (x & 0x80) != 0 {
		v ^= 0x11b
	}
	return byte(v & 0xFF)
}
func gmul2(x byte) byte { return xtime(x) }
func gmul3(x byte) byte { return gmul2(x) ^ x }
func gmul4(x byte) byte { return gmul2(gmul2(x)) }
func gmul5(x byte) byte { return gmul4(x) ^ x }

// ShiftBytes for P permutation (row shifts to the left by offsets {0,1,2,3,4,5,6,7})
func shiftBytesP(s state) state {
	var out state
	for r := 0; r < gr512Rows; r++ {
		shift := r
		for c := 0; c < gr512Cols; c++ {
			out[r][c] = s[r][(c+shift)%gr512Cols]
		}
	}
	return out
}

// ShiftBytes for Q permutation (row shifts to the right by offsets {1,2,3,4,5,6,7,0}) and byte complement on even rows initial additive constant
func shiftBytesQ(s state) state {
	var out state
	for r := 0; r < gr512Rows; r++ {
		shift := (r + 1) % gr512Cols
		for c := 0; c < gr512Cols; c++ {
			out[r][c] = s[r][(c+gr512Cols-shift)%gr512Cols]
		}
	}
	return out
}

// SubBytes
func subBytes(s state) state {
	var out state
	for r := 0; r < gr512Rows; r++ {
		for c := 0; c < gr512Cols; c++ {
			out[r][c] = sbox[s[r][c]]
		}
	}
	return out
}

// AddRoundConstant for P: XOR round number to first column increasing values per row
func addRoundConstantP(s state, r int) state {
	var out state
	for i := 0; i < gr512Rows; i++ {
		copy(out[i][:], s[i][:])
	}
	for i := 0; i < gr512Rows; i++ {
		out[i][0] ^= byte((i << 4) ^ r)
	}
	return out
}

// AddRoundConstant for Q: XOR round number and 0xFF pattern
func addRoundConstantQ(s state, r int) state {
	var out state
	for i := 0; i < gr512Rows; i++ {
		copy(out[i][:], s[i][:])
	}
	for i := 0; i < gr512Rows; i++ {
		for c := 0; c < gr512Cols; c++ {
			out[i][c] ^= 0xFF
		}
		out[i][0] ^= byte(((7 - i) << 4) ^ r)
	}
	return out
}

// Permutation P
func (g *Groestl512) permutationP(a state) state {
	s := a
	for r := 0; r < gr512Rounds; r++ {
		s = addRoundConstantP(s, r)
		s = subBytes(s)
		s = shiftBytesP(s)
		s = mixBytes(s)
	}
	return s
}

// Permutation Q
func (g *Groestl512) permutationQ(a state) state {
	s := a
	for r := 0; r < gr512Rounds; r++ {
		s = addRoundConstantQ(s, r)
		s = subBytes(s)
		s = shiftBytesQ(s)
		s = mixBytes(s)
	}
	return s
}

// Utility to xor two states
func xorState(a, b state) state {
	var out state
	for r := 0; r < gr512Rows; r++ {
		for c := 0; c < gr512Cols; c++ {
			out[r][c] = a[r][c] ^ b[r][c]
		}
	}
	return out
}

// Convenience helper
func SumGroestl512(data []byte) [64]byte {
	h := NewGroestl512()
	_, _ = h.Write(data)
	sum := h.Sum(nil)
	var out [64]byte
	copy(out[:], sum)
	return out
}
