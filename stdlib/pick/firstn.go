package pick

import (
	"container/heap"
	"golang.org/x/exp/constraints"
)

type dataHeap[E constraints.Ordered] struct {
	heap []int
	data []E
}

func (h *dataHeap[E]) Len() int           { return len(h.heap) }
func (h *dataHeap[E]) Less(i, j int) bool { return h.data[h.heap[j]] < h.data[h.heap[i]] }
func (h *dataHeap[E]) Swap(i, j int)      { h.heap[i], h.heap[j] = h.heap[j], h.heap[i] }
func (h *dataHeap[E]) Push(x interface{}) {
	h.heap = append(h.heap, x.(int))
}
func (h *dataHeap[E]) Pop() interface{} {
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

// FirstN returns min(n, data.Len()) indices of elements that would occur in data[:n] if the data was sorted.
// The indices are returned in arbitrary order.
func FirstN[E constraints.Ordered](data []E, n int) []int {
	k := min(n, len(data))
	hp := make([]int, k)
	for i := 0; i < k; i++ {
		hp[i] = i
	}
	if len(data) < n {
		return hp
	}
	dh := &dataHeap[E]{hp, data}
	heap.Init(dh)
	for i := n; i < len(data); i++ {
		if dh.data[i] < dh.data[dh.heap[0]] {
			dh.heap[0] = i
			heap.Fix(dh, 0)
		}
	}
	return dh.heap
}
