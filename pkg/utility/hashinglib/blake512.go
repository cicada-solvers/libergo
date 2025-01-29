package hashinglib

import (
	"encoding/binary"
)

var sigma = [16][16]uint8{
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
	{14, 10, 4, 8, 9, 15, 13, 6, 1, 12, 0, 2, 11, 7, 5, 3},
	{11, 8, 12, 0, 5, 2, 15, 13, 10, 14, 3, 6, 7, 1, 9, 4},
	{7, 9, 3, 1, 13, 12, 11, 14, 2, 6, 5, 10, 4, 0, 15, 8},
	{9, 0, 5, 7, 2, 4, 10, 15, 14, 1, 11, 12, 6, 8, 3, 13},
	{2, 12, 6, 10, 0, 11, 8, 3, 4, 13, 7, 5, 15, 14, 1, 9},
	{12, 5, 1, 15, 14, 13, 4, 10, 0, 7, 6, 3, 9, 2, 8, 11},
	{13, 11, 7, 14, 12, 1, 3, 9, 5, 0, 15, 4, 8, 6, 2, 10},
	{6, 15, 14, 9, 11, 3, 0, 8, 12, 2, 13, 7, 1, 4, 10, 5},
	{10, 2, 8, 4, 7, 6, 1, 5, 15, 11, 9, 14, 3, 12, 13, 0},
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
	{14, 10, 4, 8, 9, 15, 13, 6, 1, 12, 0, 2, 11, 7, 5, 3},
	{11, 8, 12, 0, 5, 2, 15, 13, 10, 14, 3, 6, 7, 1, 9, 4},
	{7, 9, 3, 1, 13, 12, 11, 14, 2, 6, 5, 10, 4, 0, 15, 8},
	{9, 0, 5, 7, 2, 4, 10, 15, 14, 1, 11, 12, 6, 8, 3, 13},
	{2, 12, 6, 10, 0, 11, 8, 3, 4, 13, 7, 5, 15, 14, 1, 9},
}

var u512 = [16]uint64{
	0x243f6a8885a308d3, 0x13198a2e03707344,
	0xa4093822299f31d0, 0x082efa98ec4e6c89,
	0x452821e638d01377, 0xbe5466cf34e90c6c,
	0xc0ac29b7c97c50dd, 0x3f84d5b5b5470917,
	0x9216d5d98979fb1b, 0xd1310ba698dfb5ac,
	0x2ffd72dbd01adfb7, 0xb8e1afed6a267e96,
	0xba7c9045f12c7f99, 0x24a19947b3916cf7,
	0x0801f2e2858efc16, 0x636920d871574e69,
}

var padding = [129]byte{
	0x80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

type Blake512State struct {
	h      [8]uint64
	s      [4]uint64
	t      [2]uint64
	buflen int
	nullt  int
	buf    [128]byte
}

func rot(x uint64, n uint) uint64 {
	return (x << (64 - n)) | (x >> n)
}

func g(v *[16]uint64, m *[16]uint64, a, b, c, d, e int, i int) {
	v[a] += (m[sigma[i][e]] ^ u512[sigma[i][e+1]]) + v[b]
	v[d] = rot(v[d]^v[a], 32)
	v[c] += v[d]
	v[b] = rot(v[b]^v[c], 25)
	v[a] += (m[sigma[i][e+1]] ^ u512[sigma[i][e]]) + v[b]
	v[d] = rot(v[d]^v[a], 16)
	v[c] += v[d]
	v[b] = rot(v[b]^v[c], 11)
}

func blake512Compress(S *Blake512State, block *[128]byte) {
	var v [16]uint64
	var m [16]uint64

	for i := 0; i < 16; i++ {
		m[i] = binary.BigEndian.Uint64(block[i*8:])
	}

	for i := 0; i < 8; i++ {
		v[i] = S.h[i]
	}

	v[8] = S.s[0] ^ u512[0]
	v[9] = S.s[1] ^ u512[1]
	v[10] = S.s[2] ^ u512[2]
	v[11] = S.s[3] ^ u512[3]
	v[12] = u512[4]
	v[13] = u512[5]
	v[14] = u512[6]
	v[15] = u512[7]

	if S.nullt == 0 {
		v[12] ^= S.t[0]
		v[13] ^= S.t[0]
		v[14] ^= S.t[1]
		v[15] ^= S.t[1]
	}

	for i := 0; i < 16; i++ {
		g(&v, &m, 0, 4, 8, 12, 0, i)
		g(&v, &m, 1, 5, 9, 13, 2, i)
		g(&v, &m, 2, 6, 10, 14, 4, i)
		g(&v, &m, 3, 7, 11, 15, 6, i)
		g(&v, &m, 0, 5, 10, 15, 8, i)
		g(&v, &m, 1, 6, 11, 12, 10, i)
		g(&v, &m, 2, 7, 8, 13, 12, i)
		g(&v, &m, 3, 4, 9, 14, 14, i)
	}

	for i := 0; i < 16; i++ {
		S.h[i%8] ^= v[i]
	}

	for i := 0; i < 8; i++ {
		S.h[i] ^= S.s[i%4]
	}
}

func Blake512Init(S *Blake512State) {
	S.h[0] = 0x6a09e667f3bcc908
	S.h[1] = 0xbb67ae8584caa73b
	S.h[2] = 0x3c6ef372fe94f82b
	S.h[3] = 0xa54ff53a5f1d36f1
	S.h[4] = 0x510e527fade682d1
	S.h[5] = 0x9b05688c2b3e6c1f
	S.h[6] = 0x1f83d9abfb41bd6b
	S.h[7] = 0x5be0cd19137e2179
	S.t[0] = 0
	S.t[1] = 0
	S.buflen = 0
	S.nullt = 0
	S.s[0] = 0
	S.s[1] = 0
	S.s[2] = 0
	S.s[3] = 0
}

func Blake512Update(S *Blake512State, in []byte) {
	left := S.buflen
	fill := 128 - left

	if left > 0 && len(in) >= fill {
		copy(S.buf[left:], in[:fill])
		S.t[0] += 1024
		if S.t[0] == 0 {
			S.t[1]++
		}
		blake512Compress(S, &S.buf)
		in = in[fill:]
		left = 0
	}

	for len(in) >= 128 {
		S.t[0] += 1024
		if S.t[0] == 0 {
			S.t[1]++
		}
		blake512Compress(S, (*[128]byte)(in[:128]))
		in = in[128:]
	}

	if len(in) > 0 {
		copy(S.buf[left:], in)
		S.buflen = left + len(in)
	} else {
		S.buflen = 0
	}
}

func Blake512Final(S *Blake512State, out []byte) {
	var msglen [16]byte
	lo := S.t[0] + uint64(S.buflen<<3)
	hi := S.t[1]

	if lo < uint64(S.buflen<<3) {
		hi++
	}

	binary.BigEndian.PutUint64(msglen[0:], hi)
	binary.BigEndian.PutUint64(msglen[8:], lo)

	if S.buflen == 111 {
		S.t[0] -= 8
		Blake512Update(S, []byte{0x81})
	} else {
		if S.buflen < 111 {
			if S.buflen == 0 {
				S.nullt = 1
			}
			S.t[0] -= 888 - uint64(S.buflen<<3)
			Blake512Update(S, padding[:111-S.buflen])
		} else {
			S.t[0] -= 1024 - uint64(S.buflen<<3)
			Blake512Update(S, padding[:128-S.buflen])
			S.t[0] -= 888
			Blake512Update(S, padding[1:112])
			S.nullt = 1
		}
		Blake512Update(S, []byte{0x01})
		S.t[0] -= 8
	}

	S.t[0] -= 128
	Blake512Update(S, msglen[:])

	for i := 0; i < 8; i++ {
		binary.BigEndian.PutUint64(out[i*8:], S.h[i])
	}
}

func Blake512Hash(in []byte) []byte {
	var S Blake512State
	Blake512Init(&S)
	Blake512Update(&S, in)
	out := make([]byte, 64)
	Blake512Final(&S, out)
	return out
}
