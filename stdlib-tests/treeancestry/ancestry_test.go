package treeancestry_test

import (
	matrix "hsecode.com/stdlib-tests/internal/matrix/bool"
	xtree "hsecode.com/stdlib-tests/internal/tree"
	"hsecode.com/stdlib-tests/internal/utils"
	"hsecode.com/stdlib/tree"
	"hsecode.com/stdlib/tree/ancestry"
	"testing"
)

func insert(T **tree.Tree, elem int, ids map[int]int, mat *matrix.Matrix) bool {
	node := *T
	point := T
	for node != nil {
		mat.Set(ids[node.Value], ids[elem], true)
		switch {
		case node.Value > elem:
			point = &(node.Left)
			node = *point
		case node.Value < elem:
			point = &(node.Right)
			node = *point
		}
	}
	newNode := new(tree.Tree)
	newNode.Value = elem
	*point = newNode
	return true
}

func TestUnit__UnbalancedBST(t *testing.T) {
	N := 250
	for n := 1; n < N; n++ {
		elems := utils.Unique(utils.SliceRandom(n*3, n))
		ids := make(map[int]int, n)
		for i, r := range elems {
			ids[r] = i
		}
		M := matrix.New(len(ids), len(ids))

		var T *tree.Tree
		for _, e := range elems {
			insert(&T, e, ids, M)
		}
		A := ancestry.New(T)

		for parent, i := range ids {
			for child, j := range ids {
				isDesc := M.Get(i, j)
				result := A.IsDescendant(parent, child)
				if result != isDesc {
					t.Fatalf("IsDescendant(%d, %d) == %v, expected %v\nInsertion order of input BST: %v", parent, child, result, isDesc, elems)
				}
			}
		}
	}
}

func TestPerf__UnbalancedBST(t *testing.T) {
	N := 840
	for n := 800; n < N; n++ {
		var T *tree.Tree
		elems := utils.Unique(utils.SliceRandom(n*3, n))
		for _, e := range elems {
			xtree.Insert(&T, e)
		}
		A := ancestry.New(T)
		for _, parent := range elems {
			for _, child := range elems {
				A.IsDescendant(parent, child)
			}
		}
	}
}
