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
	lookup map[int]*Node
	nodes  map[*Node]sentinel
}

// New creates an empty graph of the given type (directed or undirected).
func New(type_ Type) *Graph {
	g := new(Graph)
	g.nodes = make(map[*Node]sentinel, 0)
	g.lookup = make(map[int]*Node, 0)
	g.Type = type_
	return g
}

// Node is a graph node holding a user-provided value.
type Node struct {
	Value      NodeValue
	neighbours map[*Node]*Edge
}

// Edge returns an edge (u,v) if it exists; nil otherwise.
func (u *Node) Edge(v *Node) *Edge {
	if edge, ok := u.neighbours[v]; ok {
		return edge
	}
	return nil
}

// Neighbours applies the provided function f to all nodes pointed by edges from the current node u.
func (u *Node) Neighbours(f func(*Node, *Edge)) {
	for n, edge := range u.neighbours {
		f(n, edge)
	}
}

// NodeValue is an interface for user-provided node value.
type NodeValue interface {
	// ID returns a graph-unique integer igentifier.
	ID() int
}

// Edge is a graph edge holding a user-provided value.
type Edge struct {
	Value interface{}
}

// AddNode creates a new graph node with the given value or overwrites the value of an existing node with the same ID.
func (g *Graph) AddNode(value NodeValue) *Node {
	key := value.ID()
	if node, ok := g.lookup[key]; ok {
		node.Value = value
		return node
	}
	node := &Node{Value: value, neighbours: make(map[*Node]*Edge, 0)}
	g.nodes[node] = sentinel{}
	g.lookup[key] = node
	return node
}

// Node retrieves a graph node by the given node ID. Returns nil if the ID does not exist.
func (g *Graph) Node(id int) *Node {
	if node, ok := g.lookup[id]; ok {
		return node
	}
	return nil
}

// Nodes applies the provided function f to all graph nodes in arbitrary order.
func (g *Graph) Nodes(f func(*Node)) {
	for node := range g.nodes {
		f(node)
	}
}

// Edges applies the provided function f to all graph edges in arbitrary order.
//
// NOTE: for undirected graphs the function is called once for every connected pair of nodes.
// For example, for the graph (2)--(3) the function must be called once either with (u=2,v=3) or (u=3,v=2).
func (g *Graph) Edges(f func(u, v *Node, e *Edge)) {
	edges := make(map[*Edge]sentinel)
	for u := range g.nodes {
		for v, e := range u.neighbours {
			if _, ok := edges[e]; ok {
				continue
			}
			edges[e] = sentinel{}
			f(u, v, e)
		}
	}
}

// AddEdge sets an edge augmented with the provided data between nodes u and v from the graph.
// If the edge between u and v already exists, a new edge is created an overwrites the existing one.
// Panics if the nodes do not exist in the graph.
func (g *Graph) AddEdge(u, v *Node, data interface{}) *Edge {
	if _, ok := g.nodes[u]; !ok {
		panic("source node is not in graph")
	}
	if _, ok := g.nodes[v]; !ok {
		panic("target node is not in graph")
	}
	edge := &Edge{data}
	u.neighbours[v] = edge
	if g.Type == Undirected {
		v.neighbours[u] = edge
	}
	return edge
}
