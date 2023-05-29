package dijkstra

import (
	"container/heap"
	"hsecode.com/stdlib/v2/graph"
	"hsecode.com/stdlib/v2/graph/internal/xheap"
)

const maxInt = int(^uint(0) >> 1)

// Path represents the shortest path between two nodes in edge-weighted graph
type Path struct {
	// Sequence of nodes on the path (including ends)
	Nodes []graph.Node
	// Sum of edge weights on the path
	Weight uint
}

func reverse(data []graph.Node) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

// New computes the shortest path between nodes u and v in the given edge-weighted graph.
// The given weight function turns edge data into an unsigned integer weight.
// Returns nil if the nodes are not in the graph or there is no path between them.
func New(g *graph.Graph, uid, vid int, weight func(interface{}) uint) *Path {
	_, aOk := g.Node(uid)
	_, bOk := g.Node(vid)

	if !aOk || !bOk {
		return nil
	}

	hp := xheap.New()
	heap.Init(hp)

	g.Nodes(func(node graph.Node) {
		item := &xheap.Item{Node: node, Parent: nil, Key: maxInt}
		if node.ID() == uid {
			item.Key = 0
		}
		heap.Push(hp, item)
	})

	for hp.Len() > 0 {
		u := heap.Pop(hp).(*xheap.Item)

		g.Neighbours(u.Node.ID(), func(node graph.Node, edge interface{}) {
			hi := hp.Lookup[node.ID()]

			w := int(weight(edge))
			if hi.Key-w > u.Key {
				hi.Key = u.Key + w
				hi.Parent = u.Node
				heap.Fix(hp, hi.Idx)
			}
		})
	}

	start, finish := hp.Lookup[uid], hp.Lookup[vid]
	if finish.Key == maxInt {
		return nil
	}

	nodes := make([]graph.Node, 0)
	nodes = append(nodes, finish.Node)

	total := uint(0)

	for start != finish {
		parent := finish.Parent
		if parent == nil {
			panic("unknown error")
		}

		nodes = append(nodes, parent)

		pID := parent.ID()
		if e, ok := g.Edge(pID, finish.Node.ID()); ok {
			total += weight(e)
		}

		finish = hp.Lookup[pID]
	}
	reverse(nodes)

	return &Path{nodes, total}
}
