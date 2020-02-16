package cmp_test

import (
	// imported package is renamed to avoid conflict with type `int`
	cmp "hsecode.com/stdlib/cmp/int"
)

var primes = []int{7, 11, 2, 3, 23, 5}

func Example() {
	cmp.Min(primes...)
	// Output: 2

	cmp.Max(primes...)
	// Output: 23

	cmp.Min(-4, -5)
	// Output: -5

	cmp.Max()
	// panic("Max requires at least one argument")
}
