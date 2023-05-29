package fulltext_test

import (
	"fmt"
	"hsecode.com/stdlib/v2/strings/fulltext"
)

func Example() {
	index := fulltext.New([]string{
		"this is the house that jack built", // doc #0
		"this is the rat that ate the malt", // doc #1
	})
	fmt.Println(index.Search(""))
	fmt.Println(index.Search("in the house that jack built"))
	fmt.Println(index.Search("malt rat"))
	fmt.Println(index.Search("is this the"))

	// Output:
	// []
	// []
	// [1]
	// [0 1]
}
