package matrix

// Matrix implements row-major 2D array of ValueType
type Matrix[V any] struct {
	values []V
	Rows   int // number of rows
	Cols   int // number of columns
}

// New creates a new Matrix with n rows and m columns
func New[V any](n, m int) *Matrix[V] {
	mat := new(Matrix[V])
	mat.values = make([]V, n*m)
	mat.Rows, mat.Cols = n, m
	return mat
}

func (M *Matrix[V]) idx(i, j int) int {
	if i >= M.Rows || i < 0 || j >= M.Cols || j < 0 {
		panic("out-of-range")
	}
	return M.Cols*i + j
}

// Get returns a matrix element by indices (i, j). Panics if the indices are out of range.
func (M *Matrix[V]) Get(i, j int) V {
	return M.values[M.idx(i, j)]
}

// Set assigns a value of matrix element by indices (i, j). Panics if the indices are out of range.
func (M *Matrix[V]) Set(i, j int, v V) {
	M.values[M.idx(i, j)] = v
}
