package xheap

import "hsecode.com/stdlib/v2/graph"

type Item struct {
	Node, Parent graph.Node
	Key, Idx     int
}

type Nodes struct {
	Items  []*Item
	Lookup map[int]*Item
}

func New() *Nodes {
	ns := new(Nodes)
	ns.Lookup = make(map[int]*Item)
	return ns
}

func (h *Nodes) Len() int           { return len(h.Items) }
func (h *Nodes) Less(i, j int) bool { return h.Items[i].Key < h.Items[j].Key }
func (h *Nodes) Swap(i, j int) {
	h.Items[i].Idx, h.Items[j].Idx = j, i
	h.Items[i], h.Items[j] = h.Items[j], h.Items[i]
}

// Push appends an element to the right.
func (h *Nodes) Push(x interface{}) {
	n := len(h.Items)
	p := x.(*Item)
	p.Idx = n
	h.Lookup[p.Node.ID()] = p
	h.Items = append(h.Items, p)
}

// Pop removes an element from the left.
func (h *Nodes) Pop() interface{} {
	old := h.Items
	n := len(old)
	x := old[n-1]
	x.Idx = -1
	h.Items = old[0 : n-1]
	return x
}
