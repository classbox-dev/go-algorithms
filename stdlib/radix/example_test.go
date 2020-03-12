package radix_test

import (
	"fmt"
	"hsecode.com/stdlib/radix"
)

func ExampleSort() {
	data := []uint64{13, 2, 29, 3, 7, 23, 31}
	radix.Sort(data)
	fmt.Println(data)
	// Output:
	// [2 3 7 13 23 29 31]
}
