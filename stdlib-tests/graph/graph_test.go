package graph_test

import (
	"hsecode.com/stdlib-tests/internal/utils"
	"hsecode.com/stdlib-tests/internal/xgraph"
	"hsecode.com/stdlib/graph"
	"reflect"
	"sort"
	"testing"
)

type Data struct {
	Id, Payload int
}

func (v Data) ID() int {
	return v.Id
}

func TestUnit__AddNode(t *testing.T) {
	N := 100
	g := graph.New(graph.Undirected)
	values := make([]Data, 0, N)

	for i := 0; i < N; i++ {
		data := Data{Id: utils.Rand.Int(), Payload: utils.Rand.Int()}

		values = append(values, data)
		g.AddNode(data)

		valuesGraph := make([]Data, 0, i+1)
		g.Nodes(func(node *graph.Node) {
			if node == nil {
				t.Fatal("Graph.Nodes() yielded nil unexpectedly")
			}
			valuesGraph = append(valuesGraph, node.Value.(Data))
		})

		sort.Slice(values, func(i, j int) bool { return values[i].Id < values[j].Id })
		sort.Slice(valuesGraph, func(i, j int) bool { return valuesGraph[i].Id < valuesGraph[j].Id })

		if !reflect.DeepEqual(values, valuesGraph) {
			t.Fatal("Graph.Nodes() yielded unexpected set of nodes")
		}
		for _, value := range values {
			x := g.Node(value.Id)
			if x == nil || x.Value.(Data) != value {
				t.Fatal("Graph.Node() yielded unexpected node")
			}
		}
		if g.Node(-1) != nil {
			t.Fatal("Graph.Node() unexpectedly returned a node instead of nil")
		}
	}
}

func TestUnit__AddNodeOverwrite(t *testing.T) {
	N := 100
	g := graph.New(graph.Undirected)

	values := make([]Data, 0, N)
	for i := 0; i < N; i++ {
		data := Data{Id: utils.Rand.Int(), Payload: utils.Rand.Int()}
		values = append(values, data)
		g.AddNode(data)
	}
	for i := 0; i < N; i++ {
		values[i].Payload = utils.Rand.Int()
		g.AddNode(values[i])
	}

	valuesGraph := make([]Data, 0, N)
	g.Nodes(func(node *graph.Node) {
		if node == nil {
			t.Fatal("Graph.Nodes() yielded nil unexpectedly")
		}
		valuesGraph = append(valuesGraph, node.Value.(Data))
	})

	sort.Slice(values, func(i, j int) bool { return values[i].Id < values[j].Id })
	sort.Slice(valuesGraph, func(i, j int) bool { return valuesGraph[i].Id < valuesGraph[j].Id })

	if !reflect.DeepEqual(values, valuesGraph) {
		t.Fatal("Graph.Nodes() yielded unexpected set of nodes after they have been overwritten by Graph.AddNode()")
	}

	for _, value := range values {
		x := g.Node(value.Id)
		if x == nil || x.Value.(Data) != value {
			t.Fatal("Graph.Node() yielded unexpected node")
		}
	}
}

func TestUnit__AddEdgePanic(t *testing.T) {
	g1 := graph.New(graph.Directed)
	n1 := g1.AddNode(xgraph.IntValue(1))

	g2 := graph.New(graph.Directed)
	n2 := g2.AddNode(xgraph.IntValue(2))

	msg := "Graph.AddEdge() did not panic when called with non-existing nodes"
	utils.ExpectedPanic(t, msg, func() {
		g1.AddEdge(n1, n2, nil)
	})
	utils.ExpectedPanic(t, msg, func() {
		g1.AddEdge(n2, n1, nil)
	})

	msg = "Graph.AddEdge() did not panic when called with a nil node"
	utils.ExpectedPanic(t, msg, func() {
		var e *graph.Node
		g1.AddEdge(n1, e, nil)
	})
}

func TestUnit__Edges(t *testing.T) {
	N, eN := 10, 30
	g := graph.New(graph.Directed)

	for i := 0; i < N; i++ {
		g.AddNode(xgraph.IntValue(i))
	}

	edges := make(map[*graph.Node]map[*graph.Node]*graph.Edge, 0)

	for j := 0; j < eN; j++ {
		u := g.Node(utils.Rand.Intn(N))
		v := g.Node(utils.Rand.Intn(N))
		e := g.AddEdge(u, v, nil)
		if _, ok := edges[u]; !ok {
			edges[u] = make(map[*graph.Node]*graph.Edge, 0)
		}
		edges[u][v] = e
	}

	for u, adj := range edges {
		for v, edge := range adj {
			if u.Edge(v) != edge {
				t.Fatal("Node.Edge() yielded unexpected edge")
			}
		}
	}

	// Traverse all edges via Edges
	edgesGraph := make(map[*graph.Node]map[*graph.Node]*graph.Edge, 0)
	g.Edges(func(u, v *graph.Node, e *graph.Edge) {
		if _, ok := edgesGraph[u]; !ok {
			edgesGraph[u] = make(map[*graph.Node]*graph.Edge, 0)
		}
		edgesGraph[u][v] = e
	})

	if !reflect.DeepEqual(edges, edgesGraph) {
		t.Fatal("Graph.Edges() yielded unexpected set of nodes")
	}

	// Traverse all edges via Nodes + Neighbours
	edgesGraph = make(map[*graph.Node]map[*graph.Node]*graph.Edge, 0)
	g.Nodes(func(u *graph.Node) {
		u.Neighbours(func(v *graph.Node, e *graph.Edge) {
			if _, ok := edgesGraph[u]; !ok {
				edgesGraph[u] = make(map[*graph.Node]*graph.Edge, 0)
			}
			edgesGraph[u][v] = e
		})
	})

	if !reflect.DeepEqual(edges, edgesGraph) {
		t.Fatal("Graph.Nodes() + Node.Neighbours() yielded unexpected set of nodes")
	}

	extra := g.AddNode(xgraph.IntValue(2 * N))
	if g.Node(0).Edge(extra) != nil {
		t.Fatal("Node.Edge() unexpectedly returned a node instead of nil")
	}
}

func TestUnit__EdgesUndirected(t *testing.T) {
	g := graph.New(graph.Undirected)
	a, b := g.AddNode(xgraph.IntValue(42)), g.AddNode(xgraph.IntValue(43))
	g.AddEdge(a, b, nil)
	c := 0
	g.Edges(func(u, v *graph.Node, e *graph.Edge) {
		c++
	})
	if c != 1 {
		t.Fatalf("Graph.Edges() yielded %v edges for a 1-edge undirected graph", c)
	}

	c = 0
	g.Nodes(func(u *graph.Node) {
		u.Neighbours(func(v *graph.Node, e *graph.Edge) {
			c++
		})
	})
	if c != 2 {
		t.Fatalf("Graph.Nodes() + Node.Neighbours() expected to yield twice for a 1-edge undirected graph, got %v", c)
	}
}

func TestPerf__Neighbours(t *testing.T) {
	g := xgraph.RandomConnected(graph.Directed, 1100, xgraph.Ordinary, 0.4)
	g.Nodes(func(u *graph.Node) {
		u.Neighbours(func(v *graph.Node, edge *graph.Edge) {
			utils.Use(g.Node(u.Value.ID()))
			utils.Use(g.Node(v.Value.ID()))
			utils.Use(u.Edge(v))
		})
	})
	g = xgraph.RandomConnected(graph.Undirected, 1100, xgraph.Ordinary, 0.4)
	g.Nodes(func(u *graph.Node) {
		u.Neighbours(func(v *graph.Node, edge *graph.Edge) {
			utils.Use(g.Node(u.Value.ID()))
			utils.Use(g.Node(v.Value.ID()))
			utils.Use(u.Edge(v))
		})
	})
}
