package pick

import (
	"math/rand"
	"sort"
)

var random = rand.New(rand.NewSource(0xDEADC0DE))

// NthElement transforms the given data so that the element at nth index
// is changed to whatever element would occur in that position if the data were sorted.
// Order of the other elements is not important.
// Panics if nth is out of range.
//
// The behaviour is similar to that of std::nth_element in C++.
func NthElement(data sort.Interface, nth int) {
	n := data.Len()
	if nth < 0 || nth >= n {
		panic("nth is out of range")
	}
	lo, hi := 0, n
	for lo < hi {
		q := partition(data, lo, hi)
		switch {
		case q == nth:
			return
		case nth < q:
			hi = q
		default:
			lo = q + 1
		}
	}
}

func partition(data sort.Interface, lo, hi int) int {
	r := lo + random.Intn(hi-lo)
	data.Swap(r, hi-1)
	p := hi - 1
	q := lo
	for j := lo; j < p; j++ {
		if data.Less(j, p) {
			data.Swap(q, j)
			q++
		}
	}
	data.Swap(q, p)
	return q
}
