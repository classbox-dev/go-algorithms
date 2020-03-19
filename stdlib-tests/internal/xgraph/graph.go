package xgraph

import (
	"fmt"
	"hsecode.com/stdlib-tests/internal/utils"
	"hsecode.com/stdlib/graph"
	"sort"
	"strings"
)

type Cycle uint8

const (
	ForceCycle Cycle = iota
	Ordinary
)

type sentinel struct{}

type IntValue int

func (v IntValue) ID() int {
	return int(v)
}

func String(g *graph.Graph) string {
	var output strings.Builder
	left := ""
	if g.Type == graph.Undirected {
		left = "<"
	}
	nodes := make([]*graph.Node, 0)
	g.Nodes(func(node *graph.Node) { nodes = append(nodes, node) })
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Value.ID() < nodes[j].Value.ID()
	})
	edges := make(map[*graph.Edge]sentinel, 0)
	g.Nodes(func(u *graph.Node) {
		u.Neighbours(func(v *graph.Node, edge *graph.Edge) {
			if _, ok := edges[edge]; ok {
				return
			}
			edges[edge] = sentinel{}
			output.WriteString(fmt.Sprintf("%v %s-[%v]-> %v\n", u.Value, left, edge.Value, v.Value))
		})
	})
	return output.String()
}

func RandomConnected(type_ graph.Type, n int, cycle Cycle, density float64) *graph.Graph {
	G := graph.New(type_)
	for src := 1; src <= n; src++ {
		u := G.AddNode(IntValue(src))
		dst := utils.Rand.Intn(src)
		v := G.AddNode(IntValue(dst))
		weight := utils.Rand.Intn(3 * (n + 1))
		G.AddEdge(u, v, weight)
	}

	if cycle == ForceCycle {
		cn := make([]int, n)
		for i := range cn {
			cn[i] = i + 1
		}
		utils.Rand.Shuffle(n, func(i, j int) {
			cn[i], cn[j] = cn[j], cn[i]
		})
		cn = cn[:n/2]
		for i := 0; i < n/2; i++ {
			weight := utils.Rand.Intn(3 * (n + 1))
			G.AddEdge(G.Node(cn[i]), G.Node(cn[(i+1)%(n/2)]), weight)
		}
	}

	if density != 0 {
		nEdges := int(density*float64(n)*float64(n-1)) - n - 1
		for i := 0; i < nEdges; i++ {
			src, dst := utils.Rand.Intn(n), utils.Rand.Intn(n)
			if src == dst {
				i--
				continue
			}
			u, v := G.AddNode(IntValue(src)), G.AddNode(IntValue(dst))
			if u.Edge(v) != nil {
				i--
				continue
			}
			G.AddEdge(u, v, utils.Rand.Intn(3*(n+1)))
		}
	}

	return G
}
