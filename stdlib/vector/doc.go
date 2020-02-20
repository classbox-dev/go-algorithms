// Package vector implements dynamic array.
/*
Memory management

On creation, an underlying slice of the given capacity is allocated.

On insertion (Push, Insert), new elements are put into unoccupied space of the underlying slice.
If there is no room, the underlying slice must be expanded by a constant factor α>1. (α=4/3 in the reference implementation.)

On deletion (Pop, Delete), the underlying slice must be shrinked if only β*capacity (β<1) elements are occupied.
Note that the threshold (β*capacity) must be strictly smaller than (capacity/α) to avoid repeated growing/shrinking on successive insertions/deletions.

The growing/shriking strategy must not allow memory leaks: vector with all elements deleted must always have an underlying slice of near-zero length.

Restrictions

Built-in append() function is not permitted in this package.
*/
package vector
