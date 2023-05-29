package bitset_test

import (
	"fmt"
	"hsecode.com/stdlib/v2/bitset"
)

func ExampleBitset_Set() {
	bs := bitset.New(3)
	_ = bs.Set(1, true)

	value, _ := bs.Test(1) // true
	fmt.Println(value)

	err := bs.Set(3, false)
	fmt.Println(err) // out of range

	// Output:
	// true
	// out of range
}

func ExampleBitset_Test() {
	bs := bitset.New(3)
	_ = bs.Set(1, true)

	value, _ := bs.Test(0) // false
	fmt.Println(value)

	value, _ = bs.Test(1) // true
	fmt.Println(value)

	_, err := bs.Test(3)
	fmt.Println(err) // out of range

	// Output:
	// false
	// true
	// out of range
}

func ExampleBitset_Flip() {
	bs := bitset.New(3)
	_ = bs.Set(0, true)
	_ = bs.Set(1, false)

	bs.Flip()
	v1, _ := bs.Test(0) // false
	v2, _ := bs.Test(1) // true
	fmt.Println(v1, v2)

	// Output:
	// false true
}

func ExampleBitset_Reset() {
	bs := bitset.New(3)
	_ = bs.Set(1, true)

	bs.Reset()
	value, _ := bs.Test(1) // false
	fmt.Println(value)

	// Output:
	// false
}

func ExampleBitset_Any() {
	bs := bitset.New(3)

	fmt.Println(bs.Any()) // false
	_ = bs.Set(1, true)
	fmt.Println(bs.Any()) // true

	// Output:
	// false
	// true
}

func ExampleBitset_All() {
	bs := bitset.New(2)

	fmt.Println(bs.All()) // false
	_ = bs.Set(0, true)
	_ = bs.Set(1, true)
	fmt.Println(bs.All()) // true

	// Output:
	// false
	// true
}

func ExampleBitset_Count() {
	bs := bitset.New(3)

	fmt.Println(bs.Count()) // 0

	_ = bs.Set(0, true)
	fmt.Println(bs.Count()) // 1

	_ = bs.Set(1, true)
	fmt.Println(bs.Count()) // 2

	// Output:
	// 0
	// 1
	// 2
}
