package xlist

import (
	"container/list"
)

type sublist struct {
	head   *list.Element
	length int
}

// Sort performs in-place sorting of the given linked list according to the provided less function.
// "In-place" here means that the resulting list can only contain the originally allocated elements with their original values.
//
// The reference implementation uses a non-recursive variant of merge-sort.
func Sort(data *list.List, less func(a, b *list.Element) bool) {
	parts := make([]sublist, 0, 16)
	i := 0
	for e := data.Front(); e != nil; e = e.Next() {
		if (i % 2) == 0 {
			for n := len(parts); n > 3 && parts[n-2].length <= parts[n-1].length; n = len(parts) {
				p1, p2 := parts[n-2], parts[n-1]
				parts = parts[:n-2]
				parts = append(parts, merge(data, p1, p2, less))
			}
			length := 2
			en := e.Next()
			if en == nil {
				length--
			} else if less(en, e) {
				data.MoveBefore(en, e)
				e = en
			}
			parts = append(parts, sublist{e, length})
		}
		i++
	}

	for n := len(parts); n > 1; n = len(parts) {
		p1, p2 := parts[n-2], parts[n-1]
		parts = parts[:n-2]
		parts = append(parts, merge(data, p1, p2, less))
	}
}

func merge(data *list.List, left, right sublist, less func(a, b *list.Element) bool) sublist {
	dummy := data.InsertBefore(nil, left.head)
	current := dummy
	i, j := 0, 0
	for l, r := left.head, right.head; i < left.length && j < right.length; {
		if less(l, r) {
			tmp := l.Next()
			data.MoveAfter(l, current)
			l = tmp
			i++
		} else {
			tmp := r.Next()
			data.MoveAfter(r, current)
			r = tmp
			j++
		}
		current = current.Next()
	}
	head := dummy.Next()
	data.Remove(dummy)
	return sublist{head, left.length + right.length}
}
