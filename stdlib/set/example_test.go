package set_test

import (
	"fmt"
	"hsecode.com/stdlib/set"
)

type Int int

func (v Int) Less(other set.Element) bool {
	return int(v) < int(other.(Int))
}

func (v Int) Equal(other set.Element) bool {
	return int(v) == int(other.(Int))
}

func Example() {
	s := set.New()
	for _, p := range []int{5, 3, 2} {
		s.Insert(Int(p))
	}

	v, ok := s.Find(Int(3)) // Exact lookup
	fmt.Println(v, ok)

	it := s.LowerBound(Int(3)) // Range lookup
	for it.Next() {
		fmt.Println(it.Value())
	}

	ok = s.Insert(Int(7))
	fmt.Println(ok)

	ok = s.Delete(Int(2))
	fmt.Println(ok)

	// Output:
	// 3 true
	// 3
	// 5
	// true
	// true
}
