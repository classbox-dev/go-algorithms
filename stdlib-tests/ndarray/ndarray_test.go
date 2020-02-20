package ndarray

import (
	"fmt"
	ref "hsecode.com/stdlib-tests/internal/ndarray"
	"hsecode.com/stdlib-tests/internal/utils"
	"hsecode.com/stdlib/ndarray"
	"math/rand"
	"strings"
	"testing"
)

type TestCase struct {
	axes, nidx []int
	idx        int
}

func sprintSlice(arr []int) string {
	ss := make([]string, 0, len(arr))
	for _, e := range arr {
		ss = append(ss, fmt.Sprintf("%v", e))
	}
	return strings.Join(ss, ", ")
}

func TestUnit__Basic(t *testing.T) {
	matrix := ndarray.New(3, 3)
	idx := matrix.Idx(1, 1)
	if idx != 4 {
		t.Fatalf("Invalid result: New(3, 3, 3).Idx(1, 1, 1) == %v, expected 4", idx)
	}
}

func TestUnit__ValidIdx(t *testing.T) {
	cases := []TestCase{
		{[]int{64}, []int{0}, 0},
		{[]int{64}, []int{63}, 63},
		{[]int{2, 2}, []int{0, 0}, 0},
		{[]int{2, 2}, []int{1, 1}, 3},
		{[]int{3, 3, 3, 3}, []int{0, 0, 0, 0}, 0},
		{[]int{3, 3, 3, 3}, []int{2, 2, 2, 2}, 80},
	}
	for _, cs := range cases {
		idx := ndarray.New(cs.axes...).Idx(cs.nidx...)
		if idx != cs.idx {
			t.Fatalf("Invalid result: New(%v).Idx(%v) == %v, expected %v", sprintSlice(cs.axes), sprintSlice(cs.nidx), idx, cs.idx)
		}
	}
}

func TestUnit__InvalidNew(t *testing.T) {
	shape := []int{-1, 3}
	msg := fmt.Sprintf("expected panic: New(%v)", sprintSlice(shape))
	utils.ExpectedPanic(t, msg, func() {
		ndarray.New(shape...)
	})
}

func TestUnit__InvalidIdx(t *testing.T) {
	cases := []TestCase{
		{[]int{64}, []int{0, 0}, 0},
		{[]int{64}, []int{0, 0}, 0},
		{[]int{64}, []int{65}, 0},
		{[]int{2, 2}, []int{0, 0, 1}, 0},
		{[]int{2, 2}, []int{1, 1, 3}, 0},
		{[]int{3, 3, 3, 3}, []int{0, 0}, 0},
		{[]int{3, 3, 3, 3}, []int{5, 2, 2, 2}, 0},
	}
	for _, cs := range cases {
		msg := fmt.Sprintf("expected panic: New(%v).Idx(%v)", sprintSlice(cs.axes), sprintSlice(cs.nidx))
		utils.ExpectedPanic(t, msg, func() {
			ndarray.New(cs.axes...).Idx(cs.nidx...)
		})
	}
}

func TestUnit__RandomValidIdx(t *testing.T) {
	utils.InitSeed()
	for i := 0; i < 1000; i++ {
		length := rand.Intn(7) + 1
		var axes, nidx []int
		for j := 0; j < length; j++ {
			a := rand.Intn(8) + 1
			axes = append(axes, a)
			k := rand.Int() % a
			nidx = append(nidx, k)
		}
		expected := ref.New(axes...).Idx(nidx...)
		idx := ndarray.New(axes...).Idx(nidx...)

		if idx != expected {
			t.Fatalf("Invalid result: New(%v).Idx(%v) == %v, expected %v", sprintSlice(axes), sprintSlice(nidx), idx, expected)
		}
	}
}

func TestUnit__RandomInvalidLength(t *testing.T) {
	utils.InitSeed()

	for i := 0; i < 1000; i++ {

		length := rand.Intn(7) + 3
		var axes, nidx []int

		for j := 0; j < length; j++ {
			a := rand.Intn(8) + 1
			axes = append(axes, a)

			k := rand.Int() % a
			nidx = append(nidx, k)
		}

		nidxLong := append(nidx, 10)
		nidxShort := nidxLong[:2]

		msg := fmt.Sprintf("expected panic on New(%v).Idx(%v)", sprintSlice(axes), sprintSlice(nidxLong))
		utils.ExpectedPanic(t, msg, func() {
			ndarray.New(axes...).Idx(nidxLong...)
		})

		msg = fmt.Sprintf("expected panic on New(%v).Idx(%v)", sprintSlice(axes), sprintSlice(nidxShort))
		utils.ExpectedPanic(t, msg, func() {
			ndarray.New(axes...).Idx(nidxShort...)
		})
	}
}

func TestUnit__RandomInvalidValues(t *testing.T) {
	utils.InitSeed()

	for i := 0; i < 1000; i++ {
		length := rand.Intn(7) + 3
		var axes, nidx []int

		for j := 0; j < length; j++ {
			a := rand.Intn(8) + 1
			axes = append(axes, a)
			k := a + rand.Int()%a
			nidx = append(nidx, k)
		}

		msg := fmt.Sprintf("expected panic on New(%v).Idx(%v)", sprintSlice(axes), sprintSlice(nidx))
		utils.ExpectedPanic(t, msg, func() {
			ndarray.New(axes...).Idx(nidx...)
		})
	}
}

func TestPerf__NewOnceIdxMany(t *testing.T) {
	utils.InitSeed()
	matrix := ndarray.New(10, 10, 10, 10, 10, 10, 10, 10, 10, 10)
	var s int
	for i := 0; i < 200000; i++ {
		var nidx []int
		for j := 0; j < 10; j++ {
			k := rand.Intn(10)
			nidx = append(nidx, k)
		}
		idx := matrix.Idx(nidx...)
		s += idx
	}
	utils.Use(s)
}

func TestPerf__NewManyIdxMany(t *testing.T) {
	utils.InitSeed()
	var s int
	for i := 0; i < 200000; i++ {
		matrix := ndarray.New(10, 10, 10, 10, 10, 10, 10, 10, 10, 10)
		var nidx []int
		for j := 0; j < 10; j++ {
			k := rand.Intn(10)
			nidx = append(nidx, k)
		}
		idx := matrix.Idx(nidx...)
		s += idx
	}
	utils.Use(s)
}
