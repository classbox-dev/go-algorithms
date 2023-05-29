package tsort_test

import (
	"fmt"
	"hsecode.com/stdlib/v2/graph"
	"hsecode.com/stdlib/v2/graph/tsort"
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
	G.AddEdge(2, 1, nil)
	G.AddEdge(1, 0, nil)
	G.AddEdge(3, 1, nil)
	G.AddEdge(3, 0, nil)

	nodes, _ := tsort.New(G)
	fmt.Println(nodes)
	// Output: [3 2 1 0]
}
