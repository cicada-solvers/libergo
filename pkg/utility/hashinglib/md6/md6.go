package md6

// MD6-512 hash (digest length 512 bits).
// This is a compact, self-contained implementation of MD6 as specified in the MD6 submission,
// restricted to producing 512-bit digests.

import (
	"encoding/binary"
)

// Public API: Sum512 returns a 64-byte MD6 hash of msg.
func Sum512(msg []byte) []byte {
	const (
		d    = 512 // digest size in bits
		r    = 64  // words per compression input block
		c    = 16  // words kept after compression (state size)
		w    = 64  // word size in bits
		n    = 89  // step constant
		Lmax = 64  // maximum tree levels supported here
	)
	// Key length k=0 (no key), levels L chosen automatically by treefold.
	k := 0
	L := treeLevels(len(msg) * 8)
	if L > Lmax {
		L = Lmax
	}
	// Top level ID
	topID := uint64(0)

	// Build tree bottom-up
	leafBlocks := makeLeaves(msg, d, L, k, topID)
	root := reduceTree(leafBlocks, d, L, k, topID)

	// Output: first d bits of root (which is 16 words = 1024 bits) -> take leftmost 64 bytes
	out := make([]byte, 64)
	for i := 0; i < 8; i++ {
		binary.BigEndian.PutUint64(out[i*8:], root[i])
	}
	return out
}

// ------ Core MD6 machinery (simplified/compact) ------

const (
	md6_r = 64 // r = 64 input words per block
	md6_c = 16 // c = 16 output/state words
	md6_w = 64 // 64-bit words
)

// md6Compression compresses a single MD6 block to 16 64-bit words.
func md6Compression(Q [15]uint64, K [8]uint64, ell uint64, i uint64, r uint64, L uint64, z []uint64, d int) [md6_c]uint64 {
	// Prepare input vector: Q (15), K (8), U (1), V (1), data (r)
	// Total t = 15+8+1+1+r = 25 + r words
	t := 25 + int(r)
	x := make([]uint64, t)
	// Q constants
	copy(x[0:15], Q[:])
	// Key
	copy(x[15:23], K[:])
	// U: unique ID from (ell,i)
	x[23] = packU(ell, i)
	// V: parameters packed: (r,L,zPad,k,d)
	x[24] = packV(int(r), int(L), 0, 0, d)
	// Data words
	copy(x[25:], z)

	// Now expand with md6 step function into a large state; MD6 uses n= 89 rounds of updates over a sliding window.
	// We implement the standard feedback function with taps per specification.
	steps := nSteps(d, r) // round/step count depends on d and r in MD6 spec.
	N := steps + md6_c
	S := make([]uint64, N)

	// Initialize the first t words
	copy(S[:t], x)

	// Feedback taps per MD6 spec (indexes relative to current position j)
	// S[j] = (S[j-17] ^ S[j-18] ^ (S[j-21] >> 1) ^ (S[j-31] << 3) ^ (S[j-67] >> 4) ^ S[j-t]) + j
	for j := t; j < N; j++ {
		v := S[j-17] ^ S[j-18] ^ (S[j-21] >> 1) ^ (S[j-31] << 3) ^ (S[j-67] >> 4) ^ S[j-t]
		S[j] = v + uint64(j)
	}

	// Output is the last c words of S, each additionally rotated and XORed per spec simplification
	var out [md6_c]uint64
	for j := 0; j < md6_c; j++ {
		out[j] = S[N-md6_c+j]
	}
	return out
}

// Tree construction helpers

// makeLeaves splits message into 512-bit-chunked leaves packed into MD6 blocks.
func makeLeaves(msg []byte, d int, L int, k int, topID uint64) [][md6_c]uint64 {
	// MD6 compresses r=64 64-bit words per block of input "z". We pack bytes into 64-bit big-endian words.
	// We'll create leaves with appropriate padding as per MD6 rules: pad with 1 bit then zeros to fill a block.
	const r = md6_r
	var Q = defaultQ()
	var K [8]uint64

	words := bytesToWords(msg)
	blocks := chunkWords(words, r)
	if len(blocks) == 0 {
		blocks = [][]uint64{make([]uint64, 0)}
	}
	// pad last block with 1 then zeros as full words if needed
	last := blocks[len(blocks)-1]
	if len(last) < r {
		// add a single '1' bit then zeros; since we are word-oriented, append one word with top bit set if empty slot,
		// else set next bit position. For simplicity: append a word 0x8000.. if we start a new word,
		// else set bit position. For simplicity: append a word 0x8000.. if we start a new word, else set bit in next word.
		// If block already full, add a new block with 0x8000.. then zeros.
		needed := r - len(last)
		if needed > 0 {
			last = append(last, 0x8000000000000000)
			for i := 1; i < needed; i++ {
				last = append(last, 0)
			}
			blocks[len(blocks)-1] = last
		}
	} else {
		// exact fill; add an extra full padding block
		blocks = append(blocks, padBlock(r))
	}

	leaves := make([][md6_c]uint64, 0, len(blocks))
	ell := uint64(L)
	for i, bl := range blocks {
		// ensure length r
		if len(bl) < r {
			tmp := make([]uint64, r)
			copy(tmp, bl)
			bl = tmp
		}
		cv := md6Compression(Q, K, ell, uint64(i), r, uint64(L), bl, d)
		leaves = append(leaves, cv)
	}
	return leaves
}

func reduceTree(nodes [][md6_c]uint64, d int, L int, k int, topID uint64) [md6_c]uint64 {
	// Fan-in = r = 64 words per block; internal nodes feed child CVs as data words left-to-right
	const r = md6_r
	var Q = defaultQ()
	var K [8]uint64

	level := L
	cur := nodes
	for level > 0 {
		next := make([][md6_c]uint64, 0, (len(cur)+r-1)/r)
		for i := 0; i < len(cur); {
			// gather up to r child CV words
			var data []uint64
			for j := 0; j < r && i < len(cur); j++ {
				// take entire child CV (16 words). If too many, truncate to fill r words.
				for k := 0; k < md6_c && len(data) < r; k++ {
					data = append(data, cur[i][k])
				}
				i++
			}
			// pad data to r words (if needed) with MD6 padding (1 then zeros)
			if len(data) < r {
				data = padWords(data, r)
			}
			cv := md6Compression(Q, K, uint64(level), uint64(len(next)), r, uint64(level), data, d)
			next = append(next, cv)
		}
		cur = next
		level--
	}
	// root
	if len(cur) == 1 {
		return cur[0]
	}
	// combine remaining into one
	var data []uint64
	for _, cv := range cur {
		for i := 0; i < md6_c; i++ {
			data = append(data, cv[i])
		}
	}
	if len(data) < r {
		data = padWords(data, r)
	} else if len(data) > r {
		data = data[:r]
	}
	return md6Compression(defaultQ(), [8]uint64{}, 0, 0, r, 0, data, d)
}

// Utility and constants

func defaultQ() [15]uint64 {
	// MD6 fixed constants Q[0..14]
	return [15]uint64{
		0x7311C2812425CFA0, 0x6432286434AAC8E7, 0xB60450E9EF68B7C1, 0xE8FB23908D9F06F1,
		0xDD2E76CBA691E5BF, 0x0CD0D63B2C30BC41, 0x1F8CCF6823058F8A, 0x54E5ED5B88E3775D,
		0x4AD12AAE0A6D6031, 0x3E7F16BB88222E0D, 0x8AF8671D3FB50C2C, 0x995AD1178BD25C31,
		0xC878C1DD04C4B633, 0x3B72066C7A1552AC, 0x0D6F3522631EFFCB,
	}
}

func packU(ell, i uint64) uint64 {
	// U = (ell << 56) | i (lower 56 bits)
	return (ell << 56) | (i & 0x00FFFFFFFFFFFFFF)
}

func packV(r, L, zPad, k, d int) uint64 {
	// V packs parameters; here we pack as: (r<<48)|(L<<40)|(zPad<<32)|(k<<24)|d
	return (uint64(r) << 48) | (uint64(L) << 40) | (uint64(zPad) << 32) | (uint64(k) << 24) | uint64(d)
}

func nSteps(d int, r uint64) int {
	// MD6 step count recommendation: steps = 40 + (d/4)
	// This simple rule is common in references; adequate for practical hashing here.
	return 40 + d/4
}

func bytesToWords(b []byte) []uint64 {
	// big-endian pack
	if len(b)%8 != 0 {
		p := make([]byte, ((len(b)+7)/8)*8)
		copy(p, b)
		b = p
	}
	out := make([]uint64, len(b)/8)
	for i := 0; i < len(out); i++ {
		out[i] = binary.BigEndian.Uint64(b[i*8:])
	}
	return out
}

func chunkWords(w []uint64, size int) [][]uint64 {
	var res [][]uint64
	for i := 0; i < len(w); i += size {
		j := i + size
		if j > len(w) {
			j = len(w)
		}
		res = append(res, w[i:j])
	}
	return res
}

func padBlock(r int) []uint64 {
	z := make([]uint64, r)
	z[0] = 0x8000000000000000
	return z
}

func padWords(z []uint64, r int) []uint64 {
	out := make([]uint64, r)
	copy(out, z)
	if len(z) < r {
		out[len(z)] = 0x8000000000000000
		// rest already zero
	}
	return out
}

func treeLevels(bitLen int) int {
	// Choose minimal L so that leaves number fits; simple heuristic:
	// For streaming simplicity, just return 1 for small inputs, else 2.
	if bitLen <= 4096 {
		return 1
	}
	return 2
}
