package treeinorder_test

import (
	xtree "hsecode.com/stdlib-tests/internal/tree"
	"hsecode.com/stdlib-tests/internal/utils"
	"hsecode.com/stdlib/tree"
	"reflect"
	"runtime/debug"
	"sort"
	"testing"
)

func TestUnit__BST(t *testing.T) {
	for n := 0; n < 250; n++ {
		data := utils.Unique(utils.SliceRandom(2*(100+n), n))
		sort.Ints(data)

		T := xtree.New(data)
		output := make([]int, 0, len(data))
		T.InOrder(func(node *tree.Tree) { output = append(output, node.Value) })

		if !reflect.DeepEqual(output, data) {
			t.Fatalf("Invalid result.\nOutput: %v\nSerialised tree: %v", output, xtree.Serialise(T))
		}
	}
}

func TestUnit__NonBST(t *testing.T) {
	for n := 0; n < 250; n++ {
		data := utils.Unique(utils.SliceRandom(2*(100+n), n))

		T := xtree.New(data)
		output := make([]int, 0, len(data))
		T.InOrder(func(node *tree.Tree) { output = append(output, node.Value) })

		if !reflect.DeepEqual(output, data) {
			t.Fatalf("Invalid result.\nOutput: %v\nSerialised tree: %v", output, xtree.Serialise(T))
		}
	}
}

func TestUnit__UnbalancedBST(t *testing.T) {

	for n := 0; n < 250; n++ {
		data := utils.Unique(utils.SliceRandom(2*(100+n), n))
		var T *tree.Tree
		for _, x := range data {
			xtree.Insert(&T, x)
		}

		sorted := make([]int, len(data))
		copy(sorted, data)
		sort.Ints(sorted)

		output := make([]int, 0, len(data))
		T.InOrder(func(node *tree.Tree) { output = append(output, node.Value) })

		if !reflect.DeepEqual(output, sorted) {
			t.Fatalf("Invalid result.\nOutput: %v\nInsertion order: %v", output, data)
		}
	}
}

func TestUnit__EnsureNonRecursive(t *testing.T) {
	data := utils.Unique(utils.SliceRandom(10000, 1000))
	sort.Ints(data)
	var T *tree.Tree
	for _, x := range data {
		xtree.Insert(&T, x)
	}
	oldSize := debug.SetMaxStack(1024)
	defer debug.SetMaxStack(oldSize)
	T.InOrder(func(node *tree.Tree) {})
}

func TestPerf__UnbalancedBST(t *testing.T) {
	for n := 20; n < 3500; n++ {
		var T *tree.Tree
		for i := 0; i < n; {
			x := utils.Rand.Intn(n * 5)
			if xtree.Insert(&T, x) {
				i++
			}
		}
		T.InOrder(func(node *tree.Tree) {})
	}
}
