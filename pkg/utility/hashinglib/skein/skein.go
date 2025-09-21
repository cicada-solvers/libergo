package skein

import (
	"encoding/binary"
	"hash"
)

// Skein-512 parameters
const (
	skeinBlockBytes = 64 // 512-bit block size
	skeinOutBytes   = 64 // 512-bit output
	threefishWords  = 8  // 8 x 64-bit words
	rounds          = 72
	keyScheduleLen  = threefishWords + 1 + 19 // 9 + 19 = 28 subkeys sets
)

// tweak flags
const (
	tWeakBitFinal = 1 << 63
	tWeakBitFirst = 1 << 62

	tWeakTypeMsg = 48 // message type per spec for message blocks
	tWeakTypeOut = 63 // output stage
)

// Skein512 implements hash.Hash for Skein-512-512.
type Skein512 struct {
	// configuration and chaining state
	chaining [threefishWords]uint64
	t        [2]uint64

	// buffer for partial blocks
	buf    [skeinBlockBytes]byte
	bufLen int

	// total bytes processed in current tree level (into tweak)
	// We maintain in tweak[0] as a 64-bit counter.
}

// NewSkein512 returns a new hash.Hash computing Skein-512-512.
func NewSkein512() hash.Hash {
	var s Skein512
	s.Reset()
	return &s
}

func (s *Skein512) Size() int      { return skeinOutBytes }
func (s *Skein512) BlockSize() int { return skeinBlockBytes }

func (s *Skein512) Reset() {
	// Initialize with configuration block per Skein spec (schema "SHA3")
	// Configuration block is processed as a message block with FIRST|FINAL and type = 4 (CFG).
	for i := range s.chaining {
		s.chaining[i] = 0
	}
	for i := range s.t {
		s.t[i] = 0
	}
	s.bufLen = 0

	// Build configuration block (32 bytes): Schema ID "SHA3", version 1, output length 512 bits.
	var cfg [32]byte
	copy(cfg[:8], []byte("SHA3"))               // schema identifier (ASCII), null-padded
	binary.LittleEndian.PutUint32(cfg[8:12], 1) // version
	binary.LittleEndian.PutUint64(cfg[16:24], 512)

	// Tweak for config block: type = 4 (CFG), FIRST|FINAL set, count = len(cfg)=32
	s.t[0] = 32
	s.t[1] = (1 << 62) | (1 << 63) | uint64(4) // FIRST|FINAL|type=CFG(4)

	var block [skeinBlockBytes]byte
	copy(block[:], cfg[:]) // rest zero
	s.ubi(block[:], true)

	// Prepare for message: reset tweak and buffer
	s.t[0] = 0
	s.t[1] = uint64(tWeakTypeMsg) | (1 << 62) // type=MSG, FIRST set
	s.bufLen = 0
}

func (s *Skein512) Write(p []byte) (int, error) {
	n := len(p)
	if n == 0 {
		return 0, nil
	}

	// If we have buffered bytes, fill to a block first
	if s.bufLen > 0 {
		r := skeinBlockBytes - s.bufLen
		if r > n {
			r = n
		}
		copy(s.buf[s.bufLen:], p[:r])
		s.bufLen += r
		p = p[r:]
		if s.bufLen == skeinBlockBytes {
			s.processMsgBlock(s.buf[:], false)
			s.bufLen = 0
			// after the first processed block, FIRST flag must be cleared
			s.t[1] &^= tWeakBitFirst
		}
	}

	// Process full blocks directly
	for len(p) >= skeinBlockBytes {
		s.processMsgBlock(p[:skeinBlockBytes], false)
		p = p[skeinBlockBytes:]
		s.t[1] &^= tWeakBitFirst
	}

	// Buffer remainder
	if len(p) > 0 {
		copy(s.buf[:], p)
		s.bufLen = len(p)
	}
	return n, nil
}

func (s *Skein512) Sum(in []byte) []byte {
	// Clone state to avoid modifying receiver
	cp := *s

	// Process final message block (with FINAL flag). If buffer is empty, we still must process zero-length final block.
	var finalBlock [skeinBlockBytes]byte
	copy(finalBlock[:], cp.buf[:cp.bufLen])
	// pad with zeros (already zeroed). Set FINAL flag.
	cp.processMsgBlock(finalBlock[:], true)

	// Output stage: generate 64 bytes using UBI with type=OUT over a counter sequence.
	var out [skeinOutBytes]byte
	var ctrBlk [skeinBlockBytes]byte
	// Tweak for OUT: FIRST|FINAL both set for each output block; type = OUT; count = 8 bytes (counter length).
	cp.t[0] = 8
	cp.t[1] = uint64(tWeakTypeOut) | tWeakBitFirst | tWeakBitFinal

	// counter = 0 for first 64 bytes
	binary.LittleEndian.PutUint64(ctrBlk[:8], 0)
	cp.ubi(ctrBlk[:], true)

	// result is chaining value after OUT UBI, little-endian words
	for i := 0; i < threefishWords; i++ {
		binary.LittleEndian.PutUint64(out[i*8:(i+1)*8], cp.chaining[i])
	}
	return append(in, out[:]...)
}

func (s *Skein512) Sum64() [64]byte {
	sum := s.Sum(nil)
	var out [64]byte
	copy(out[:], sum)
	return out
}

// processMsgBlock processes a message block with MSG type.
// If final is true, sets FINAL flag for the block (including if it's a partial block; zero-padded per UBI).
func (s *Skein512) processMsgBlock(block []byte, final bool) {
	if final {
		s.t[1] |= tWeakBitFinal
	}
	// increment processed bytes counter in tweak
	s.t[0] += uint64(len(block))
	s.ubi(block, false)
}

// ubi performs the UBI chaining: G = Threefish(K=G, T=tweak, M=block) XOR M
// where K is current chaining value s.chaining and tweak is s.t.
// For config/out blocks shorter than 64 bytes, block must already be zero-padded.
func (s *Skein512) ubi(block []byte, preserveMsgFlags bool) {
	// Load message as 8 little-endian words (zero-padded if block shorter)
	var m [threefishWords]uint64
	for i := 0; i < threefishWords; i++ {
		start := i * 8
		if start+8 <= len(block) {
			m[i] = binary.LittleEndian.Uint64(block[start : start+8])
		} else if start < len(block) {
			var tmp [8]byte
			copy(tmp[:], block[start:])
			m[i] = binary.LittleEndian.Uint64(tmp[:])
		} else {
			m[i] = 0
		}
	}

	// Run Threefish-512 with key = chaining, tweak = s.t
	out := threefish512Encrypt(s.chaining, s.t, m)

	// XOR with M and set as new chaining value
	for i := 0; i < threefishWords; i++ {
		s.chaining[i] = out[i] ^ m[i]
	}

	// After finishing a UBI block, per UBI: clear FIRST flag; retain FINAL only for the block that had it.
	// The caller controls FINAL; here we only clear FIRST.
	s.t[1] &^= tWeakBitFirst

	// For OUT stage and CFG stage, their own counters/types are set by caller. For MSG, caller maintains t[0] progression.
	// If preserveMsgFlags is false, do nothing else here.
	_ = preserveMsgFlags
}

// threefish512Encrypt: core permutation for Skein-512.
// K: 8-word key; T: 2-word tweak; M: 8-word message block.
// Returns 8-word result.
func threefish512Encrypt(K [8]uint64, T [2]uint64, M [8]uint64) [8]uint64 {
	// Key schedule: k0..k8 with k8 = C240 ^ k0 ^ ... ^ k7
	var ks [9]uint64
	ks[8] = 0x1BD11BDAA9FC1A22
	for i := 0; i < 8; i++ {
		ks[i] = K[i]
		ks[8] ^= K[i]
	}
	// Tweak schedule: t0, t1, t2 where t2 = t0 ^ t1
	var ts [3]uint64
	ts[0], ts[1] = T[0], T[1]
	ts[2] = ts[0] ^ ts[1]

	// Initialize state with plaintext plus first subkey
	var x [8]uint64
	for i := 0; i < 8; i++ {
		x[i] = M[i] + ks[i]
	}
	x[5] += ts[0]
	x[6] += ts[1]
	x[7] += 0

	// Rotation constants per Threefish-512 (Eight mixes per two rounds group)
	rc := [8][4]uint64{
		{46, 36, 19, 37},
		{33, 27, 14, 42},
		{17, 49, 36, 39},
		{44, 9, 54, 56},
		{39, 30, 34, 24},
		{13, 50, 10, 17},
		{25, 29, 39, 43},
		{8, 35, 56, 22},
	}

	// 72 rounds; every 4 rounds inject subkey
	for d := 0; d < rounds; d++ {
		// Mix step (four parallel mixes per round)
		s := rc[d%8]
		// Mix pairs: (0,1), (2,3), (4,5), (6,7)
		x[0] += x[1]
		x[1] = rotl64(x[1], s[0]) ^ x[0]

		x[2] += x[3]
		x[3] = rotl64(x[3], s[1]) ^ x[2]

		x[4] += x[5]
		x[5] = rotl64(x[5], s[2]) ^ x[4]

		x[6] += x[7]
		x[7] = rotl64(x[7], s[3]) ^ x[6]

		// Permutation of words after mixes
		x = [8]uint64{x[0], x[3], x[2], x[1], x[4], x[7], x[6], x[5]}

		// Subkey injection every 4 rounds
		if (d+1)%4 == 0 {
			sk := (d + 1) / 4
			x[0] += ks[(sk+0)%9]
			x[1] += ks[(sk+1)%9]
			x[2] += ks[(sk+2)%9]
			x[3] += ks[(sk+3)%9]
			x[4] += ks[(sk+4)%9]
			x[5] += ks[(sk+5)%9] + ts[sk%3]
			x[6] += ks[(sk+6)%9] + ts[(sk+1)%3]
			x[7] += ks[(sk+7)%9] + uint64(sk)
		}
	}

	return x
}

func rotl64(x uint64, n uint64) uint64 {
	return (x << n) | (x >> (64 - n))
}
