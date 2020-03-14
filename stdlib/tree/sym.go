package tree

// IsSym tests whether the tree is a mirror of itself (i.e. symmetric around its center).
// For example, this tree is symmetric:
//
//      1
//     / \
//    2   2
//   / \ / \
//  3  4 4  3
//
// And this one is not:
//
//    1
//   / \
//  2   2
//   \   \
//   3    3
func (T *Tree) IsSym() bool {
	if T == nil {
		return true
	}
	return twoSym(T.Left, T.Right)
}

func twoSym(l *Tree, r *Tree) bool {
	if l == nil && r == nil {
		return true
	} else if l == nil || r == nil {
		return false
	}
	return l.Value == r.Value && twoSym(l.Left, r.Right) && twoSym(l.Right, r.Left)
}
