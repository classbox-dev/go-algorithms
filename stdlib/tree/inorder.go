package tree

// InOrder performs in-order traversal of the tree
// and calls the provided visit function for each tree node.
// The function must be implemented non-recursively.
func (T *Tree) InOrder(visit func(node *Tree)) {
	stack := make([]*Tree, 0)
	node := T
	for len(stack) > 0 || node != nil {
		if node != nil {
			stack = append(stack, node)
			node = node.Left
		} else {
			node = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			visit(node)
			node = node.Right
		}
	}
}
