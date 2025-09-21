package cube

// CubeHash implementation (CubeHash16/32-256 and CubeHash16/32-512).
// Parameters:
//   - r = 16 (rounds per block)
//   - b = 32 (bytes per block)
//   - h = 32 or 64 (output size in bytes)
//
// The code provides one-shot and streaming APIs.

import (
	"encoding/binary"
)

// cubeState holds 32 little-endian 32-bit words (1024-bit state).
type cubeState struct {
	v [32]uint32
}

// quarterRound performs one CubeHash inner round step.
func (s *cubeState) quarterRound() {
	// This follows the reference round structure:
	// for i=0..15: x[i] += x[i+16]
	// for i=0..15: x[i+16] = ROTL32(x[i+16],7)
	// permute words, then XORs, then rotates and permutes again
	// We implement the standard CubeHash round used in Bernstein's submission.

	// 1) x[i] += x[i+16]
	for i := 0; i < 16; i++ {
		s.v[i] += s.v[i+16]
	}
	// 2) x[i+16] = ROTL32(x[i+16], 7)
	for i := 16; i < 32; i++ {
		s.v[i] = (s.v[i] << 7) | (s.v[i] >> (32 - 7))
	}
	// 3) Swap pairs (this is the "shuffle" step)
	for i := 0; i < 8; i++ {
		s.v[i+16], s.v[i+24] = s.v[i+24], s.v[i+16]
	}
	// 4) x[i+16] ^= x[i]
	for i := 0; i < 16; i++ {
		s.v[i+16] ^= s.v[i]
	}
	// 5) Swap quads
	for i := 0; i < 4; i++ {
		base := i * 4
		s.v[base+16], s.v[base+20] = s.v[base+20], s.v[base+16]
		s.v[base+17], s.v[base+21] = s.v[base+21], s.v[base+17]
		s.v[base+18], s.v[base+22] = s.v[base+22], s.v[base+18]
		s.v[base+19], s.v[base+23] = s.v[base+23], s.v[base+19]
	}
	// 6) x[i] += x[i+16]
	for i := 0; i < 16; i++ {
		s.v[i] += s.v[i+16]
	}
	// 7) x[i+16] = ROTL32(x[i+16], 11)
	for i := 16; i < 32; i++ {
		s.v[i] = (s.v[i] << 11) | (s.v[i] >> (32 - 11))
	}
	// 8) Swap pairs again (different pattern)
	for i := 0; i < 8; i++ {
		s.v[i], s.v[i+8] = s.v[i+8], s.v[i]
	}
	// 9) x[i+16] ^= x[i]
	for i := 0; i < 16; i++ {
		s.v[i+16] ^= s.v[i]
	}
}

// rounds applies r rounds.
func (s *cubeState) rounds(r int) {
	for i := 0; i < r; i++ {
		s.quarterRound()
	}
}

// cubeHash implements the streaming state for CubeHash with fixed parameters.
type cubeHash struct {
	st    cubeState
	r     int // rounds per block
	b     int // block size in bytes (32)
	h     int // output size in bytes (32 or 64)
	buf   [32]byte
	n     int
	total uint64
}

// initCube initializes the state with parameters (r,b,h).
// Follows the spec: set state to zero, set x[0]=h/8, x[1]=b, x[2]=r, then 10*r rounds.
func (c *cubeHash) initCube(r, b, h int) {
	c.r, c.b, c.h = r, b, h
	for i := range c.st.v {
		c.st.v[i] = 0
	}
	// Parameter inject (little-endian words):
	// x[0] = h (in bits) / 8 -> output bytes, but spec uses output bytes in x[0] for some variants.
	// Widely used setting: x[0]=h, x[1]=b, x[2]=r for the reference parameterization.
	c.st.v[0] = uint32(h)
	c.st.v[1] = uint32(b)
	c.st.v[2] = uint32(r)
	c.st.rounds(10 * r)
}

// blockAbsorb XORs the 32-byte block into the first 8 words and performs r rounds.
func (c *cubeHash) blockAbsorb(block []byte) {
	// XOR little-endian 32 bytes into x[0..7]
	for i := 0; i < 8; i++ {
		w := binary.LittleEndian.Uint32(block[i*4:])
		c.st.v[i] ^= w
	}
	c.st.rounds(c.r)
}

// Write absorbs input bytes.
func (c *cubeHash) Write(p []byte) {
	c.total += uint64(len(p))
	for len(p) > 0 {
		if c.n == c.b {
			c.blockAbsorb(c.buf[:])
			c.n = 0
		}
		n := copy(c.buf[c.n:], p)
		c.n += n
		p = p[n:]
	}
}

// finalize performs padding and finalization, producing h bytes.
func (c *cubeHash) finalize() []byte {
	// Padding: append 0x80, then zeros, then process final block.
	c.buf[c.n] = 0x80
	for i := c.n + 1; i < c.b; i++ {
		c.buf[i] = 0
	}
	c.blockAbsorb(c.buf[:])
	// Finalization bit: flip a bit in state (spec: x[31] ^= 1) then 10*r rounds.
	c.st.v[31] ^= 1
	c.st.rounds(10 * c.r)

	// Output: serialize x[0..(h/4-1)] as little-endian.
	outWords := c.h / 4
	out := make([]byte, c.h)
	for i := 0; i < outWords; i++ {
		binary.LittleEndian.PutUint32(out[i*4:], c.st.v[i])
	}
	return out
}

// CubeHash256 computes CubeHash16/32-256 (32-byte digest) in one shot.
func CubeHash256(data []byte) []byte {
	var c cubeHash
	c.initCube(16, 32, 32)
	c.Write(data)
	return c.finalize()
}

// CubeHash512 computes CubeHash16/32-512 (64-byte digest) in one shot.
func CubeHash512(data []byte) []byte {
	var c cubeHash
	c.initCube(16, 32, 64)
	c.Write(data)
	return c.finalize()
}

// CubeHash256Hex returns CubeHash256 as lowercase hex.
func CubeHash256Hex(data []byte) string {
	sum := CubeHash256(data)
	const hextable = "0123456789abcdef"
	dst := make([]byte, len(sum)*2)
	for i, v := range sum {
		dst[i*2] = hextable[v>>4]
		dst[i*2+1] = hextable[v&0x0f]
	}
	return string(dst)
}

// CubeHash512Hex returns CubeHash512 as lowercase hex.
func CubeHash512Hex(data []byte) string {
	sum := CubeHash512(data)
	const hextable = "0123456789abcdef"
	dst := make([]byte, len(sum)*2)
	for i, v := range sum {
		dst[i*2] = hextable[v>>4]
		dst[i*2+1] = hextable[v&0x0f]
	}
	return string(dst)
}
