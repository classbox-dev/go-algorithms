package mst_test

import (
	"fmt"
	"hsecode.com/stdlib/v2/graph"
	"hsecode.com/stdlib/v2/graph/mst"
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
	G.AddEdge(0, 1, 50)
	G.AddEdge(0, 2, 50)
	G.AddEdge(0, 3, 30)

	// "Light" nodes
	G.AddEdge(1, 2, 2)
	G.AddEdge(2, 3, 5)
	G.AddEdge(3, 1, 8)

	tree := mst.New(G, func(edge interface{}) int { return edge.(int) })

	tree.Edges(func(u, v graph.Node, _ interface{}) {
		fmt.Println(u.ID(), v.ID())
	})

	// Output:
	// 0 3
	// 3 2
	// 2 1
}
