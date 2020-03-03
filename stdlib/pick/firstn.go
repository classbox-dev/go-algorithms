package pick

import "container/heap"

type dataHeap struct {
	heap []int
	data Ordered
}

func (h dataHeap) Len() int           { return len(h.heap) }
func (h dataHeap) Less(i, j int) bool { return h.data.Less(h.heap[j], h.heap[i]) }
func (h dataHeap) Swap(i, j int)      { h.heap[i], h.heap[j] = h.heap[j], h.heap[i] }
func (h dataHeap) Push(x interface{}) {
	h.heap = append(h.heap, x.(int))
}
func (h dataHeap) Pop() interface{} {
	n := len(h.heap)
	x := h.heap[n-1]
	h.heap = h.heap[0 : n-1]
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// The type describes an indexed sequence (typically, slice) of comparable elements.
type Ordered interface {
	// Len returns the number of elements in the collection.
	Len() int
	// Less reports whether the element with index i should sort before the element with index j.
	Less(i, j int) bool
}

// FirstN returns min(n, data.Len()) indices of elements that would occur in data[:n] if the data was sorted.
// The indices are returned in arbitrary order.
func FirstN(data Ordered, n int) []int {
	k := min(n, data.Len())
	hp := make([]int, k)
	for i := 0; i < k; i++ {
		hp[i] = i
	}
	if data.Len() < n {
		return hp
	}
	dh := dataHeap{hp, data}
	heap.Init(dh)
	for i := n; i < data.Len(); i++ {
		if dh.data.Less(i, dh.heap[0]) {
			dh.heap[0] = i
			heap.Fix(dh, 0)
		}
	}
	return dh.heap
}
