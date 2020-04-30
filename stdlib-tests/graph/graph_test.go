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
		g.Nodes(func(node graph.Node) {
			if node == nil {
				t.Fatal("Nodes() yielded nil unexpectedly")
			}
			valuesGraph = append(valuesGraph, node.(Data))
		})

		sort.Slice(values, func(i, j int) bool { return values[i].Id < values[j].Id })
		sort.Slice(valuesGraph, func(i, j int) bool { return valuesGraph[i].Id < valuesGraph[j].Id })

		if !reflect.DeepEqual(values, valuesGraph) {
			t.Fatal("Nodes() yielded unexpected set of nodes")
		}
		for _, value := range values {
			x, ok := g.Node(value.Id)
			if !ok || x.(Data) != value {
				t.Fatal("Node() yielded unexpected node")
			}
		}
		if _, ok := g.Node(-1); ok {
			t.Fatal("Node() indicated success for a missing node")
		}
	}
}

func TestUnit__AddNodeOverwrite(t *testing.T) {
	N := 10
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
	g.Nodes(func(node graph.Node) {
		if node == nil {
			t.Fatal("Nodes() yielded nil unexpectedly")
		}
		valuesGraph = append(valuesGraph, node.(Data))
	})

	sort.Slice(values, func(i, j int) bool { return values[i].Id < values[j].Id })
	sort.Slice(valuesGraph, func(i, j int) bool { return valuesGraph[i].Id < valuesGraph[j].Id })

	if !reflect.DeepEqual(values, valuesGraph) {
		t.Fatal("Nodes() yielded unexpected set of nodes after they have been overwritten by Graph.AddNode()")
	}

	for _, value := range values {
		x, ok := g.Node(value.Id)
		if !ok || x.(Data) != value {
			t.Fatal("Node() yielded unexpected node")
		}
	}
}

func TestUnit__AddEdgePanic(t *testing.T) {
	g1 := graph.New(graph.Directed)
	g1.AddNode(xgraph.IntValue(1))

	g2 := graph.New(graph.Directed)
	g2.AddNode(xgraph.IntValue(2))

	msg := "AddEdge() did not panic when called with missing nodes"
	utils.ExpectedPanic(t, msg, func() {
		g1.AddEdge(1, 2, nil)
	})
	utils.ExpectedPanic(t, msg, func() {
		g1.AddEdge(2, 1, nil)
	})
}

func TestUnit__NeighboursPanic(t *testing.T) {
	g := graph.New(graph.Directed)
	g.AddNode(xgraph.IntValue(1))

	msg := "Neighbours() did not panic when called for a missing node"
	utils.ExpectedPanic(t, msg, func() {
		g.Neighbours(2, func(v graph.Node, _ interface{}) {})
	})
}

func TestUnit__Edges(t *testing.T) {
	N, eN := 10, 30
	g := graph.New(graph.Directed)

	for i := 0; i < N; i++ {
		g.AddNode(xgraph.IntValue(i))
	}

	edges := make(map[int]map[int]float64, 0)

	for j := 0; j < eN; j++ {
		u := utils.Rand.Intn(N)
		v := utils.Rand.Intn(N)
		data := utils.Rand.Float64()
		g.AddEdge(u, v, data)

		if _, ok := edges[u]; !ok {
			edges[u] = make(map[int]float64, 0)
		}
		edges[u][v] = data

	}

	for u, adj := range edges {
		for v, edge := range adj {
			data, ok := g.Edge(u, v)
			if !ok {
				t.Fatal("Edge() did not find an existing edge")
			}
			if data != edge {
				t.Fatal("Edge() returned unexpected edge data")
			}
		}
	}

	// Traverse all edges via Edges
	edgesGraph := make(map[int]map[int]float64, 0)

	g.Edges(func(u, v graph.Node, e interface{}) {
		if _, ok := edgesGraph[u.ID()]; !ok {
			edgesGraph[u.ID()] = make(map[int]float64, 0)
		}
		edgesGraph[u.ID()][v.ID()] = e.(float64)
	})

	if !reflect.DeepEqual(edges, edgesGraph) {
		t.Fatal("Edges() yielded unexpected set of nodes")
	}

	// Traverse all edges via Nodes + Neighbours
	edgesGraph = make(map[int]map[int]float64, 0)
	g.Nodes(func(u graph.Node) {
		g.Neighbours(u.ID(), func(v graph.Node, e interface{}) {
			if _, ok := edgesGraph[u.ID()]; !ok {
				edgesGraph[u.ID()] = make(map[int]float64, 0)
			}
			edgesGraph[u.ID()][v.ID()] = e.(float64)
		})
	})
	if !reflect.DeepEqual(edges, edgesGraph) {
		t.Fatal("Nodes()+Neighbours() yielded unexpected set of nodes")
	}

	g.AddNode(xgraph.IntValue(2 * N))
	if _, ok := g.Edge(0, 2*N); ok {
		t.Fatal("Edge() indicated success for a missing edge")
	}

	if _, ok := g.Edge(2*N+1, 2*N+2); ok {
		t.Fatal("Edge() indicated success for a missing node")
	}
}

func TestUnit__EdgesUndirected(t *testing.T) {
	g := graph.New(graph.Undirected)
	g.AddNode(xgraph.IntValue(42))
	g.AddNode(xgraph.IntValue(43))
	g.AddEdge(42, 43, nil)
	c := 0
	g.Edges(func(u, v graph.Node, e interface{}) {
		c++
	})
	if c != 1 {
		t.Fatalf("Edges() yielded %v edges for a 1-edge undirected graph", c)
	}

	c = 0
	g.Nodes(func(u graph.Node) {
		g.Neighbours(u.ID(), func(v graph.Node, e interface{}) {
			c++
		})
	})
	if c != 2 {
		t.Fatalf("Nodes()+Neighbours() expected to yield twice for a 1-edge undirected graph, got %v", c)
	}
}

func TestPerf__Neighbours(t *testing.T) {
	g := xgraph.RandomConnected(graph.Directed, 1100, xgraph.Ordinary, 0.4)
	g.Nodes(func(u graph.Node) {
		g.Neighbours(u.ID(), func(v graph.Node, edge interface{}) {
			g.Node(u.ID())
			n, _ := g.Node(v.ID())
			utils.Use(n)
			e, _ := g.Edge(u.ID(), v.ID())
			utils.Use(e)
		})
	})
	g = xgraph.RandomConnected(graph.Undirected, 1100, xgraph.Ordinary, 0.4)
	g.Nodes(func(u graph.Node) {
		g.Neighbours(u.ID(), func(v graph.Node, edge interface{}) {
			g.Node(u.ID())
			n, _ := g.Node(v.ID())
			utils.Use(n)
			e, _ := g.Edge(u.ID(), v.ID())
			utils.Use(e)
		})
	})
}
