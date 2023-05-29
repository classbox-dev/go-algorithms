package xgraph

import (
	"fmt"
	"hsecode.com/stdlib-tests/v2/internal/utils"
	"hsecode.com/stdlib/v2/graph"
	"sort"
	"strings"
)

type Cycle uint8

const (
	ForceCycle Cycle = iota
	Ordinary
)

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
	nodes := make([]IntValue, 0)
	g.Nodes(func(node graph.Node) { nodes = append(nodes, node.(IntValue)) })
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].ID() < nodes[j].ID()
	})
	g.Edges(func(src, dst graph.Node, data interface{}) {
		output.WriteString(fmt.Sprintf("%v %s-[%v]-> %v\n", src, left, data, dst))
	})
	return output.String()
}

func RandomConnected(type_ graph.Type, n int, cycle Cycle, density float64) *graph.Graph {
	G := graph.New(type_)
	for src := 1; src <= n; src++ {
		G.AddNode(IntValue(src))
		dst := utils.Rand.Intn(src)
		G.AddNode(IntValue(dst))
		weight := utils.Rand.Intn(3 * (n + 1))
		G.AddEdge(src, dst, weight)
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
			G.AddEdge(cn[i], cn[(i+1)%(n/2)], weight)
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
			G.AddNode(IntValue(src))
			G.AddNode(IntValue(dst))
			if _, ok := G.Edge(src, dst); ok {
				i--
				continue
			}
			G.AddEdge(src, dst, utils.Rand.Intn(3*(n+1)))
		}
	}

	return G
}
