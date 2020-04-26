package set

import (
	"math/bits"
	"math/rand"
)

const (
	maxLevel   = 25
	levelsMask = ^((uint32(1) << maxLevel) - 1)
)

type fatNext [maxLevel + 1]*listElement

// Element interface has to be implemented by Set elements.
type Element interface {
	// Less returns true if the element is less than the other
	Less(other Element) bool
	// Equal returns true if the element is equivalent to the other
	Equal(other Element) bool
}

type listElement struct {
	prev  *listElement
	next  fatNext
	value Element
}

// Set represents an ordered collection of unique elements
type Set struct {
	head   *listElement
	length int
}

// Iterator is a stateful iterator pointing to a set element or an imaginary "past-the-end" element.
//
// See the diagram visualising set iterators and iteration order: https://hsecode.com/.static/set-iterator.png
type Iterator struct {
	head    *listElement
	curr    *listElement
	started bool
}

// Value returns an element the iterator is pointing to.
// Panics if neither Next() nor Prev() was called beforehand.
// The return value can only be trusted if the preceding Next() or Prev() returned true.
func (it *Iterator) Value() Element {
	if !it.started {
		panic("iterator is not initialised with Next() call")
	}
	return it.curr.value
}

// Next advances the iterator to the next greater element or,
// if called for the first time, simply initialises the iterator without advancing it.
// Returns true on success, or false if the iteration is finished.
func (it *Iterator) Next() bool {
	if it.curr.next[0] == nil {
		return false
	}
	it.curr = it.curr.next[0]
	it.started = true
	return true
}

// Prev advances the iterator to the next smaller element.
// Returns true on success, or false if the iteration is finished.
func (it *Iterator) Prev() bool {
	if !it.started {
		it.started = true
		return it.curr != it.head
	}
	if it.curr.prev == it.head {
		return false
	}
	it.started = true
	it.curr = it.curr.prev
	return true
}

// New creates an empty set
func New() *Set {
	s := new(Set)
	s.head = new(listElement)
	return s
}

func (s *Set) lookup(e Element, stack *fatNext, upper bool) (ok bool) {
	currElem := s.head
	for i := maxLevel; i >= 0; i-- {
		prev, curr := currElem, currElem.next[i]
		for curr != nil {
			isLess := curr.value.Less(e)
			isEqual := !isLess && curr.value.Equal(e)
			if (!upper && isLess) || (upper && (isLess || isEqual)) {
				prev = curr
				curr = curr.next[i]
			} else if isEqual {
				ok = true
				break
			} else {
				break
			}
		}
		stack[i] = prev
		currElem = prev
	}
	return
}

// Insert adds a new element to the set.
// Does nothing if an equivalent element is already in the set.
// Returns true if the actual insertion happens.
// The running time is O(log N) for N elements.
func (s *Set) Insert(e Element) bool {
	var stack fatNext
	if ok := s.lookup(e, &stack, false); ok {
		return false
	}
	prev := stack[0]

	le := new(listElement)
	le.value = e

	le.prev = prev
	oldNext := prev.next[0]
	if oldNext != nil {
		oldNext.prev = le
	}
	prev.next[0] = le
	le.next[0] = oldNext

	rnd := rand.Uint32() | levelsMask
	topLevel := bits.TrailingZeros32(rnd)

	for i := 1; i <= topLevel; i++ {
		oldNext := stack[i].next[i]
		stack[i].next[i] = le
		le.next[i] = oldNext
	}
	s.length++
	return true
}

// Delete removes an element equivalent to the given one from the set.
// Does nothing if there is no equivalent element in the set.
// Returns true if the actual deletion happens.
// The running time is O(log N) for N elements.
func (s *Set) Delete(e Element) bool {
	var stack fatNext
	if ok := s.lookup(e, &stack, false); !ok {
		return false
	}
	elem := stack[0].next[0]
	nextNext := elem.next[0]
	if nextNext != nil {
		nextNext.prev = stack[0]
	}
	for i, p := range &stack {
		if p.next[i] == elem {
			nextElem := elem.next[i]
			p.next[i] = nextElem
		}
	}
	s.length--
	return true
}

// Find returns an element from the set that is equivalent to the given one, or nil if such element is not present.
// The ok result indicates whether the element was found.
// The running time is O(log N) for N elements.
func (s *Set) Find(e Element) (Element, bool) {
	var stack fatNext
	if ok := s.lookup(e, &stack, false); ok {
		return stack[0].next[0].value, true
	}
	return nil, false
}

// Begin() returns an iterator pointing to the first (minimum) element of the set.
// The running time is O(log N) for N elements.
func (s *Set) Begin() *Iterator {
	return &Iterator{curr: s.head, head: s.head}
}

// End() returns an iterator pointing to the element following the last (maximum) element.
// The running time is O(log N) for N elements.
func (s *Set) End() *Iterator {
	currElem := s.head
	for i := maxLevel; i >= 0; i-- {
		prev, curr := currElem, currElem.next[i]
		for curr != nil {
			prev = curr
			curr = curr.next[i]
		}
		currElem = prev
	}
	return &Iterator{curr: currElem, head: s.head}
}

// LowerBound returns an iterator pointing to the first element not less than the given one.
// The running time is O(log N) for N elements.
func (s *Set) LowerBound(e Element) *Iterator {
	var stack fatNext
	_ = s.lookup(e, &stack, false)
	return &Iterator{curr: stack[0], head: s.head}
}

// UpperBound returns an iterator pointing to the first element greater than the given one.
// The running time is O(log N) for N elements.
func (s *Set) UpperBound(e Element) *Iterator {
	var stack fatNext
	_ = s.lookup(e, &stack, true)
	return &Iterator{curr: stack[0], head: s.head}
}

// Len returns the number of elements in the set. The running time is O(1).
func (s *Set) Len() int {
	return s.length
}
