package treebalanced_test

import (
	xtree "hsecode.com/stdlib-tests/internal/tree"
	"hsecode.com/stdlib-tests/internal/utils"
	"hsecode.com/stdlib/tree"
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
		inorder := xtree.InOrder(bst)

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
