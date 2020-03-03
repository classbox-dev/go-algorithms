package firstn_test

import (
	"hsecode.com/stdlib-tests/internal/utils"
	"hsecode.com/stdlib/pick"
	"reflect"
	"sort"
	"testing"
)

func getSorted(a []int, idx []int) []int {
	result := make([]int, len(idx))
	for i := 0; i < len(idx); i++ {
		result[i] = a[idx[i]]
	}
	sort.Ints(result)
	return result
}

type ordered []int

func (p ordered) Len() int           { return len(p) }
func (p ordered) Less(i, j int) bool { return p[i] < p[j] }

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func TestUnit_SmallN(t *testing.T) {
	data := utils.Range(0, 100)
	testData := make([]int, 100)
	copy(testData, data)
	idxs := pick.FirstN(ordered(testData), 105)
	sort.Ints(idxs)
	if !reflect.DeepEqual(idxs, data) {
		t.Fatal("FirstN() returned unexpected result when n >= data.Len()")
	}
}

func TestUnit_FirstN(t *testing.T) {
	for size := 2; size <= 65536; size *= 2 {
		data := utils.SliceRandom(1024, size)

		testData := make([]int, size)
		copy(testData, data)

		sorted := make([]int, size)
		copy(sorted, data)
		sort.Ints(sorted)

		for i := 1; i < size; i += max(size/5, 1) {
			expected := sorted[:i]
			idxs := pick.FirstN(ordered(testData), i)

			if !reflect.DeepEqual(testData, data) {
				t.Fatal("FirstN() cannot change the given data")
			}
			result := getSorted(testData, idxs)

			if !reflect.DeepEqual(expected, result) {
				t.Fatalf("FirstN() returned unexpected result (length=%d)", size)
			}
		}
	}
}

func TestPerf_BigData(t *testing.T) {
	data := utils.RangeShuffled(0, 10*1024*1024)
	for g := 1; g < 30; g++ {
		x := pick.FirstN(ordered(data), g*2)
		utils.Use(x)
	}
}

func TestPerf_FirstN(t *testing.T) {
	for g := 0; g < 6; g++ {
		for size := 2; size <= 65536; size = size * 2 {
			data := utils.SliceRandom(1024, size)
			for i := 1; i < size; i += max(size/10, 1) {
				x := pick.FirstN(ordered(data), i)
				utils.Use(x)
			}
		}
	}
}
