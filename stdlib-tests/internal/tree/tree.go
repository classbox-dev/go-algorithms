package tree

import (
	"fmt"
	"hsecode.com/stdlib/v2/tree"
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

func Encode(T *tree.Tree) []string {
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
	return Normalise(result)
}

func Decode(data []string) (*tree.Tree, error) {
	var root *tree.Tree
	level := []**tree.Tree{&root}
	var i int
	n := len(data)

	for len(level) > 0 {
		nextLevel := make([]**tree.Tree, 0)
		for _, node := range level {
			if i < n {
				if data[i] != "nil" {
					if node == nil {
						return nil, fmt.Errorf("invalid representation: node %v under a nil node", data[i])
					}
					val, err := strconv.Atoi(data[i])
					if err != nil {
						return nil, fmt.Errorf("one of the items is neither iterger or nil: %v", data[i])
					}
					*node = &tree.Tree{val, nil, nil}
					nextLevel = append(nextLevel, &((*node).Left), &((*node).Right))
				} else {
					nextLevel = append(nextLevel, nil, nil)
				}
				i++
			} else {
				nextLevel = nextLevel[:0]
				break
			}
		}
		level = nextLevel
	}
	return root, nil
}

func Normalise(data []string) []string {
	n := 0
	for i := len(data) - 1; i >= 0 && data[i] == "nil"; i-- {
		n++
	}
	return data[:len(data)-n]
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

func InOrder(T *tree.Tree, visit func(node *tree.Tree)) {
	stack := make([]*tree.Tree, 0)
	node := T
	for len(stack) > 0 || node != nil {
		if node != nil {
			stack = append(stack, node)
			node = node.Left
		} else {
			node = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			visit(node)
			node = node.Right
		}
	}
}
