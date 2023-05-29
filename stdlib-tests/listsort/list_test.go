package listsort_test

import (
	"container/list"
	"hsecode.com/stdlib-tests/v2/internal/utils"
	"hsecode.com/stdlib/v2/xlist"
	"sort"
	"testing"
)

const MaxInt = int(^uint(0) >> 1)

func listToSlice(data *list.List) []int {
	output := make([]int, data.Len())
	i := 0
	for e := data.Front(); e != nil; e = e.Next() {
		output[i] = e.Value.(int)
		i++
	}
	return output
}

func sliceToList(data []int) *list.List {
	output := list.New()
	for _, c := range data {
		output.PushBack(c)
	}
	return output
}

func randomData(length int, generator func() int) []int {
	data := make([]int, length)
	for i := 0; i < length; i++ {
		data[i] = generator()
	}
	return data
}

func TestUnit__RandomSmall(t *testing.T) {
	for length := 1; length < 40; length += 1 {
		for bound := length; bound < (MaxInt / 4); bound <<= 2 {
			data := randomData(length, func() int { return utils.Rand.Int() % (2*bound + 1) })

			lst := sliceToList(data)
			xlist.Sort(lst, func(a, b *list.Element) bool { return a.Value.(int) < b.Value.(int) })
			result := listToSlice(lst)

			if !sort.SliceIsSorted(result, func(i, j int) bool { return result[i] < result[j] }) {
				t.Fatalf("sorting error: Sort(%v) == %v", data, result)
			}
		}
	}
}

func TestUnit__RandomCustomLess(t *testing.T) {
	for length := 1; length < 40; length += 1 {
		for bound := length; bound < (MaxInt / 4); bound <<= 2 {
			data := randomData(length, func() int { return utils.Rand.Int() % (2*bound + 1) })

			lst := sliceToList(data)
			xlist.Sort(lst, func(a, b *list.Element) bool { return a.Value.(int) > b.Value.(int) })
			result := listToSlice(lst)

			if !sort.SliceIsSorted(result, func(i, j int) bool { return result[i] > result[j] }) {
				t.Fatalf("reverse sorting error: Sort(%v) == %v", data, result)
			}
		}
	}
}

func TestUnit__RandomLarge(t *testing.T) {
	for lengthBound := 21; lengthBound <= 128*1024; lengthBound <<= 2 {
		for n := 0; n < 3; n++ {
			length := lengthBound + utils.Rand.Intn(lengthBound)

			data := randomData(length, func() int { return utils.Rand.Int() })

			lst := sliceToList(data)
			xlist.Sort(lst, func(a, b *list.Element) bool { return a.Value.(int) < b.Value.(int) })
			result := listToSlice(lst)

			if !sort.SliceIsSorted(result, func(i, j int) bool { return result[i] < result[j] }) {
				t.Fatalf("invalid sort for random list of length %v", length)
			}
		}
	}
}

func TestUnit__InPlace(t *testing.T) {

	data := randomData(5000, func() int { return utils.Rand.Int() })
	lst := sliceToList(data)
	orig := make(map[*list.Element]int)

	for e := lst.Front(); e != nil; e = e.Next() {
		orig[e] = e.Value.(int)
	}

	xlist.Sort(lst, func(a, b *list.Element) bool { return a.Value.(int) < b.Value.(int) })

	for e := lst.Front(); e != nil; e = e.Next() {
		v, ok := orig[e]
		if !ok {
			t.Fatalf("Sort() is not in-place: newly allocated nodes detected in the result")
		}
		if v != e.Value.(int) {
			t.Fatalf("Sort() is not in-place: values of the original elements have been changed")
		}
	}
}

func TestPerf__RandomLarge(t *testing.T) {
	for lengthBound := 21; lengthBound <= 128*1024; lengthBound <<= 2 {
		for n := 0; n < 15; n++ {
			length := lengthBound + utils.Rand.Intn(lengthBound)
			lst := list.New()
			for i := 0; i < length; i++ {
				lst.PushBack(utils.Rand.Int())
			}
			xlist.Sort(lst, func(a, b *list.Element) bool { return a.Value.(int) < b.Value.(int) })
		}
	}
}
