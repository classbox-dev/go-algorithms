package tsort

import (
	"errors"
	"hsecode.com/stdlib/graph"
	"hsecode.com/stdlib/graph/internal/walks"
	"sort"
)

// New returns one of the possible topological orderings of the provided graph.
// Error is returned if the given graph is undirected or cyclic.
// The function must be implemented non-recursively.
func New(g *graph.Graph) ([]*graph.Node, error) {
	if g.Type == graph.Undirected {
		return []*graph.Node{}, errors.New("directed graph expected")
	}
	dfs := walks.NewDFS(g, func(node *graph.Node) {})
	if dfs.HasCycle {
		return []*graph.Node{}, errors.New("cycle detected")
	}
	nodes := make([]*graph.Node, 0, len(dfs.Nodes))
	for node := range dfs.Nodes {
		nodes = append(nodes, node)
	}

	sort.Slice(nodes, func(i, j int) bool {
		return dfs.Nodes[nodes[i]].Exit > dfs.Nodes[nodes[j]].Exit
	})
	return nodes, nil
}
