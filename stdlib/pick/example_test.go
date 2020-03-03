package pick_test

import (
	// imported package is renamed to avoid conflict with type `int`
	"fmt"
	"hsecode.com/stdlib/pick"
	"sort"
)

func ExampleFirstN() {
	// Indices of top-3 elements: 0, 2, 5
	data := []int{13, 2, 29, 3, 7, 23}
	indices := pick.FirstN(sort.Reverse(sort.IntSlice(data)), 3)

	// FirstN returns indices in arbitrary order, sort them to ensure stable result
	sort.IntSlice(indices).Sort()

	fmt.Println(indices)
	// Output:
	// [0 2 5]
}

func ExampleNthElement() {
	data := []int{13, 2, 29, 3, 7, 23, 31}
	m := len(data) / 2 // median index if the array was sorted
	pick.NthElement(sort.IntSlice(data), m)
	fmt.Println(data[m])
	// Output:
	// 13
}
