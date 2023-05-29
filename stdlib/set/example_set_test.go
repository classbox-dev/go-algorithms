package set_test

import (
	"fmt"
	"hsecode.com/stdlib/set"
)

func ExampleSet_Insert() {
	s := set.New[int]()
	fmt.Println(s.Insert(7)) // true, the element was inserted
	fmt.Println(s.Insert(7)) // false, duplicated element, do nothing
	// Output:
	// true
	// false
}

func ExampleSet_Delete() {
	s := set.New[int]()
	for _, p := range []int{7, 11, 2, 3, 13, 5} {
		s.Insert(p)
	}
	fmt.Println(s.Delete(7)) // true, the element was deleted
	fmt.Println(s.Delete(7)) // false, no such element

	it := s.Begin()
	for it.Next() {
		fmt.Println(it.Value())
	}
	// Output:
	// true
	// false
	// 2
	// 3
	// 5
	// 11
	// 13
}

func ExampleSet_Find() {
	s := set.New[int]()
	for _, p := range []int{7, 11, 2, 3, 13, 5} {
		s.Insert(p)
	}
	fmt.Println(s.Find(13))
	fmt.Println(s.Find(10))
	// Output:
	// 13 true
	// 0 false
}

func ExampleSet_Begin() {
	s := set.New[int]()
	for _, p := range []int{7, 11, 2, 3, 13, 5} {
		s.Insert(p)
	}
	it := s.Begin()
	for it.Next() {
		fmt.Println(it.Value())
	}
	// Output:
	// 2
	// 3
	// 5
	// 7
	// 11
	// 13
}

func ExampleSet_End() {
	s := set.New[int]()
	for _, p := range []int{7, 11, 2, 3, 13, 5} {
		s.Insert(p)
	}
	it := s.End()
	for it.Prev() {
		fmt.Println(it.Value())
	}
	// Output:
	// 13
	// 11
	// 7
	// 5
	// 3
	// 2
}

func ExampleIterator_Value() {
	s := set.New[int]()
	for _, p := range []int{7, 11, 2, 3, 13, 5} {
		s.Insert(p)
	}
	it := s.Begin()

	// it.Value() // would panic, iterator is not initialised

	it.Next()
	fmt.Println(it.Value())
	// Output: 2
}

func ExampleSet_LowerBound_exact() {
	s := set.New[int]()
	for _, p := range []int{7, 11, 2, 3, 13, 5} {
		s.Insert(p)
	}
	it := s.LowerBound(11) // points to 11
	for it.Next() {
		fmt.Println(it.Value())
	}
	// Output:
	// 11
	// 13
}

func ExampleSet_LowerBound_missing() {
	s := set.New[int]()
	for _, p := range []int{7, 11, 2, 3, 13, 5} {
		s.Insert(p)
	}
	it := s.LowerBound(4) // 4 is missing, points to 5
	for it.Prev() {
		fmt.Println(it.Value())
	}
	// Output:
	// 3
	// 2
}

func ExampleSet_UpperBound() {
	s := set.New[int]()
	for _, p := range []int{7, 11, 2, 3, 13, 5} {
		s.Insert(p)
	}
	it := s.UpperBound(5) // points to 7 (first element greater than 5)
	for it.Next() {
		fmt.Println(it.Value())
	}
	// Output:
	// 7
	// 11
	// 13
}

func ExampleSet_Len() {
	s := set.New[int]()
	for _, p := range []int{7, 11, 2, 3, 13, 5} {
		s.Insert(p)
	}

	fmt.Println(s.Len())

	s.Insert(5)
	s.Delete(100)

	fmt.Println(s.Len()) // set is unchanged, length should be the same

	// Output:
	// 6
	// 6
}
