package set_test

import (
	"fmt"
	"hsecode.com/stdlib/set"
)

func Example() {
	s := set.New[int]()
	for _, p := range []int{5, 3, 2} {
		s.Insert(p)
	}

	v, ok := s.Find(3) // Exact lookup
	fmt.Println(v, ok)

	it := s.LowerBound(3) // Range lookup
	for it.Next() {
		fmt.Println(it.Value())
	}

	ok = s.Insert(7)
	fmt.Println(ok)

	ok = s.Delete(2)
	fmt.Println(ok)

	// Output:
	// 3 true
	// 3
	// 5
	// true
	// true
}
