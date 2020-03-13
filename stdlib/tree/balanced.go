package tree

import "sort"

// NewBST constructs a balanced binary search tree from the elements of the given array in O(n log n) time.
// The tree can only contain unique elements even if the input array has duplicates.
func NewBST(elements []int) *Tree {
	unique := make([]int, 0, len(elements))
	same := make(map[int]struct{}, len(elements))
	for _, e := range elements {
		if _, ok := same[e]; !ok {
			same[e] = struct{}{}
			unique = append(unique, e)
		}
	}
	sort.Ints(unique)
	return fromSorted(unique)
}

func fromSorted(unique []int) *Tree {
	n := len(unique)
	if n == 0 {
		return nil
	}
	node := new(Tree)
	node.Value = unique[n/2]
	node.Left = fromSorted(unique[0 : n/2])
	node.Right = fromSorted(unique[n/2+1:])
	return node
}
