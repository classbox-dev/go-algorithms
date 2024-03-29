package matrix_test

import (
	"reflect"
	"testing"

	"hsecode.com/stdlib-tests/internal/utils"
	"hsecode.com/stdlib/matrix"
)

func TestUnit__RowsCols(t *testing.T) {
	n, m := utils.Rand.Int()%256, utils.Rand.Int()%256
	mat := matrix.New[int](n, m)
	if mat.Rows != n {
		t.Fatal("unexpected value of Rows")
	}
	if mat.Cols != m {
		t.Fatal("unexpected value of Cols")
	}
}

func TestUnit__OutOfRange(t *testing.T) {
	mat := matrix.New[int](64, 64)
	indices := []struct{ a, b int }{
		{64, 63},
		{64, 65},
		{10, 64},
		{64, 10},
		{10, -1},
		{-1, 10},
		{100, 100},
	}
	sum := 0
	for _, idx := range indices {
		utils.ExpectedPanic(t, "out-of-range Set did not panic", func() {
			mat.Set(idx.a, idx.b, idx.a+idx.b)
		})
		utils.ExpectedPanic(t, "out-of-range Get did not panic", func() {
			sum += mat.Get(idx.a, idx.b)
		})
	}
	utils.Use(sum)
}

func TestUnit__FillAndRead(t *testing.T) {
	n, m := utils.Rand.Int()%256, utils.Rand.Int()%256
	mat := matrix.New[int](n, m)
	expected := make([]int, n*m)
	result := make([]int, n*m)

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			v := utils.Rand.Int()
			expected = append(expected, v)
			mat.Set(i, j, v)
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			result = append(result, mat.Get(i, j))
		}
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatal("filled the matrix with n*m values, but could not read them all back")
	}
}

func TestPerf__Iterate(t *testing.T) {
	n, m := 2096000, 4
	mat := matrix.New[int](n, m)

	sum := 0
	for c := 0; c < 30; c++ {
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				mat.Set(i, j, 23)
			}
		}
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				mat.Set(j, i, 11)
			}
		}
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				sum += mat.Get(i, j)
			}
		}
	}
	utils.Use(sum)
}
