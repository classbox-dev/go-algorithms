package ancestry

import "hsecode.com/stdlib/v2/tree"

type time struct {
	exit  int
	enter int
}

// Ancestry is a data structure that can determine if a certain tree node is a proper descendant of another.
type Ancestry struct {
	nodes map[int]time
}

// New creates an instance of Ancestry for the given binary tree with unique integer nodes.
// The running time is O(N) for N nodes.
func New(T *tree.Tree) *Ancestry {
	a := new(Ancestry)
	a.nodes = make(map[int]time)
	ts := 0
	var dfs func(T *tree.Tree)
	dfs = func(T *tree.Tree) {
		if T == nil {
			return
		}
		t := new(time)
		t.enter = ts
		ts++
		dfs(T.Left)
		dfs(T.Right)
		t.exit = ts
		ts++
		a.nodes[T.Value] = *t
	}
	dfs(T)
	return a
}

// IsDescendant determines if node b is a proper descendant of node a.
// The running time is guaranteed to be Θ(1).
// Panics if the given nodes are not in the tree.
func (A *Ancestry) IsDescendant(a, b int) bool {
	ta, ok1 := A.nodes[a]
	tb, ok2 := A.nodes[b]
	if !ok1 || !ok2 {
		panic("unknown nodes")
	}
	return (ta.enter < tb.enter) && (tb.enter < tb.exit) && (tb.exit < ta.exit)
}
