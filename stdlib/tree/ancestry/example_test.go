package ancestry_test

import (
	"fmt"
	"hsecode.com/stdlib/v2/tree"
	"hsecode.com/stdlib/v2/tree/ancestry"
)

func Example() {
	// Create BST:
	//     6
	//    / \
	//   2   8
	//  / \
	// 1   3
	T := &tree.Tree{
		Value: 6,
		Left:  &tree.Tree{Value: 2, Left: &tree.Tree{Value: 1}, Right: &tree.Tree{Value: 3}},
		Right: &tree.Tree{Value: 8},
	}

	A := ancestry.New(T)

	fmt.Println(A.IsDescendant(2, 2)) // false, not a _proper_ descentant
	fmt.Println(A.IsDescendant(2, 3)) // true
	fmt.Println(A.IsDescendant(6, 3)) // true
	fmt.Println(A.IsDescendant(1, 8)) // false

	// Output:
	// false
	// true
	// true
	// false
}
