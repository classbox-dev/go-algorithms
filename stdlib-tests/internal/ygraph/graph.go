package ygraph

import (
	ggraph "gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"hsecode.com/stdlib/graph"
)

func ReferenceWeightedUndirected(g *graph.Graph) *simple.WeightedUndirectedGraph {
	ref := simple.NewWeightedUndirectedGraph(-1, -1)
	copyGraph(ref, g)
	return ref
}

func ReferenceWeightedDirected(g *graph.Graph) *simple.WeightedDirectedGraph {
	ref := simple.NewWeightedDirectedGraph(-1, -1)
	copyGraph(ref, g)
	return ref
}

func copyGraph(dst ggraph.WeightedBuilder, src *graph.Graph) {
	src.Nodes(func(node *graph.Node) {
		dst.AddNode(simple.Node(node.Value.ID()))
	})
	src.Edges(func(u, v *graph.Node, e *graph.Edge) {
		edge := dst.NewWeightedEdge(simple.Node(u.Value.ID()), simple.Node(v.Value.ID()), float64(e.Value.(int)))
		dst.SetWeightedEdge(edge)
	})
}
