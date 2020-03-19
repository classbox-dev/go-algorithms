package ygraph

import (
	"gonum.org/v1/gonum/graph/simple"
	"hsecode.com/stdlib/graph"
)

func ReferenceWeightedUndirected(g *graph.Graph) *simple.WeightedUndirectedGraph {
	ref := simple.NewWeightedUndirectedGraph(-1, -1)
	g.Nodes(func(node *graph.Node) {
		ref.AddNode(simple.Node(node.Value.ID()))
	})
	g.Edges(func(u, v *graph.Node, e *graph.Edge) {
		edge := ref.NewWeightedEdge(simple.Node(u.Value.ID()), simple.Node(v.Value.ID()), float64(e.Value.(int)))
		ref.SetWeightedEdge(edge)
	})
	return ref
}
