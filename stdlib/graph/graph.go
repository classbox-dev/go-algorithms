package graph

type sentinel struct{}

// Type represents a choice of graph types.
type Type int

const (
	Directed Type = iota
	Undirected
)

// Graph represents either directed or undirected graph.
type Graph struct {
	// Type identifies whether the graph is directed or undirected.
	Type   Type
	lookup map[int]Node
	nodes  map[int]map[int]interface{}
}

// New creates an empty graph of the given type (directed or undirected).
func New(graphType Type) *Graph {
	g := new(Graph)
	g.nodes = make(map[int]map[int]interface{})
	g.lookup = make(map[int]Node)
	g.Type = graphType

	return g
}

// Node is an interface for user-provided node type.
type Node interface {
	// ID returns a graph-unique integer identifier.
	ID() int
}

// AddNode creates a new node or overwrites an existing one with the same ID.
func (g *Graph) AddNode(node Node) {
	id := node.ID()
	g.lookup[id] = node

	if _, ok := g.nodes[id]; !ok {
		g.nodes[id] = make(map[int]interface{})
	}
}

// Node retrieves a node by the given node ID.
// The boolean indicates whether the node was found.
func (g *Graph) Node(id int) (Node, bool) {
	node, ok := g.lookup[id]
	return node, ok
}

// Neighbours applies the provided function f to all nodes adjacent to u.
// Panics if the node does not exist.
func (g *Graph) Neighbours(u int, f func(v Node, edgeData interface{})) {
	adj, ok := g.nodes[u]
	if !ok {
		panic("node does not exist")
	}

	for id, data := range adj {
		f(g.lookup[id], data)
	}
}

// Nodes applies the provided function f to all nodes in arbitrary order.
func (g *Graph) Nodes(f func(Node)) {
	for _, node := range g.lookup {
		f(node)
	}
}

// Edge returns data from the edge (u,v).
// The boolean indicates whether the edge was found.
func (g *Graph) Edge(u, v int) (interface{}, bool) {
	adj, ok := g.nodes[u]
	if !ok {
		return nil, false
	}

	data, ok := adj[v]
	if !ok {
		return nil, false
	}

	return data, true
}

// Edges applies the provided function f to all graph edges in arbitrary order.
//
// NOTE: for undirected graphs the function is called once for every connected pair of nodes.
// For example, for the graph (2)--(3) the function must be called once either with (u=2,v=3) or (u=3,v=2).
func (g *Graph) Edges(f func(u, v Node, edgeData interface{})) {
	switch g.Type {
	case Directed:
		for u, adj := range g.nodes {
			for v, data := range adj {
				f(g.lookup[u], g.lookup[v], data)
			}
		}
	case Undirected:
		type edge struct{ u, v int }

		edges := make(map[edge]sentinel)

		for u, adj := range g.nodes {
			for v, data := range adj {
				if _, ok := edges[edge{v, u}]; ok {
					continue
				}

				edges[edge{u, v}] = sentinel{}

				f(g.lookup[u], g.lookup[v], data)
			}
		}
	}
}

// AddEdge creates an edge between nodes u and v augmented with the provided data.
// If the edge between u and v already exists, the method overwrites the edge data.
// Panics if the given nodes do not exist.
func (g *Graph) AddEdge(u, v int, edgeData interface{}) {
	if _, ok := g.lookup[u]; !ok {
		panic("source node is not in graph")
	}
	if _, ok := g.lookup[v]; !ok {
		panic("target node is not in graph")
	}

	g.nodes[u][v] = edgeData

	if g.Type == Undirected {
		g.nodes[v][u] = edgeData
	}
}
