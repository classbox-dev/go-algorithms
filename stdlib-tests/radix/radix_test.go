package radix_test

import (
	"hsecode.com/stdlib-tests/internal/utils"
	"hsecode.com/stdlib/radix"
	"sort"
	"testing"
)

func randomData(length int, generator func() uint64) []uint64 {
	data := make([]uint64, length)
	for i := 0; i < length; i++ {
		data[i] = utils.Rand.Uint64()
	}
	return data
}

const (
	MaxUint64 = ^uint64(0)
)

func TestUnit__RandomSmall(t *testing.T) {
	for length := 1; length < 40; length += 1 {
		for bound := uint64(length); bound < (MaxUint64 / 4); bound <<= 2 {
			data := randomData(length, func() uint64 { return utils.Rand.Uint64() % (bound + 1) })
			original := make([]uint64, length)
			copy(original, data)
			radix.Sort(data)
			if !sort.SliceIsSorted(data, func(i, j int) bool { return data[i] < data[j] }) {
				t.Fatalf("sorting error: Sort(%v) == %v", original, data)
			}
		}
	}
}

func TestUnit__RandomLarge(t *testing.T) {
	for lengthBound := 21; lengthBound <= 1024*1024; lengthBound <<= 2 {
		for n := 0; n < 3; n++ {
			length := lengthBound + utils.Rand.Intn(lengthBound)
			data := randomData(length, func() uint64 { return utils.Rand.Uint64() })
			radix.Sort(data)
			if !sort.SliceIsSorted(data, func(i, j int) bool { return data[i] < data[j] }) {
				t.Fatalf("invalid sort for random array of length %v", length)
			}
		}
	}
}

func TestPerf__Sort(t *testing.T) {
	for lengthBound := 65536; lengthBound <= 256*1024; lengthBound <<= 2 {
		length := lengthBound + utils.Rand.Intn(lengthBound)
		for j := 0; j < 50; j++ {
			data := make([]uint64, length)
			for i := 0; i < length; i++ {
				data[i] = utils.Rand.Uint64()
			}
			radix.Sort(data)
		}
	}
}
