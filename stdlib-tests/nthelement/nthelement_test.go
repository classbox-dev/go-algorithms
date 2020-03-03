package nthelement_test

import (
	"hsecode.com/stdlib-tests/internal/utils"
	"hsecode.com/stdlib/pick"
	"sort"
	"testing"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func TestUnit__Panic(t *testing.T) {

	msg := "NthElement did not get panic with out-of-range nth"
	utils.ExpectedPanic(t, msg, func() {
		pick.NthElement(sort.IntSlice{1, 2, 3}, -1)
	})
	utils.ExpectedPanic(t, msg, func() {
		pick.NthElement(sort.IntSlice{1, 2, 3}, 5)
	})
}

func TestUnit_Nth(t *testing.T) {
	for size := 2; size <= 65536; size *= 2 {
		data := utils.SliceRandom(1024, size)

		testData := make([]int, size)
		copy(testData, data)

		sorted := make([]int, size)
		copy(sorted, data)
		sort.Ints(sorted)

		for i := 0; i < size; i += max(size/5, 1) {
			pick.NthElement(sort.IntSlice(testData), i)
			if testData[i] != sorted[i] {
				t.Fatalf("NthElement() returned unexpected result (length=%d)", size)
			}
			copy(testData, data)
		}
	}
}

func TestPerf_Nth(t *testing.T) {
	test := func(size int, original, data []int) {
		i := utils.Rand.Intn(size)
		copy(data, original)
		pick.NthElement(sort.IntSlice(data), i)
	}

	for size := 10000; size < 60000; size += 10000 {
		data := sort.IntSlice(make([]int, size))
		test(size, utils.Range(0, size), data)

		bomb := utils.RangeReversed(0, size/2)
		bomb = append(bomb, utils.Range(0, size/2)...)
		test(size, bomb, data)

		test(size, utils.RangeReversed(0, size), data)
		for a := 0; a < 100; a++ {
			test(size, utils.RangeShuffled(0, size), data)
			test(size, utils.SliceRandom(100, size), data)
		}
	}
}
