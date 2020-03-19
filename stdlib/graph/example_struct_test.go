package graph_test

import (
	"fmt"
	"hsecode.com/stdlib/graph"
)

type Person struct {
	Id   int
	Name string
}

func (d Person) ID() int { return d.Id }

func ExampleGraph_AddNode_overwrite() {
	G := graph.New(graph.Undirected)

	G.AddNode(Person{Id: 1, Name: "Fred Weasley"})
	fmt.Println(G.Node(1).Value.(Person).Name)

	// data will be overwritten due to the same ID
	G.AddNode(Person{Id: 1, Name: "George Weasley"})
	fmt.Println(G.Node(1).Value.(Person).Name)

	// Output:
	// Fred Weasley
	// George Weasley
}
