package tsort

import (
	"errors"
	"hsecode.com/stdlib/graph"
	"hsecode.com/stdlib/graph/internal/walks"
	"sort"
)

// New returns one of the possible topological orderings of the provided directed acyclic graph.
// Returns an error if the given graph is undirected or cyclic.
// The function must be implemented non-recursively.
func New(g *graph.Graph) ([]graph.Node, error) {
	if g.Type == graph.Undirected {
		return nil, errors.New("directed graph expected")
	}

	dfs := walks.NewDFS(g, func(node graph.Node) {})
	if dfs.HasCycle {
		return nil, errors.New("cycle detected")
	}

	nodes := make([]int, 0, len(dfs.Nodes))
	for id := range dfs.Nodes {
		nodes = append(nodes, id)
	}

	sort.Slice(nodes, func(i, j int) bool {
		return dfs.Nodes[nodes[i]].Exit > dfs.Nodes[nodes[j]].Exit
	})

	sorted := make([]graph.Node, len(dfs.Nodes))
	for i, id := range nodes {
		sorted[i], _ = g.Node(id)
	}

	return sorted, nil
}
