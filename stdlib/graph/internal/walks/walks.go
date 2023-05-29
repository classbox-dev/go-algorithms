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
	Nodes    map[int]*DFSItem
	HasCycle bool
}

type DFSItem struct {
	Colour Colour
	Enter  int
	Exit   int
}

func NewDFS(g *graph.Graph, f func(graph.Node)) *DFS {
	dfs := new(DFS)
	dfs.Nodes = make(map[int]*DFSItem)

	g.Nodes(func(node graph.Node) {
		dfs.Nodes[node.ID()] = &DFSItem{Colour: White}
	})
	n := len(dfs.Nodes)
	time := 0
	stack := make([]int, 0, n)

	g.Nodes(func(node graph.Node) {
		if dfs.Nodes[node.ID()].Colour != White {
			return
		}

		stack = append(stack, node.ID())

		for len(stack) > 0 {
			xn := stack[len(stack)-1]
			item := dfs.Nodes[xn]

			switch item.Colour {

			case White:
				item.Enter = time
				time++
				f(node) // visit time!
				item.Colour = Grey

				g.Neighbours(xn, func(node graph.Node, data interface{}) {
					// skip back edges in undirected graphs
					if _, ok := g.Edge(node.ID(), xn); ok && g.Type == graph.Undirected {
						return
					}
					switch dfs.Nodes[node.ID()].Colour {
					case White:
						stack = append(stack, node.ID())
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
