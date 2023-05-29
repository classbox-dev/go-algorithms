package dijkstra_test

import (
	"gonum.org/v1/gonum/graph/path"
	"hsecode.com/stdlib-tests/v2/internal/utils"
	"hsecode.com/stdlib-tests/v2/internal/xgraph"
	"hsecode.com/stdlib-tests/v2/internal/ygraph"
	"hsecode.com/stdlib/v2/graph"
	"hsecode.com/stdlib/v2/graph/dijkstra"
	"math"
	"testing"
)

func TestUnit__Random(t *testing.T) {
	for n := 5; n < 100; n += 5 {
		g := xgraph.RandomConnected(graph.Directed, n, xgraph.Ordinary, 0.4)
		ref := ygraph.ReferenceWeightedDirected(g)
		for j := 0; j < 40; j++ {
			a, b := utils.Rand.Intn(n)+1, utils.Rand.Intn(n)+1

			p := dijkstra.New(g, a, b, func(edge interface{}) uint { return uint(edge.(int)) })
			_, correctWeight := path.DijkstraFrom(ref.Node(int64(a)), ref).To(int64(b))

			if p == nil {
				if correctWeight != math.Inf(1) {
					t.Fatal("dijkstra.New() returned nil instead of a possible path")
				}
				continue
			}

			for i := 0; i < len(p.Nodes)-1; i++ {
				if _, ok := g.Edge(p.Nodes[i].ID(), p.Nodes[i+1].ID()); !ok {
					t.Fatal("dijkstra.New() returned a disconnected path")
				}
			}

			if math.Round(correctWeight) != math.Round(float64(p.Weight)) {
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
			p := dijkstra.New(g, a, b, func(edge interface{}) uint { return uint(edge.(int)) })
			utils.Use(p)
		}
	}
}
