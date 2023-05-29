package tsort_test

import (
	"hsecode.com/stdlib-tests/v2/internal/utils"
	"hsecode.com/stdlib-tests/v2/internal/xgraph"
	"hsecode.com/stdlib/v2/graph"
	"hsecode.com/stdlib/v2/graph/tsort"
	"runtime/debug"
	"testing"
)

func TestUnit__Undirected(t *testing.T) {
	g := xgraph.RandomConnected(graph.Undirected, 5, xgraph.Ordinary, 0)
	if _, err := tsort.New(g); err == nil {
		t.Fatalf("tsort.New() did not return an error for an undirected graph")
	}
}

func TestUnit__Cyclic(t *testing.T) {
	for i := 0; i < 400; i += 5 {
		n := i + 3
		g := xgraph.RandomConnected(graph.Directed, n, xgraph.ForceCycle, 0)

		if _, err := tsort.New(g); err == nil {
			t.Fatalf("tsort.New() did not return an error for a cyclic graph of %v nodes", n)
		}
	}
}

func TestUnit__NonRecursive(t *testing.T) {
	g := xgraph.RandomConnected(graph.Directed, 1000, xgraph.Ordinary, 0)
	oldSize := debug.SetMaxStack(1024)
	defer debug.SetMaxStack(oldSize)
	ns, _ := tsort.New(g)
	utils.Use(&ns)
}

func TestUnit__Random(t *testing.T) {
	for i := 0; i < 400; i += 5 {
		n := i + 3
		g := xgraph.RandomConnected(graph.Directed, n, xgraph.Ordinary, 0)

		ordering, err := tsort.New(g)
		if err != nil {
			t.Fatal("tsort.New() returned unexpected error")
		}

		uniq := make(map[graph.Node]struct{})
		pos := make(map[graph.Node]int)

		for i, node := range ordering {
			uniq[node] = struct{}{}
			pos[node] = i
		}

		nNodes := 0
		g.Nodes(func(_ graph.Node) {
			nNodes++
		})

		if len(uniq) != nNodes {
			t.Fatal("tsort.New() returned less nodes than expected")
		}

		g.Edges(func(u, v graph.Node, _ interface{}) {
			if pos[u] > pos[v] {
				t.Fatal("tsort.New() returned an invalid topological ordering")
			}
		})

	}
}

func TestPerf__Random(t *testing.T) {
	for i := 3; i < 500; i += 1 {
		g := xgraph.RandomConnected(graph.Directed, i, xgraph.Ordinary, 0)
		for j := 0; j < 10; j++ {
			nodes, _ := tsort.New(g)
			utils.Use(&nodes)
		}
	}
}
