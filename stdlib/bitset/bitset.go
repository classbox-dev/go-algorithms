package bitset

import (
	"errors"
	"math/bits"
)

const bucketSize = 32

// Bitset is fixed-size sequence of bits
type Bitset struct {
	buckets []uint32
	size    int
}

func (b *Bitset) init(size int) *Bitset {
	if size < 1 {
		panic("non-positive bitset size")
	}
	bx := 1 + (size-1)/bucketSize
	buckets := make([]uint32, bx)
	b.buckets = buckets
	b.size = size
	return b
}

// New creates a new Bitset of the given size. All bits of the bitset are initially false.
func New(size int) *Bitset { return new(Bitset).init(size) }

// Set sets a bit value on the provided position. Returns an error if the position is out of range.
func (b *Bitset) Set(pos int, value bool) error {
	if pos < 0 || pos >= b.size {
		return errors.New("out of range")
	}
	idx := pos / bucketSize
	off := pos % bucketSize
	var tv uint32 = 1 << uint(off)
	if value {
		b.buckets[idx] |= tv
	} else {
		b.buckets[idx] &= ^tv
	}
	return nil
}

// Test returns a bit value on the provided position.
// Returns an error if the position is out of range.
func (b *Bitset) Test(pos int) (bool, error) {
	if pos < 0 || pos >= b.size {
		return false, errors.New("out of range")
	}
	idx := pos / bucketSize
	off := pos % bucketSize
	var tv uint32 = 1 << uint(off)
	return (b.buckets[idx] & tv) > 0, nil
}

// Count returns the number of bits set to true
func (b *Bitset) Count() int {
	c := 0
	for _, bucket := range b.buckets {
		c += bits.OnesCount32(bucket)
	}
	return c
}

// All checks if all bits are set to true
func (b *Bitset) All() bool {
	allBits := ^uint32(0)
	for _, bucket := range b.buckets[:len(b.buckets)-1] {
		if bucket != allBits {
			return false
		}
	}
	var allBitsLast uint32 = (1 << uint(b.size%32)) - 1
	return b.buckets[len(b.buckets)-1] == allBitsLast
}

// Any checks if there is at least one bit set to true
func (b *Bitset) Any() bool {
	for _, bucket := range b.buckets {
		if bucket != 0 {
			return true
		}
	}
	return false
}

// Flip changes all bits to the opposite values
func (b *Bitset) Flip() {
	for i, bucket := range b.buckets {
		b.buckets[i] = ^bucket
	}
}

// Reset sets all bits to false
func (b *Bitset) Reset() {
	for i := range b.buckets {
		b.buckets[i] = uint32(0)
	}
}
