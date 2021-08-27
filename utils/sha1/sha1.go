// Taken from https://pkg.go.dev/crypto/sha1
// I simply removed the assembly-optimized bits.

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sha1 implements the SHA-1 hash algorithm as defined in RFC 3174.
//
// SHA-1 is cryptographically broken and should not be used for secure
// applications.
package sha1

import (
	"crypto"
	"encoding/binary"
	"errors"
	"hash"
)

var block = blockGeneric

func init() {
	crypto.RegisterHash(crypto.SHA1, New)
}

// The size of a SHA-1 checksum in bytes.
const Size = 20

// The blocksize of SHA-1 in bytes.
const BlockSize = 64

const (
	Chunk = 64
	Init0 = 0x67452301
	Init1 = 0xEFCDAB89
	Init2 = 0x98BADCFE
	Init3 = 0x10325476
	Init4 = 0xC3D2E1F0
)

// Digest represents the partial evaluation of a checksum.
type Digest struct {
	h   [5]uint32
	x   [Chunk]byte
	nx  int
	len uint64
}

const (
	magic         = "sha\x01"
	marshaledSize = len(magic) + 5*4 + Chunk + 8
)

func (d *Digest) MarshalBinary() ([]byte, error) {
	b := make([]byte, 0, marshaledSize)
	b = append(b, magic...)
	b = appendUint32(b, d.h[0])
	b = appendUint32(b, d.h[1])
	b = appendUint32(b, d.h[2])
	b = appendUint32(b, d.h[3])
	b = appendUint32(b, d.h[4])
	b = append(b, d.x[:d.nx]...)
	b = b[:len(b)+len(d.x)-int(d.nx)] // already zero
	b = appendUint64(b, d.len)
	return b, nil
}

func (d *Digest) UnmarshalBinary(b []byte) error {
	if len(b) < len(magic) || string(b[:len(magic)]) != magic {
		return errors.New("crypto/sha1: invalid hash state identifier")
	}
	if len(b) != marshaledSize {
		return errors.New("crypto/sha1: invalid hash state size")
	}
	b = b[len(magic):]
	b, d.h[0] = consumeUint32(b)
	b, d.h[1] = consumeUint32(b)
	b, d.h[2] = consumeUint32(b)
	b, d.h[3] = consumeUint32(b)
	b, d.h[4] = consumeUint32(b)
	b = b[copy(d.x[:], b):]
	b, d.len = consumeUint64(b)
	d.nx = int(d.len % Chunk)
	return nil
}

func appendUint64(b []byte, x uint64) []byte {
	var a [8]byte
	binary.BigEndian.PutUint64(a[:], x)
	return append(b, a[:]...)
}

func appendUint32(b []byte, x uint32) []byte {
	var a [4]byte
	binary.BigEndian.PutUint32(a[:], x)
	return append(b, a[:]...)
}

func consumeUint64(b []byte) ([]byte, uint64) {
	_ = b[7]
	x := uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
		uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
	return b[8:], x
}

func consumeUint32(b []byte) ([]byte, uint32) {
	_ = b[3]
	x := uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
	return b[4:], x
}

// Original Reset.
func (d *Digest) Reset() {
	d.h[0] = Init0
	d.h[1] = Init1
	d.h[2] = Init2
	d.h[3] = Init3
	d.h[4] = Init4
	d.nx = 0
	d.len = 0
}

func HashMakeFrom(A, B, C, D, E uint32, l uint64) hash.Hash {
	d := new(Digest)
	d.h[0] = A
	d.h[1] = B
	d.h[2] = C
	d.h[3] = D
	d.h[4] = E
	d.nx = 0
	d.len = l
	return d
}

// TODO.
func MakeFromHash(hash [Size]byte, l uint64) *Digest {
	d := new(Digest)
	d.h[0] = binary.BigEndian.Uint32(hash[:4])
	d.h[1] = binary.BigEndian.Uint32(hash[4:8])
	d.h[2] = binary.BigEndian.Uint32(hash[8:12])
	d.h[3] = binary.BigEndian.Uint32(hash[12:16])
	d.h[4] = binary.BigEndian.Uint32(hash[16:])	
	d.nx = 0
	d.len = l
	return d
}

func MakeFromHashState(A, B, C, D, E uint32, l uint64) *Digest {
	d := new(Digest)
	d.h[0] = A
	d.h[1] = B
	d.h[2] = C
	d.h[3] = D
	d.h[4] = E
	d.nx = 0
	d.len = l
	return d
}

func MakeFromDigest(c *Digest) *Digest {
	d := new(Digest)
	d.h[0] = c.h[0]
	d.h[1] = c.h[1]
	d.h[2] = c.h[2]
	d.h[3] = c.h[3]
	d.h[4] = c.h[4]
	d.nx = c.nx
	d.len = c.len
	return d


}
// New returns a new hash.Hash computing the SHA1 checksum. The Hash also
// implements encoding.BinaryMarshaler and encoding.BinaryUnmarshaler to
// marshal and unmarshal the internal state of the hash.
func New() hash.Hash {
	d := new(Digest)
	d.Reset()
	return d
}

func (d *Digest) Size() int { return Size }

func (d *Digest) BlockSize() int { return BlockSize }

func (d *Digest) Write(p []byte) (nn int, err error) {
	nn = len(p)
	d.len += uint64(nn)
	if d.nx > 0 {
		n := copy(d.x[d.nx:], p)
		d.nx += n
		if d.nx == Chunk {
			block(d, d.x[:])
			d.nx = 0
		}
		p = p[n:]
	}
	if len(p) >= Chunk {
		n := len(p) &^ (Chunk - 1)
		block(d, p[:n])
		p = p[n:]
	}
	if len(p) > 0 {
		d.nx = copy(d.x[:], p)
	}
	return
}

func (d *Digest) Sum(in []byte) []byte {
	// Make a copy of d so that caller can keep writing and summing.
	d0 := *d
	hash := d0.CheckSum()
	return append(in, hash[:]...)
}

func (d *Digest) CheckSum() [Size]byte {
	len := d.len
	// Padding.  Add a 1 bit and 0 bits until 56 bytes mod 64.
	var tmp [64]byte
	tmp[0] = 0x80
	if len%64 < 56 {
		d.Write(tmp[0 : 56-len%64])
	} else {
		d.Write(tmp[0 : 64+56-len%64])
	}

	// Length in bits.
	len <<= 3
	binary.BigEndian.PutUint64(tmp[:], len)
	d.Write(tmp[0:8])

	if d.nx != 0 {
		panic("d.nx != 0")
	}

	var Digest [Size]byte

	binary.BigEndian.PutUint32(Digest[0:], d.h[0])
	binary.BigEndian.PutUint32(Digest[4:], d.h[1])
	binary.BigEndian.PutUint32(Digest[8:], d.h[2])
	binary.BigEndian.PutUint32(Digest[12:], d.h[3])
	binary.BigEndian.PutUint32(Digest[16:], d.h[4])

	return Digest
}

// ConstantTimeSum computes the same result of Sum() but in constant time
func (d *Digest) ConstantTimeSum(in []byte) []byte {
	d0 := *d
	hash := d0.constSum()
	return append(in, hash[:]...)
}

func (d *Digest) constSum() [Size]byte {
	var length [8]byte
	l := d.len << 3
	for i := uint(0); i < 8; i++ {
		length[i] = byte(l >> (56 - 8*i))
	}

	nx := byte(d.nx)
	t := nx - 56                 // if nx < 56 then the MSB of t is one
	mask1b := byte(int8(t) >> 7) // mask1b is 0xFF iff one block is enough

	separator := byte(0x80) // gets reset to 0x00 once used
	for i := byte(0); i < Chunk; i++ {
		mask := byte(int8(i-nx) >> 7) // 0x00 after the end of data

		// if we reached the end of the data, replace with 0x80 or 0x00
		d.x[i] = (^mask & separator) | (mask & d.x[i])

		// zero the separator once used
		separator &= mask

		if i >= 56 {
			// we might have to write the length here if all fit in one block
			d.x[i] |= mask1b & length[i-56]
		}
	}

	// compress, and only keep the Digest if all fit in one block
	block(d, d.x[:])

	var Digest [Size]byte
	for i, s := range d.h {
		Digest[i*4] = mask1b & byte(s>>24)
		Digest[i*4+1] = mask1b & byte(s>>16)
		Digest[i*4+2] = mask1b & byte(s>>8)
		Digest[i*4+3] = mask1b & byte(s)
	}

	for i := byte(0); i < Chunk; i++ {
		// second block, it's always past the end of data, might start with 0x80
		if i < 56 {
			d.x[i] = separator
			separator = 0
		} else {
			d.x[i] = length[i-56]
		}
	}

	// compress, and only keep the Digest if we actually needed the second block
	block(d, d.x[:])

	for i, s := range d.h {
		Digest[i*4] |= ^mask1b & byte(s>>24)
		Digest[i*4+1] |= ^mask1b & byte(s>>16)
		Digest[i*4+2] |= ^mask1b & byte(s>>8)
		Digest[i*4+3] |= ^mask1b & byte(s)
	}

	return Digest
}

// Sum returns the SHA-1 checksum of the data.
func Sum(data []byte) [Size]byte {
	var d Digest
	d.Reset()
	d.Write(data)
	return d.CheckSum()
}


