package xlist_test

import (
	"container/list"
	"fmt"
	"hsecode.com/stdlib/xlist"
)

var xs = []int{13, 2, 29, 3, 7, 23, 31}

func Example_sort() {
	// Create linked list from slice
	data := list.New()
	for _, x := range xs {
		data.PushBack(x)
	}

	xlist.Sort(data, func(a, b *list.Element) bool {
		return a.Value.(int) < b.Value.(int)
	})

	// Iterate over sorted list
	for e := data.Front(); e != nil; e = e.Next() {
		fmt.Printf("%v ", e.Value)
	}
	// Output:
	// 2 3 7 13 23 29 31
}
