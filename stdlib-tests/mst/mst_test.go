package mst_test

import (
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
	"hsecode.com/stdlib-tests/internal/utils"
	"hsecode.com/stdlib-tests/internal/xgraph"
	"hsecode.com/stdlib-tests/internal/ygraph"
	"hsecode.com/stdlib/graph"
	"hsecode.com/stdlib/graph/mst"
	"math"
	"testing"
)

func TestUnit__Directed(t *testing.T) {
	g := xgraph.RandomConnected(graph.Directed, 5, xgraph.Ordinary, 0)
	msg := "mst.New() did not panic for a directed graph"
	utils.ExpectedPanic(t, msg, func() {
		mst.New(g, func(edge *graph.Edge) int { return 0 })
	})
}

func TestUnit__Random(t *testing.T) {
	for k := 20; k < 150; k += 3 {
		G := xgraph.RandomConnected(graph.Undirected, k, xgraph.Ordinary, 0.2)
		ref := ygraph.ReferenceWeightedUndirected(G)

		b := simple.NewWeightedUndirectedGraph(-1, -1)
		totalExpected := path.Prim(b, ref)

		ng := mst.New(G, func(edge *graph.Edge) int { return edge.Value.(int) })

		total := 0
		ng.Edges(func(u, v *graph.Node, e *graph.Edge) { total += e.Value.(int) })

		if math.Round(float64(total)) != totalExpected {
			t.Fatal("mst.New() returned a graph with non-minimal edge sum")
		}

		tt := ygraph.ReferenceWeightedUndirected(ng)
		if len(topo.ConnectedComponents(tt)) != 1 {
			t.Fatal("mst.New() returned a disconnected graph")
		}
	}
}

func TestPerf__Random(t *testing.T) {
	for k := 10; k < 150; k += 3 {
		G := xgraph.RandomConnected(graph.Undirected, k, xgraph.Ordinary, 0.3)
		for l := 0; l < 50; l++ {
			ng := mst.New(G, func(edge *graph.Edge) int { return edge.Value.(int) })
			utils.Use(&ng)
		}
	}
}
