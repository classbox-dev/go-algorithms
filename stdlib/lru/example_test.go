package lru_test

import (
	"fmt"
	"hsecode.com/stdlib/lru"
)

func ExampleCache_Get() {
	c := lru.New(3) // cache of size 3
	c.Put(1, 100)
	c.Put(2, 200)
	c.Put(3, 300)

	fmt.Println(c.Get(1))
	fmt.Println(c.Get(5))
	// Output:
	// 100 true
	// 0 false
}

func ExampleCache_Put() {
	c := lru.New(2) // cache of size 3

	c.Put(1, 100)
	c.Put(2, 200)
	c.Put(3, 300)

	fmt.Println(c.Get(1)) // -> not found, the key was replaced

	c.Get(2) // "use" the key

	c.Put(4, 400) // remove `3` because `2` was used recently.
	c.Put(2, 200) // "use" `2`
	c.Put(5, 500) // remove `4`, `2` is still there

	fmt.Println(c.Get(2)) // -> found, `2` was used recently
	fmt.Println(c.Get(4)) // not found

	// Output:
	// 0 false
	// 200 true
	// 0 false
}
