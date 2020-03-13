package tree

import (
	"hsecode.com/stdlib/tree"
	"sort"
	"strconv"
)

func NewBST(elements []int) *tree.Tree {
	unique := make([]int, 0, len(elements))
	same := make(map[int]struct{}, len(elements))
	for _, e := range elements {
		if _, ok := same[e]; !ok {
			same[e] = struct{}{}
			unique = append(unique, e)
		}
	}
	sort.Ints(unique)
	return New(unique)
}

func New(data []int) *tree.Tree {
	n := len(data)
	if n == 0 {
		return nil
	}
	node := new(tree.Tree)
	node.Value = data[n/2]
	node.Left = New(data[0 : n/2])
	node.Right = New(data[n/2+1:])
	return node
}

func Insert(T **tree.Tree, elem int) bool {
	node := *T
	var point = T
	for node != nil {
		switch {
		case node.Value > elem:
			point = &(node.Left)
			node = *point
		case node.Value < elem:
			point = &(node.Right)
			node = *point
		case node.Value == elem:
			return false
		}
	}
	newNode := new(tree.Tree)
	newNode.Value = elem
	*point = newNode
	return true
}

func Serialise(T *tree.Tree) []string {
	if T == nil {
		return []string{}
	}

	md := MaxDepth(T)
	depth := 1

	var result []string

	level := []*tree.Tree{T}

	for len(level) > 0 {
		nextLevel := make([]*tree.Tree, 0)

		for _, node := range level {
			if node != nil {
				result = append(result, strconv.Itoa(node.Value))
				if depth < md {
					nextLevel = append(nextLevel, node.Left, node.Right)
				}
			} else {
				result = append(result, "nil")
				if depth < md {
					nextLevel = append(nextLevel, nil, nil)
				}
			}
		}

		depth++
		level = nextLevel
	}

	n := 0
	for i := len(result) - 1; i >= 0 && result[i] == "nil"; i-- {
		n++
	}
	return result[:len(result)-n]
}

func MaxDepth(tree *tree.Tree) int {
	if tree == nil {
		return 0
	}
	lm := MaxDepth(tree.Left)
	rm := MaxDepth(tree.Right)
	if lm > rm {
		return 1 + lm
	}
	return 1 + rm
}

func InOrder(bst *tree.Tree) []int {
	if bst == nil {
		return []int{}
	}
	result := make([]int, 0)
	result = append(result, InOrder(bst.Left)...)
	result = append(result, bst.Value)
	result = append(result, InOrder(bst.Right)...)
	return result
}
