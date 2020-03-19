package walks

import (
	"hsecode.com/stdlib/graph"
)

type Colour uint8

const (
	White Colour = iota
	Grey
	Black
)

type DFS struct {
	Nodes    map[*graph.Node]*DFSItem
	HasCycle bool
}

type DFSItem struct {
	Colour Colour
	Enter  int
	Exit   int
}

func NewDFS(g *graph.Graph, f func(*graph.Node)) *DFS {
	dfs := new(DFS)
	dfs.Nodes = make(map[*graph.Node]*DFSItem, 0)

	g.Nodes(func(node *graph.Node) {
		dfs.Nodes[node] = &DFSItem{Colour: White}
	})
	n := len(dfs.Nodes)
	time := 0
	stack := make([]*graph.Node, 0, n)

	g.Nodes(func(node *graph.Node) {
		if dfs.Nodes[node].Colour != White {
			return
		}
		stack = append(stack, node)
		for len(stack) > 0 {
			xn := stack[len(stack)-1]
			item := dfs.Nodes[xn]
			switch item.Colour {
			case White:
				item.Enter = time
				time++
				f(node) // visit time!
				item.Colour = Grey
				xn.Neighbours(func(node *graph.Node, edge *graph.Edge) {
					if node.Edge(xn) != nil && g.Type == graph.Undirected {
						return
					}

					switch dfs.Nodes[node].Colour {
					case White:
						stack = append(stack, node)
					case Grey:
						dfs.HasCycle = true
					}
				})
			case Grey:
				time++
				stack = stack[0 : len(stack)-1]
				item.Exit = time
				item.Colour = Black
			default:
				stack = stack[:len(stack)-1]
			}
		}
	})
	return dfs
}
