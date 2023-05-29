package graph_test

import (
	"fmt"
	"sort"

	"hsecode.com/stdlib/v2/graph"
)

type Int int

func (v Int) ID() int {
	return int(v)
}

func Example() {
	G := graph.New(graph.Undirected)

	// Create the following graph with no edge data:
	//
	//  1---2 \
	//  | / |  3
	//  5---4 /

	for i := 1; i <= 5; i++ {
		G.AddNode(Int(i))
	}
	G.AddEdge(1, 2, nil)
	G.AddEdge(2, 3, nil)
	G.AddEdge(3, 4, nil)
	G.AddEdge(4, 5, nil)
	G.AddEdge(5, 1, nil)
	G.AddEdge(2, 4, nil)
	G.AddEdge(2, 5, nil)

	// Compute neighbour nodes for 2
	nodes := make([]int, 0)
	G.Neighbours(2, func(v graph.Node, e interface{}) {
		nodes = append(nodes, v.ID())
	})

	sort.Ints(nodes)
	fmt.Println(nodes)
	// Output: [1 3 4 5]
}
