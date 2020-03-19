package dijkstra

import (
	"container/heap"
	"hsecode.com/stdlib/graph"
	"hsecode.com/stdlib/graph/internal/xheap"
)

const maxInt = int(^uint(0) >> 1)

// Path represents the shortest path between two nodes in edge-weighted graph
type Path struct {
	// Sequence of nodes on the path (including ends)
	Nodes []*graph.Node
	// Sum of edge weights on the path
	Weight uint
}

func reverse(data []*graph.Node) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

// New computes the shortest path between nodes u and v in the given edge-weighted graph
// with regard to the provided weight function.
// Returns nil if the nodes are not in the graph or there is no path between them.
func New(g *graph.Graph, uid, vid int, weight func(*graph.Edge) uint) *Path {
	an, bn := g.Node(uid), g.Node(vid)
	if an == nil || bn == nil {
		return nil
	}

	hp := &xheap.Nodes{Lookup: make(map[*graph.Node]*xheap.Item, 0)}
	heap.Init(hp)

	g.Nodes(func(node *graph.Node) {
		item := &xheap.Item{Node: node, Parent: nil, Key: maxInt}
		if node == an {
			item.Key = 0
		}
		heap.Push(hp, item)
	})

	for hp.Len() > 0 {
		u := heap.Pop(hp).(*xheap.Item)

		u.Node.Neighbours(func(node *graph.Node, edge *graph.Edge) {
			hi := hp.Lookup[node]
			w := int(weight(edge))
			if hi.Key-w > u.Key {
				hi.Key = u.Key + w
				hi.Parent = u.Node
				heap.Fix(hp, hi.Idx)
			}
		})
	}

	start, finish := hp.Lookup[an], hp.Lookup[bn]
	if finish.Key == maxInt {
		return nil
	}

	nodes := make([]*graph.Node, 0)
	nodes = append(nodes, finish.Node)

	total := uint(0)
	for start != finish {
		parent := finish.Parent
		if parent == nil {
			panic("unknown error")
		}
		nodes = append(nodes, parent)
		if e := parent.Edge(finish.Node); e != nil {
			total += weight(e)
		}
		finish = hp.Lookup[parent]
	}
	reverse(nodes)

	return &Path{nodes, total}
}
