package tree_test

import (
	"fmt"
	"hsecode.com/stdlib/tree"
)

func ExampleDecode() {
	T, _ := tree.Decode([]string{"2", "1", "3", "nil", "nil", "nil", "10"})

	// Resulting tree:
	//    2
	//   / \
	//  1    3
	//        \
	//         10

	T.InOrder(func(node *tree.Tree) { fmt.Println(node.Value) })
	// Output:
	// 1
	// 2
	// 3
	// 10
}

func ExampleNewBST() {
	// Create BST from slice:
	elements := []int{11, 2, 5, 7, 3, 3, 5, 2}
	T := tree.NewBST(elements)
	// Print all nodes in sorted order
	T.InOrder(func(node *tree.Tree) { fmt.Println(node.Value) })
	// Output:
	// 2
	// 3
	// 5
	// 7
	// 11
}

func ExampleTree_InOrder() {
	// Create BST:
	//    2
	//   / \
	//  1    3
	//        \
	//         10
	T := &tree.Tree{
		Value: 2,
		Left:  &tree.Tree{Value: 1},
		Right: &tree.Tree{Value: 3, Right: &tree.Tree{Value: 10}},
	}
	// Print all nodes in sorted order
	T.InOrder(func(node *tree.Tree) { fmt.Println(node.Value) })
	// Output:
	// 1
	// 2
	// 3
	// 10
}

func ExampleTree_Encode() {
	// Create BST:
	//    2
	//   / \
	//  1    3
	//        \
	//         10
	T := &tree.Tree{
		Value: 2,
		Left:  &tree.Tree{Value: 1},
		Right: &tree.Tree{Value: 3, Right: &tree.Tree{Value: 10}},
	}
	// Print serialised representation
	fmt.Println(T.Encode())
	// Output: [2 1 3 nil nil nil 10]
}
