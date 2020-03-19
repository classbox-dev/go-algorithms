package mst_test

import (
	"fmt"
	"hsecode.com/stdlib/graph"
	"hsecode.com/stdlib/graph/mst"
)

type Int int

func (v Int) ID() int {
	return int(v)
}

func Example() {
	G := graph.New(graph.Undirected)

	for i := 0; i <= 3; i++ {
		G.AddNode(Int(i))
	}

	// "Heavy" nodes
	G.AddEdge(G.Node(0), G.Node(1), 50)
	G.AddEdge(G.Node(0), G.Node(2), 50)
	G.AddEdge(G.Node(0), G.Node(3), 30)

	// "Light" nodes
	G.AddEdge(G.Node(1), G.Node(2), 2)
	G.AddEdge(G.Node(2), G.Node(3), 5)
	G.AddEdge(G.Node(3), G.Node(1), 8)

	tree := mst.New(G, func(edge *graph.Edge) int { return edge.Value.(int) })

	tree.Edges(func(u, v *graph.Node, e *graph.Edge) {
		fmt.Println(u.Value, v.Value)
	})

	// Output:
	// 0 3
	// 3 2
	// 2 1
}
