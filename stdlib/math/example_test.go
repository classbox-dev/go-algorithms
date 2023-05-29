package math_test

import (
	// imported package is renamed to avoid conflict with type `int`
	"fmt"
	"hsecode.com/stdlib/math"
)

func ExampleNthPrime() {
	fmt.Println(math.NthPrime(1), math.NthPrime(2), math.NthPrime(3))

	// https://www.wolframalpha.com/input/?i=Prime[24368]
	fmt.Println(math.NthPrime(24368))

	// Output:
	// 2 3 5
	// 279121
}
