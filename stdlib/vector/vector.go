package vector

import (
	"github.com/cheekybits/genny/generic"
)

//go:generate genny -in=$GOFILE -out=int/dont_edit.go gen "ValueType=int"

// ValueType is a generic type of a vector element (imported from github.com/cheekybits/genny/generic).
//
// The package contains subpackages where ValueType is automatically replaced with concrete types.
//
// The int subpackage is required for tests. Make sure to include the following comment in your source code:
//
//	//go:generate genny -in=$GOFILE -out=int/dont_edit.go gen "ValueType=int"
//
type ValueType generic.Type

// Vector is an implementation of dynamic array with fast insertion to and deletion from the end
type Vector struct {
	slice []ValueType
	Len   int // actual number of occupied elements in the underlying slice
}

func (a *Vector) init(cap int) *Vector {
	a.resize(cap)
	return a
}

func (a *Vector) resize(cap int) {
	elems := make([]ValueType, cap, cap)
	for i := 0; i < a.Len; i++ {
		elems[i] = a.slice[i]
	}
	a.slice = elems
}

// New creates a new empty Vector with a given capacity
func New(cap int) *Vector { return new(Vector).init(cap) }

// Get retrieves an element by index idx
func (a *Vector) Get(idx int) ValueType {
	if idx < 0 || idx >= a.Len {
		panic("index error")
	}
	return a.slice[idx]
}

// Set writes the given element to index idx
func (a *Vector) Set(idx int, x ValueType) {
	if idx < 0 || idx >= a.Len {
		panic("index error")
	}
	a.slice[idx] = x
}

// Insert inserts new element before index idx
func (a *Vector) Insert(idx int, x ValueType) {
	if idx < 0 || idx > a.Len {
		panic("insert index is out of range")
	}
	if a.Len+1 > len(a.slice) {
		if len(a.slice) > 2 {
			a.resize(4 * len(a.slice) / 3)
		} else {
			a.resize(len(a.slice) + 1)
		}
	}
	for j := a.Len - 1; j >= idx; j-- {
		a.slice[j+1] = a.slice[j]
	}
	a.slice[idx] = x
	a.Len++
}

// Delete removes an element at index idx
func (a *Vector) Delete(idx int) {
	if idx < 0 || idx >= a.Len {
		panic("insert index is out of range")
	}
	for j := a.Len - 1; j > idx; j-- {
		a.slice[j-1] = a.slice[j]
	}
	a.Len--
	if 2*a.Len < len(a.slice) && len(a.slice) > 1 {
		a.resize(3 * len(a.slice) / 4)
	}
}

// Push inserts the given element to the end of the vector
func (a *Vector) Push(x ValueType) {
	a.Insert(a.Len, x)
}

// Pop deletes and returns an element from the end of the vector
func (a *Vector) Pop() ValueType {
	v := a.Get(a.Len - 1)
	a.Delete(a.Len - 1)
	return v
}