package treesym_test

import (
	"fmt"
	xtree "hsecode.com/stdlib-tests/internal/tree"
	"hsecode.com/stdlib-tests/internal/utils"
	"hsecode.com/stdlib/tree"
	"strconv"
	"testing"
)

func TestUnit__SymWithNil(t *testing.T) {
	cases := [][]string{
		{"nil"},
		{"10", "1", "1", "nil", "3", "3", "nil"},
		{"10", "1", "1", "nil", "3", "3", "nil", "nil", "nil", "100", "nil", "nil", "100", "nil", "nil"},
	}

	for _, data := range cases {
		T, err := xtree.Decode(data)
		if err != nil {
			t.Fatal("Unexpected error in test code")
		}
		if !T.IsSym() {
			t.Fatal("IsSymTree(tree) == false for a symmetric tree")
		}
	}
}

func TestUnit__NonSymWithNil(t *testing.T) {
	cases := [][]string{
		{"23", "32"},
		{"10", "1", "1", "nil", "nil", "3"},
		{"10", "1", "1", "nil", "3", "3", "nil", "nil", "nil", "100"},
	}
	for _, data := range cases {
		T, err := xtree.Decode(data)
		if err != nil {
			t.Fatal("Unexpected error in test code")
		}
		if T.IsSym() {
			t.Fatal("IsSymTree(tree) == true for a non-symmetric tree")
		}
	}
}

func TestPerf__RandomSym(t *testing.T) {
	for i := 8; i < 20; i++ {
		T := symTree(i)
		for j := 0; j < 3; j++ {
			if !T.IsSym() {
				t.Fatalf("IsSymTree(tree) == false for a symmetric tree of height %d", i)
			}
		}
	}
}

func TestPerf__RandomNonSym(t *testing.T) {
	for i := 10; i < 1000; i++ {
		var T *tree.Tree
		for j := 0; j < i; j++ {
			xtree.Insert(&T, utils.Rand.Int())
		}
		if T.IsSym() {
			t.Fatalf("IsSymTree(tree) == true for a non-symmetric tree with %d nodes", i)
		}
	}
}

func symTree(height int) *tree.Tree {
	n := 1
	data := make([]string, 0, 2<<height-1)
	data = append(data, strconv.Itoa(utils.Rand.Int()%1000))
	for i := 0; i < height; i++ {
		for j := 0; j < n; j++ {
			v := utils.Rand.Int() % 1000
			data = append(data, strconv.Itoa(v))
		}
		r := len(data) - 1
		for j := 0; j < n; j++ {
			data = append(data, data[r-j])
		}
		n <<= 1
	}
	T, err := xtree.Decode(data)
	if err != nil {
		panic(fmt.Sprintf("unexpected error: %v", err))
	}
	return T
}

func randomTree(n int) *tree.Tree {
	data := make([]string, 0, n)
	for i := 0; i < n; i++ {
		data = append(data, strconv.Itoa(utils.Rand.Int()%2000))
	}
	T, err := xtree.Decode(data)
	if err != nil {
		panic(fmt.Sprintf("unexpected error: %v", err))
	}
	return T
}
