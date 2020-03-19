package tsort_test

import (
	"fmt"
	"hsecode.com/stdlib/graph"
	"hsecode.com/stdlib/graph/tsort"
)

type Int int

func (v Int) ID() int {
	return int(v)
}

func Example() {
	G := graph.New(graph.Directed)

	for i := 0; i <= 3; i++ {
		G.AddNode(Int(i))
	}
	G.AddEdge(G.Node(2), G.Node(1), nil)
	G.AddEdge(G.Node(1), G.Node(0), nil)
	G.AddEdge(G.Node(3), G.Node(1), nil)
	G.AddEdge(G.Node(3), G.Node(0), nil)

	nodes, _ := tsort.New(G)
	for _, node := range nodes {
		fmt.Println(node.Value)
	}
	// Output:
	// 3
	// 2
	// 1
	// 0
}
