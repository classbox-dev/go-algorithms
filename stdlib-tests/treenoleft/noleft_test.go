package treenoleft_test

import (
	xtree "hsecode.com/stdlib-tests/v2/internal/tree"
	"hsecode.com/stdlib-tests/v2/internal/utils"
	"hsecode.com/stdlib/v2/tree"
	"sort"
	"testing"
)

const (
	MinInt = -int(^uint(0)>>1) - 1
)

type Addrs map[*tree.Tree]int

func testTree(t *testing.T, T *tree.Tree, addrs Addrs) {
	cnt := 0
	prev := MinInt
	xtree.InOrder(T, func(node *tree.Tree) {
		if v, ok := addrs[node]; !ok || v != node.Value {
			t.Fatal("NoLeft() is not inplace: tree nodes have been reallocated, duplicated, or changed")
		}
		if node.Value < prev {
			t.Fatal("NoLeft() returned an invalid binary search tree")
		}
		prev = node.Value
		cnt++
	})
	if cnt != len(addrs) {
		t.Fatal("NoLeft() returned a tree with less nodes than expected")
	}
}

func TestUnit__BST(t *testing.T) {
	for n := 0; n < 250; n++ {
		elems := utils.Unique(utils.SliceRandom((n+1)*5, n))
		sort.Ints(elems)
		T := xtree.NewBST(elems)
		addrs := make(Addrs, n)
		xtree.InOrder(T, func(node *tree.Tree) {
			addrs[node] = node.Value
		})
		result := T.NoLeft()
		testTree(t, result, addrs)
	}
}

func TestUnit__UnbalancedBST(t *testing.T) {
	for n := 1; n < 3000; n += 5 {
		var T *tree.Tree
		for j := 0; j < n; j++ {
			xtree.Insert(&T, utils.Rand.Int())
		}
		addrs := make(Addrs, n)
		xtree.InOrder(T, func(node *tree.Tree) {
			addrs[node] = node.Value
		})
		result := T.NoLeft()
		testTree(t, result, addrs)
	}
}

func TestPerf__UnbalancedBST(t *testing.T) {
	N := 3000
	randoms := make([]int, 2*N)
	for i := 0; i < 2*N; i++ {
		randoms[i] = utils.Rand.Intn(2 * N)
	}
	for n := 1; n < N; n += 1 {
		var T *tree.Tree
		offset := utils.Rand.Intn(N)
		for j := 0; j < n; j++ {
			xtree.Insert(&T, randoms[offset+j])
		}
		r := T.NoLeft()
		utils.Use(r.Value)
	}
}
