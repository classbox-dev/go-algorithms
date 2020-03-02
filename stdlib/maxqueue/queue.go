package maxqueue

import (
	"errors"
	"github.com/cheekybits/genny/generic"
)

//go:generate genny -in=$GOFILE -out=int/dont_edit.go gen "ValueType=int"

// ValueType is a generic type for MaxQueue values (imported from github.com/cheekybits/genny/generic).
//
// The package contains subpackages where ValueType is automatically replaced with concrete types.
//
// The int subpackage is required for tests. Make sure to include the following comment in your source code:
//
//	//go:generate genny -in=$GOFILE -out=int/dont_edit.go gen "ValueType=int"
//
type ValueType generic.Number

type maxStack struct {
	elems []ValueType
	maxs  []ValueType
}

func (q *maxStack) Push(value ValueType) {
	var newMax ValueType

	n := len(q.maxs)
	if q.Empty() || value > q.maxs[n-1] {
		newMax = value
	} else {
		newMax = q.maxs[n-1]
	}
	q.maxs = append(q.maxs, newMax)
	q.elems = append(q.elems, value)
}

func (q *maxStack) Pop() (ValueType, error) {
	if q.Empty() {
		return 0, errors.New("stack is empty")
	}
	n := len(q.elems)
	r := q.elems[n-1]
	q.elems = q.elems[:n-1]
	q.maxs = q.maxs[:n-1]
	return r, nil
}

func (q *maxStack) Max() (ValueType, error) {
	if q.Empty() {
		return 0, errors.New("stack is empty")
	}
	n := len(q.maxs)
	return q.maxs[n-1], nil
}

func (q *maxStack) Empty() bool {
	return len(q.elems) == 0
}

// MaxQueue is a FIFO queue that allows fast queries for the maximal element.
type MaxQueue struct {
	head maxStack
	tail maxStack
}

// New creates an instance of MaxQueue
func New() *MaxQueue {
	return new(MaxQueue)
}

// Push inserts an element to the queue tail in amortised constant time.
func (q *MaxQueue) Push(value ValueType) {
	q.tail.Push(value)
}

// Pop removes an element from the queue head in amortised constant time. Returns an error if the queue is empty.
func (q *MaxQueue) Pop() (ValueType, error) {
	if q.head.Empty() {
		for !q.tail.Empty() {
			v, _ := q.tail.Pop()
			q.head.Push(v)
		}
	}

	if q.head.Empty() {
		return 0, errors.New("queue is empty")
	}

	return q.head.Pop()
}

// Max returns the maximal element in constant time.
// Returns an error if the queue is empty.
func (q *MaxQueue) Max() (ValueType, error) {
	var m1, m2 ValueType

	m1, err1 := q.head.Max()
	m2, err2 := q.tail.Max()

	if err1 == nil && err2 == nil {
		if m1 > m2 {
			return m1, nil
		}
		return m2, nil
	} else if err1 == nil && err2 != nil {
		return m1, nil
	} else if err1 != nil && err2 == nil {
		return m2, nil
	}

	return 0, errors.New("queue is empty")
}
