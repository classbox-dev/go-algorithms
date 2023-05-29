package vector_test

import (
	// imported package is renamed to avoid conflict with type `int`
	"fmt"
	"hsecode.com/stdlib/v2/vector"
)

func ExampleVector_Push() {
	vec := vector.New[int](0)
	for _, p := range []int{2, 3, 5, 7, 11} {
		vec.Push(p)
	}
	fmt.Println(vec.Get(2))
	// Output: 5
}

func ExampleVector_Pop() {
	vec := vector.New[int](0)
	vec.Push(10)
	vec.Push(11)
	fmt.Println(vec.Pop())
	// Output: 11
}

func ExampleVector_Insert() {
	vec := vector.New[int](0)
	vec.Push(42)
	vec.Insert(0, 41)
	fmt.Println(vec.Get(0), vec.Get(1))
	// Output: 41 42
}
