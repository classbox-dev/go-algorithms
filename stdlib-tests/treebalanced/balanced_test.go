package treebalanced_test

import (
	xtree "hsecode.com/stdlib-tests/v2/internal/tree"
	"hsecode.com/stdlib-tests/v2/internal/utils"
	"hsecode.com/stdlib/v2/tree"
	"reflect"
	"sort"
	"testing"
)

func TestUnit__Random(t *testing.T) {

	for n := 0; n < 1200; n++ {
		elements := utils.SliceRandom((n+10)/2, n)
		unique := utils.Unique(elements)
		sort.Ints(unique)

		bst := tree.NewBST(elements)
		inorder := make([]int, 0, len(elements))
		xtree.InOrder(bst, func(node *tree.Tree) { inorder = append(inorder, node.Value) })

		if !reflect.DeepEqual(inorder, unique) {
			t.Fatalf("The tree contains wrong elements, duplicates, or is not a valid BST.\nInput: %v", elements)
		}

		if bst != nil {
			diff := xtree.MaxDepth(bst.Left) - xtree.MaxDepth(bst.Right)
			if diff < 0 {
				diff = -diff
			}
			if diff > 1 {
				t.Fatalf("BST is not balanced\nInput: %v", elements)
			}
		}
	}
}

func TestPerf__Random(t *testing.T) {
	for n := 6000; n < 6700; n++ {
		elements := utils.SliceRandom((n+10)/2, n)
		tree.NewBST(elements)
	}
}
