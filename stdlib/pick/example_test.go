package pick_test

import (
	// imported package is renamed to avoid conflict with type `int`
	"fmt"
	"hsecode.com/stdlib/pick"
	"sort"
)

func ExampleFirstN() {
	// Indices of min-3 elements: 1, 3, 4
	data := []int{13, 2, 29, 3, 7, 23}
	indices := pick.FirstN(data, 3)

	// FirstN returns indices in arbitrary order, sort them to ensure stable result
	sort.IntSlice(indices).Sort()

	fmt.Println(indices)
	// Output:
	// [1 3 4]
}

func ExampleNthElement() {
	data := []int{13, 2, 29, 3, 7, 23, 31}
	m := len(data) / 2 // median index if the array was sorted
	pick.NthElement(data, m)
	fmt.Println(data[m])
	// Output:
	// 13
}
