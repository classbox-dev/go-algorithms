package utils

import (
	"math/rand"
	"runtime"
	"runtime/debug"
	"sort"
	"testing"
)

func SliceRandom(rng int, length int) []int {
	output := make([]int, 0, length)
	for i := 0; i < length; i++ {
		output = append(output, (rand.Int()%(2*rng))-rng)
	}
	return output
}

func InitSeed() {
	rand.Seed(0xDEADBEEF)
}

func RangeShuffled(a, b int) []int {
	s := Range(a, b)
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
	return s
}

func RangeReversed(a, b int) []int {
	s := Range(a, b)
	sort.Slice(s, func(i, j int) bool { return i > j })
	return s
}

func Range(a, b int) []int {
	s := make([]int, 0, b-a)
	for i := a; i < b; i++ {
		s = append(s, i)
	}
	return s
}

func ExpectedPanic(t *testing.T, message string, f func()) {
	panicked := false
	defer func() {
		if r := recover(); r != nil {
			panicked = true
		}
	}()
	f()
	if !panicked {
		t.Fatal(message)
	}
}

func MemoryLeak(f func()) int64 {
	debug.SetGCPercent(-1) // Disable GC
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memStart := int64(m.Alloc)

	f()

	runtime.GC()
	runtime.ReadMemStats(&m)
	return int64(m.Alloc) - memStart
}
