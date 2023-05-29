package bitset_test

import (
	"fmt"
	"hsecode.com/stdlib-tests/v2/internal/utils"
	"hsecode.com/stdlib/v2/bitset"
	"testing"
)

func must(err error) {
	if err != nil {
		panic(fmt.Sprintf("unexpected error: %v", err))
	}
}

func TestUnit__OutOfRange(t *testing.T) {

	bs := bitset.New(10)

	err := bs.Set(-1, false)
	if err == nil {
		t.Fatalf("Set did not return error for negative position")
	}

	_, err = bs.Test(-1)
	if err == nil {
		t.Fatalf("Test did not return error for negative position")
	}

	err = bs.Set(11, false)
	if err == nil {
		t.Fatalf("Set did not return error for out of range position")
	}

	_, err = bs.Test(11)
	if err == nil {
		t.Fatalf("Test did not return error for out of range position")
	}
}

func TestUnit__All(t *testing.T) {

	bs := bitset.New(1000)

	for i := 0; i < 1000; i++ {
		if bs.All() == true {
			t.Fatalf("All() returned true for non-full bitset")
		}
		must(bs.Set(i, true))
	}
	if bs.All() == false {
		t.Fatalf("All() returned false for full bitset")
	}
}

func TestUnit__Any(t *testing.T) {
	bs := bitset.New(1000)
	if bs.Any() == true {
		t.Fatalf("Any() returned true for empty bitset")
	}
	for i := 0; i < 1000; i++ {
		must(bs.Set(i, true))
		if bs.Any() == false {
			t.Fatalf("Any() returned false for non-empty bitset")
		}
	}
}

func TestUnit__FlipEmpty(t *testing.T) {
	bs := bitset.New(1000)
	bs.Flip()
	for i := 0; i < 1000; i++ {
		if isSet, err := bs.Test(i); err != nil || !isSet {
			t.Fatalf("Test(%d) returned false for flipped empty bitset", i)
		}
	}
}

func TestUnit__FlipNonEmpty(t *testing.T) {
	bs := bitset.New(1000)
	for i := 0; i < 1000; i++ {
		if i%2 == 0 {
			must(bs.Set(i, true))
		}
	}
	bs.Flip()
	for i := 0; i < 1000; i++ {
		if i%2 != 0 {
			if isSet, err := bs.Test(i); err != nil || !isSet {
				t.Fatalf("invalid bit in flipped bitset")
			}
		}
	}
}

func TestUnit__Reset(t *testing.T) {
	bs := bitset.New(1000)
	for i := 0; i < 1000; i++ {
		if i%2 == 0 {
			must(bs.Set(i, true))
		}
	}
	bs.Reset()
	for i := 0; i < 1000; i++ {
		if isSet, err := bs.Test(i); err != nil || isSet {
			t.Fatalf("found non-false bit after Reset")
		}
	}
}

func TestUnit__SetUnsetCount(t *testing.T) {

	bs := bitset.New(2000)

	// Setting bits
	for i := 0; i < 2000; i++ {
		err := bs.Set(i, true)
		if err != nil {
			t.Fatalf("Failed to set a bit: %s", err)
		}
		if isSet, err := bs.Test(i); err != nil || !isSet {
			t.Fatalf("Failed to test a previously set bit")
		}
		if bs.Count() != (i + 1) {
			t.Fatalf("Invalid bit count")
		}
	}

	// Unsetting bits
	for i := 1999; i >= 0; i-- {
		must(bs.Set(i, false))
		if isSet, err := bs.Test(i); err != nil || isSet {
			t.Fatalf("Failed to test a previously unset bit")
		}
		if bs.Count() != i {
			t.Fatalf("Invalid bit count")
		}
	}
}

func TestPerf__SetUnsetCount(t *testing.T) {
	N := 100000
	bs := bitset.New(N)
	// Setting bits
	for i := 0; i < N; i++ {
		must(bs.Set(i, true))
		if _, err := bs.Test(i); err != nil {
			t.Fatalf("Failed to test a bit")
		}
		bs.Count()
	}
	// Unsetting bits
	for i := N - 1; i >= 0; i-- {
		must(bs.Set(i, false))
		if _, err := bs.Test(i); err != nil {
			t.Fatalf("Failed to test a bit")
		}
		bs.Count()
	}
}

func TestPerf__AllAny(t *testing.T) {
	N := 100000
	bs := bitset.New(N)

	for i := 0; i < N; i++ {
		bs.All()
		must(bs.Set(i, true))
		bs.All()
		must(bs.Set(i, false))
	}

	for i := 0; i < N; i++ {
		bs.Any()
		must(bs.Set(i, true))
		bs.Any()
		must(bs.Set(i, false))
	}
}

func TestPerf__ResetFlip(t *testing.T) {
	N := 100000
	bs := bitset.New(N)

	for i := 0; i < N; i++ {
		must(bs.Set(i, true))
		bs.Reset()
	}

	for i := 0; i < N; i++ {
		bs.Flip()
		bs.Flip()
	}
}

func TestPerf__Memory(t *testing.T) {
	n := 65536
	memSize := int64(n / 8)
	var bs *bitset.Bitset
	leak := utils.MemoryLeak(func() {
		bs = bitset.New(n)
	})
	if leak > (memSize * 3 / 2) {
		t.Fatalf("bitset of size %d is expected to take around %d bytes; got %d bytes.", n, memSize, leak)
	}
	v, _ := bs.Test(11)
	utils.Use(v)
}
