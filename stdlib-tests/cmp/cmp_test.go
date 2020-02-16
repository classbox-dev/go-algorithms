package cmp_test

import (
	"hsecode.com/stdlib-tests/internal/utils"
	cmp "hsecode.com/stdlib/cmp/int"
	"testing"
)

func TestUnit__Basic(t *testing.T) {
	v := cmp.Min(2, 3)
	if v != 2 {
		t.Fatalf("Min(2, 3) = %v, expected 2", v)
	}

	v = cmp.Min(5, 3, 1)
	if v != 1 {
		t.Fatalf("Min(5, 3, 1) = %v, expected 1", v)
	}

	v = cmp.Max(2, 3)
	if v != 3 {
		t.Fatalf("Min(2, 3) = %v, expected 3", v)
	}

	v = cmp.Max(5, 3, 1)
	if v != 5 {
		t.Fatalf("Max(5, 3, 1) = %v, expected 5", v)
	}
}

func TestUnit__MinAllPlaces(t *testing.T) {
	s := utils.Range(1, 100)
	for i, x := range s {
		s[i] = 0
		if cmp.Min(s...) != s[i] {
			t.Fatal("Min() could not find expected minimum")
		}
		s[i] = x
	}
}

func TestUnit__MaxAllPlaces(t *testing.T) {
	s := utils.Range(1, 100)
	for i, x := range s {
		s[i] = 200
		if cmp.Max(s...) != s[i] {
			t.Fatal("Max() could not find expected maximum")
		}
		s[i] = x
	}
}

func TestUnit__MinPanic(t *testing.T) {
	panicked := false
	defer func() {
		if r := recover(); r != nil {
			panicked = true
		}
	}()
	cmp.Min()
	if !panicked {
		t.Fatal("Min() with no arguments did not panic")
	}
}

func TestUnit__MaxPanic(t *testing.T) {
	panicked := false
	defer func() {
		if r := recover(); r != nil {
			panicked = true
		}
	}()
	cmp.Max()
	if !panicked {
		t.Fatal("Max() with no arguments did not panic")
	}
}

func TestUnit__Many(t *testing.T) {
	utils.InitSeed()
	for i := 1; i < 500; i += 10 {
		s := utils.RangeShuffled(i, i+i)
		if cmp.Min(s...) != i {
			t.Fatalf("Min() did not return expected value out of %d arguments", i)
		}
		if cmp.Max(s...) != i+i-1 {
			t.Fatalf("Max() did not return expected value out of %d arguments", i)
		}
	}
}

func TestPerf__Many(t *testing.T) {
	utils.InitSeed()
	for c := 0; c < 350; c++ {
		for i := 2; i < 500; i += 1 {
			s := utils.RangeShuffled(i, i+i)
			cmp.Min(s...)
			cmp.Max(s...)
		}
	}
}
