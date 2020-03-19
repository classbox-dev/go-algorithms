package graph_test

import (
	"fmt"
	"hsecode.com/stdlib/graph"
	"sort"
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
	G.AddEdge(G.Node(1), G.Node(2), nil)
	G.AddEdge(G.Node(2), G.Node(3), nil)
	G.AddEdge(G.Node(3), G.Node(4), nil)
	G.AddEdge(G.Node(4), G.Node(5), nil)
	G.AddEdge(G.Node(5), G.Node(1), nil)
	G.AddEdge(G.Node(2), G.Node(4), nil)
	G.AddEdge(G.Node(2), G.Node(5), nil)

	// Compute neighbour nodes for 2
	nodes := make([]int, 0)
	G.Node(2).Neighbours(func(v *graph.Node, e *graph.Edge) {
		nodes = append(nodes, v.Value.ID())
	})

	sort.Ints(nodes)
	fmt.Println(nodes)
	// Output: [1 3 4 5]
}

func ExampleGraph_Edges_undirected() {
	ug := graph.New(graph.Undirected)
	ug.AddEdge(ug.AddNode(Int(2)), ug.AddNode(Int(3)), "single edge")
	ug.Edges(func(u, v *graph.Node, e *graph.Edge) {
		fmt.Println(u.Value, v.Value, e.Value)
	})
	// Output: 2 3 single edge
}

func ExampleNode_Edge() {
	g := graph.New(graph.Directed)

	u, v := g.AddNode(Int(2)), g.AddNode(Int(3))
	g.AddEdge(u, v, "some edge")

	fmt.Println(u.Edge(v).Value)
	fmt.Println(v.Edge(u))
	// Output:
	// some edge
	// <nil>
}

func ExampleGraph_Node() {
	g := graph.New(graph.Directed)
	g.AddNode(Int(2))

	fmt.Println(g.Node(2).Value)
	fmt.Println(g.Node(8))
	// Output:
	// 2
	// <nil>
}

func ExampleGraph_Nodes() {
	g := graph.New(graph.Directed)
	g.AddNode(Int(2))
	g.AddNode(Int(3))

	g.Nodes(func(node *graph.Node) {
		fmt.Println(node.Value)
	})
	// Output:
	// 2
	// 3
}

func ExampleGraph_Edges_directed() {
	g := graph.New(graph.Directed)
	u, v := g.AddNode(Int(2)), g.AddNode(Int(3))
	g.AddEdge(u, v, "forward edge")
	g.AddEdge(v, u, "backward edge")

	g.Edges(func(u, v *graph.Node, e *graph.Edge) {
		fmt.Println(u.Value, v.Value, e.Value)
	})
	// Output:
	// 2 3 forward edge
	// 3 2 backward edge
}
