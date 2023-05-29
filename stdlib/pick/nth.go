package pick

import (
	"golang.org/x/exp/constraints"
	"math/rand"
)

var random = rand.New(rand.NewSource(0xDEADC0DE))

// NthElement transforms the given data so that the element at nth index
// is changed to whatever element would occur in that position if the data were sorted.
// Order of the other elements is not important.
// Panics if nth is out of range.
//
// The behaviour is similar to that of std::nth_element in C++.
func NthElement[E constraints.Ordered](data []E, nth int) {
	n := len(data)
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

func partition[E constraints.Ordered](data []E, lo, hi int) int {
	r := lo + random.Intn(hi-lo)
	data[r], data[hi-1] = data[hi-1], data[r]
	p := hi - 1
	q := lo
	for j := lo; j < p; j++ {
		if data[j] < data[p] {
			data[q], data[j] = data[j], data[q]
			q++
		}
	}
	data[q], data[p] = data[p], data[q]
	return q
}
