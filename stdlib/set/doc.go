// Package set implements an ordered collection of unique elements with logarithmic operations.
/*
Because of generic elements, the set can also be used as ordered dictionary with range lookups.

Iterators

The following diagram demonstrates the behaviour of an iterator that initially points to the value "E":

	                           UpperBound("C")
	Begin()  LowerBound("C")   LowerBound("D")            End()
	     │             │       │                            │
	     ▼             ▼       ▼                            ▼
	   ┌───┐   ┌───┐   ┌───┐ ┌───┐    ┌───┐        ┌───┐  ┌ ─ ┐
	   │ A │   │ B │   │ C │ │ E │    │ F │    ... │ Z │
	   └───┘   └───┘   └───┘ └───┘    └───┘        └───┘  └ ─ ┘
	      ◀─────  ◀─────      │▲ ──────▶  ─────▶
	       Prev()  Prev()     └┘   Next() Next()
	                       Next()

The semantics is optimised for using Next() and Prev() as conditions in for-loops:

	it := s.Begin()
	for it.Next() {
		it.Value() // in-order iteration through the whole set forward
	}

	it := s.End()
	for it.Prev() {
		it.Value() // in-order iteration through the whole set backward
	}

Implementation

The reference implementation uses skip-list with the maximum number of pointer levels hard-coded to 26.
Therefore the set can hold up to 33.6 million elements without losing its logarithmic expected running time for insertions, deletions, and lookups.
*/
package set
