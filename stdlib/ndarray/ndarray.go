package ndarray

// NDArray is a prototype of a row-major n-dimensional array. Does NOT hold the actual array, only the dimensions.
type NDArray struct {
	shape []int
}

// New creates an NDArray with given dimentions. Panics if any of the dimensions are negative.
func New(shape ...int) *NDArray {
	for _, i := range shape {
		if i < 0 {
			panic("invalid shape")
		}
	}
	return &NDArray{shape}
}

// Idx returns a linearised index for a given n-dimentional index. Panics if the index is out of range or has an invalid dimension.
func (nda *NDArray) Idx(indicies ...int) int {
	ni, ns := len(indicies), len(nda.shape)
	if ni != ns {
		panic("invalid dimension")
	}
	for i, idx := range indicies {
		if idx < 0 || (idx+1) > nda.shape[i] {
			panic("index is out-of-range")
		}
	}
	total, s := 1, 0
	for i := 1; i <= ni; i++ {
		s += total * indicies[ni-i]
		total *= nda.shape[ni-i]
	}
	return s
}
