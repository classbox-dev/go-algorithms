package islands_test

import (
	"fmt"
	"hsecode.com/stdlib/graph/islands"
	matrix "hsecode.com/stdlib/matrix/int"
)

// Three islands
var rows = [][]int{
	{1, 0, 1},
	{0, 0, 1},
	{1, 1, 0},
}

func Example() {
	grid := matrix.New(3, 3)

	// Fill in the matrix
	for i, row := range rows {
		for j, v := range row {
			grid.Set(i, j, v)
		}
	}

	fmt.Println(islands.Count(grid))
	// Output: 3
}
