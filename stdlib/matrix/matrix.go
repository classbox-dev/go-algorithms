package matrix

import (
	"github.com/cheekybits/genny/generic"
)

//go:generate genny -in=$GOFILE -out=int/dont_edit.go gen "ValueType=int"

// ValueType is a generic type for Matrix values (imported from github.com/cheekybits/genny/generic).
//
// The package contains subpackages where ValueType is automatically replaced with concrete types.
//
// The int subpackage is required for tests. Make sure to include the following comment in your source code:
//
//	//go:generate genny -in=$GOFILE -out=int/dont_edit.go gen "ValueType=int"
//
type ValueType generic.Type

// Matrix implements row-major 2D array of ValueType
type Matrix struct {
	values []ValueType
	Rows   int // number of rows
	Cols   int // number of columns
}

// New creates a new Matrix with n rows and m columns
func New(n, m int) *Matrix {
	mat := new(Matrix)
	mat.values = make([]ValueType, n*m)
	mat.Rows, mat.Cols = n, m
	return mat
}

func (M *Matrix) idx(i, j int) int {
	if i >= M.Rows || i < 0 || j >= M.Cols || j < 0 {
		panic("out-of-range")
	}
	return M.Cols*i + j
}

// Get returns a matrix element by indices (i, j). Panics if the indices are out of range.
func (M *Matrix) Get(i, j int) ValueType {
	return M.values[M.idx(i, j)]
}

// Set assigns a value of matrix element by indices (i, j). Panics if the indices are out of range.
func (M *Matrix) Set(i, j int, v ValueType) {
	M.values[M.idx(i, j)] = v
}
