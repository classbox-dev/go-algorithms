package mst

import (
	"container/heap"
	"hsecode.com/stdlib/graph"
	"hsecode.com/stdlib/graph/internal/xheap"
)

const maxInt = int(^uint(0) >> 1)

// New returns a minimum spanning tree of the given undirected graph subject to edge weights.
// The given weight function turns edge data into an integer weight.
// Panics if the graph is directed.
func New(g *graph.Graph, weight func(interface{}) int) *graph.Graph {
	if g.Type != graph.Undirected {
		panic("expected undirected graph")
	}

	hp := xheap.New()
	heap.Init(hp)

	g.Nodes(func(node graph.Node) {
		heap.Push(hp, &xheap.Item{Node: node, Parent: nil, Key: maxInt})
	})

	items := hp.Items
	for hp.Len() > 0 {
		u := heap.Pop(hp).(*xheap.Item)

		g.Neighbours(u.Node.ID(), func(node graph.Node, edge interface{}) {
			hi := hp.Lookup[node.ID()]

			if hi.Idx != -1 && weight(edge) < hi.Key {
				hi.Parent = u.Node
				hi.Key = weight(edge)
				heap.Fix(hp, hi.Idx)
			}
		})
	}

	ng := graph.New(graph.Undirected)

	for _, x := range items {
		xID := x.Node.ID()
		ng.AddNode(x.Node)
		if x.Parent != nil {
			pID := x.Parent.ID()
			ng.AddNode(x.Parent)
			e, _ := g.Edge(pID, xID)
			ng.AddEdge(pID, xID, e)
		}
	}

	return ng
}
