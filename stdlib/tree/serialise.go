package tree

import (
	"fmt"
	"strconv"
)

const nilValue = "nil"

// Decode creates a binary tree from the given serialised representation.
// Return error if any of the array values are neither intergers or nils,
// or if the array describes an impossible tree structure (e.g. ["nil", "2", "3"]).
//
// See Serialisation Format section in the package overview for more details on the format.
func Decode(data []string) (*Tree, error) {
	var root *Tree
	level := []**Tree{&root}
	var i int
	n := len(data)

	for len(level) > 0 {
		nextLevel := make([]**Tree, 0)
		for _, node := range level {
			if i < n {
				if data[i] != nilValue {
					if node == nil {
						return nil, fmt.Errorf("invalid representation: node %v under a nil node", data[i])
					}
					val, err := strconv.Atoi(data[i])
					if err != nil {
						return nil, fmt.Errorf("one of the items is neither interger or nil: %v", data[i])
					}
					*node = &Tree{val, nil, nil}
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

// Encode returns a serialised representaion of the given tree.
//
// See Serialisation Format section in the package overview for more details on the format.
func (T *Tree) Encode() []string {
	if T == nil {
		return []string{}
	}
	md := maxDepth(T)
	depth := 1
	var result []string

	level := []*Tree{T}

	for len(level) > 0 {
		nextLevel := make([]*Tree, 0)

		for _, node := range level {
			if node != nil {
				result = append(result, strconv.Itoa(node.Value))
				if depth < md {
					nextLevel = append(nextLevel, node.Left, node.Right)
				}
			} else {
				result = append(result, nilValue)
				if depth < md {
					nextLevel = append(nextLevel, nil, nil)
				}
			}
		}
		depth++
		level = nextLevel
	}
	n := 0
	for i := len(result) - 1; i >= 0 && result[i] == nilValue; i-- {
		n++
	}
	return result[:len(result)-n]
}

func maxDepth(tree *Tree) int {
	if tree == nil {
		return 0
	}
	lm := maxDepth(tree.Left)
	rm := maxDepth(tree.Right)
	if lm > rm {
		return 1 + lm
	}
	return 1 + rm
}
