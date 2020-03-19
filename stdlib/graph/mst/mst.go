package mst

import (
	"container/heap"
	"hsecode.com/stdlib/graph"
	"hsecode.com/stdlib/graph/internal/xheap"
)

const maxInt = int(^uint(0) >> 1)

// New returns a minimum spanning tree of the given undirected graph with regard to the given weight function.
// Panics if the graph is directed.
func New(g *graph.Graph, weight func(*graph.Edge) int) *graph.Graph {
	if g.Type != graph.Undirected {
		panic("expected undirected graph")
	}

	hp := &xheap.Nodes{Lookup: make(map[*graph.Node]*xheap.Item, 0)}
	heap.Init(hp)

	g.Nodes(func(node *graph.Node) {
		heap.Push(hp, &xheap.Item{Node: node, Parent: nil, Key: maxInt})
	})

	items := hp.Items
	for hp.Len() > 0 {
		u := heap.Pop(hp).(*xheap.Item)

		u.Node.Neighbours(func(node *graph.Node, edge *graph.Edge) {
			hi := hp.Lookup[node]
			if hi.Idx != -1 && weight(edge) < hi.Key {
				hi.Parent = u.Node
				hi.Key = weight(edge)
				heap.Fix(hp, hi.Idx)
			}
		})
	}

	ng := graph.New(graph.Undirected)
	for _, x := range items {
		u := ng.AddNode(x.Node.Value)
		if x.Parent != nil {
			v := ng.AddNode(x.Parent.Value)
			e := x.Parent.Edge(x.Node)
			ng.AddEdge(u, v, e.Value)
		}
	}
	return ng
}
