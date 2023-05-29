package maxqueue

import (
	"errors"
	"golang.org/x/exp/constraints"
)

type maxStack[E constraints.Ordered] struct {
	elems []E
	maxs  []E
}

func (q *maxStack[E]) Push(value E) {
	var newMax E

	n := len(q.maxs)
	if q.Empty() || value > q.maxs[n-1] {
		newMax = value
	} else {
		newMax = q.maxs[n-1]
	}
	q.maxs = append(q.maxs, newMax)
	q.elems = append(q.elems, value)
}

func (q *maxStack[E]) Pop() (E, error) {
	if q.Empty() {
		var zero E
		return zero, errors.New("stack is empty")
	}
	n := len(q.elems)
	r := q.elems[n-1]
	q.elems = q.elems[:n-1]
	q.maxs = q.maxs[:n-1]
	return r, nil
}

func (q *maxStack[E]) Max() (E, error) {
	if q.Empty() {
		var zero E
		return zero, errors.New("stack is empty")
	}
	n := len(q.maxs)
	return q.maxs[n-1], nil
}

func (q *maxStack[E]) Empty() bool {
	return len(q.elems) == 0
}

// MaxQueue is a FIFO queue that allows fast queries for the maximal element.
type MaxQueue[E constraints.Ordered] struct {
	head maxStack[E]
	tail maxStack[E]
}

// New creates an instance of MaxQueue
func New[E constraints.Ordered]() *MaxQueue[E] {
	return new(MaxQueue[E])
}

// Push inserts an element to the queue tail in amortised constant time.
func (q *MaxQueue[E]) Push(value E) {
	q.tail.Push(value)
}

// Pop removes an element from the queue head in amortised constant time. Returns an error if the queue is empty.
func (q *MaxQueue[E]) Pop() (E, error) {
	if q.head.Empty() {
		for !q.tail.Empty() {
			v, _ := q.tail.Pop()
			q.head.Push(v)
		}
	}

	if q.head.Empty() {
		var zero E
		return zero, errors.New("queue is empty")
	}

	return q.head.Pop()
}

// Max returns the maximal element in constant time.
// Returns an error if the queue is empty.
func (q *MaxQueue[E]) Max() (E, error) {
	var m1, m2 E

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

	var zero E
	return zero, errors.New("queue is empty")
}
