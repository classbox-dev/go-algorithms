package tree

// Tree is a binary tree with int values.
type Tree struct {
	Value int
	Left  *Tree
	Right *Tree
}
