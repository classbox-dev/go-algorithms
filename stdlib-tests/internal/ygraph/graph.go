package ygraph

import (
	"gonum.org/v1/gonum/graph/simple"
	"hsecode.com/stdlib/graph"
)

func ReferenceWeightedUndirected(g *graph.Graph) *simple.WeightedUndirectedGraph {
	ref := simple.NewWeightedUndirectedGraph(-1, -1)
	lookup := make(map[*graph.Edge]struct{})
	g.Nodes(func(u *graph.Node) {
		u.Neighbours(func(v *graph.Node, edge *graph.Edge) {
			if _, ok := lookup[edge]; ok {
				return
			}
			lookup[edge] = struct{}{}
			e := ref.NewWeightedEdge(simple.Node(u.Value.ID()), simple.Node(v.Value.ID()), float64(edge.Value.(int)))
			ref.SetWeightedEdge(e)
		})
	})
	return ref
}
