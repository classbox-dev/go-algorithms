package tree

// NoLeft transforms the binary search tree into an equivalent one that has no left subtrees.
// Returns a pointer to the new root.
// The tree must be transformed in-place i.e. the result can only contain the originally allocated nodes with their original values.
// The function assumes that T is already a valid binary search tree.
func (T *Tree) NoLeft() *Tree {
	if T == nil {
		return nil
	}
	var newRoot *Tree
	left := T.Left.NoLeft()
	if left == nil {
		newRoot = T
	} else {
		newRoot = left
		p := newRoot
		for p.Right != nil {
			p = p.Right
		}
		p.Right = T
	}
	T.Left = nil
	T.Right = T.Right.NoLeft()
	return newRoot
}
