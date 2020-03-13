// Package tree is a collection of algorithms on binary trees.
/*
Serialisation Format

We define a serialised format of a binary tree is an array of strings constructed with level-order traversal of the tree.
Numbers denote node values, and "nil" denotes an absence of node.

For example, an array ["6", "5", "7", "2", "nil", "nil", "8"] represents the following tree:

      __6
     /   \
    5     7
   /       \
  2         8

The representation must not contain trailing "nil"s at the end
(even though ["6", "5", "7", "2", "nil", "nil", "8", "nil", "nil", "nil"] technically denotes the same tree).

More examples can be found at https://support.leetcode.com/hc/en-us/articles/360011883654-What-does-1-null-2-3-mean-in-binary-tree-representation- (with a minor difference of using "null" instead of "nil").

Function Encode and Decode from this package use this serialised format.
Encode turns the given tree into an array of strings and Decode recreates a tree from the given array of strings.

Restrictions

InOrder must be implemented non-recursively.
*/
package tree
