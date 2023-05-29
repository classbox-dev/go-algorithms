package vector

// Vector is an implementation of dynamic array with fast insertion to and deletion from the end
type Vector[E any] struct {
	slice []E
	Len   int // actual number of occupied elements in the underlying slice
}

func (a *Vector[E]) init(cap int) *Vector[E] {
	a.resize(cap)
	return a
}

func (a *Vector[E]) resize(cap int) {
	elems := make([]E, cap)
	for i := 0; i < a.Len; i++ {
		elems[i] = a.slice[i]
	}
	a.slice = elems
}

// New creates a new empty Vector with a given capacity
func New[E any](cap int) *Vector[E] { return new(Vector[E]).init(cap) }

// Get retrieves an element by index idx
func (a *Vector[E]) Get(idx int) E {
	if idx < 0 || idx >= a.Len {
		panic("index error")
	}
	return a.slice[idx]
}

// Set writes the given element to index idx
func (a *Vector[E]) Set(idx int, x E) {
	if idx < 0 || idx >= a.Len {
		panic("index error")
	}
	a.slice[idx] = x
}

// Insert inserts new element before index idx
func (a *Vector[E]) Insert(idx int, x E) {
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
func (a *Vector[E]) Delete(idx int) {
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
func (a *Vector[E]) Push(x E) {
	a.Insert(a.Len, x)
}

// Pop deletes and returns an element from the end of the vector
func (a *Vector[E]) Pop() E {
	v := a.Get(a.Len - 1)
	a.Delete(a.Len - 1)
	return v
}
