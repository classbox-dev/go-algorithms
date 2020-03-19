package dijkstra_test

import (
	"gonum.org/v1/gonum/graph/path"
	"hsecode.com/stdlib-tests/internal/utils"
	"hsecode.com/stdlib-tests/internal/xgraph"
	"hsecode.com/stdlib-tests/internal/ygraph"
	"hsecode.com/stdlib/graph"
	"hsecode.com/stdlib/graph/dijkstra"
	"math"
	"testing"
)

func TestUnit__Random(t *testing.T) {
	for n := 5; n < 100; n += 5 {
		g := xgraph.RandomConnected(graph.Directed, n, xgraph.Ordinary, 0.4)
		ref := ygraph.ReferenceWeightedDirected(g)
		for j := 0; j < 40; j++ {
			a, b := utils.Rand.Intn(n)+1, utils.Rand.Intn(n)+1

			p := dijkstra.New(g, a, b, func(edge *graph.Edge) uint { return uint(edge.Value.(int)) })
			_, correctWeight := path.DijkstraFrom(ref.Node(int64(a)), ref).To(int64(b))

			if p == nil {
				if correctWeight != math.Inf(1) {
					t.Fatal("dijkstra.New() returned nil instead of a possible path")
				}
				continue
			}
			for i := 0; i < len(p.Nodes)-1; i++ {
				if p.Nodes[i].Edge(p.Nodes[i+1]) == nil {
					t.Fatal("dijkstra.New() returned a disconnected path")
				}
			}
			if correctWeight != float64(p.Weight) {
				t.Fatal("dijkstra.New() returned a path with non-optimal edge sum")
			}
		}
	}
}

func TestPerf__Random(t *testing.T) {
	for n := 20; n < 100; n += 1 {
		g := xgraph.RandomConnected(graph.Directed, n, xgraph.Ordinary, 0.4)
		for j := 0; j < 90; j++ {
			a, b := utils.Rand.Intn(n)+1, utils.Rand.Intn(n)+1
			p := dijkstra.New(g, a, b, func(edge *graph.Edge) uint { return uint(edge.Value.(int)) })
			utils.Use(p)
		}
	}
}
