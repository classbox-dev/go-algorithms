package islands

import matrix "hsecode.com/stdlib/matrix/int"

// Count returns the number of islands (contiguous non-zero regions) on the given grid.
// The function can (and probably should) mutate the grid. Reference solution uses depth-first search.
func Count(grid *matrix.Matrix) int {
	m := grid.Rows
	if m == 0 {
		return 0
	}
	n, res := grid.Cols, 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid.Get(i, j) != 0 {
				res++
				dfs(grid, i, j, m, n)
			}
		}
	}
	return res
}

func dfs(grid *matrix.Matrix, i, j, m, n int) {
	grid.Set(i, j, 0)
	if i-1 >= 0 && grid.Get(i-1, j) == 1 {
		dfs(grid, i-1, j, m, n)
	}
	if j-1 >= 0 && grid.Get(i, j-1) == 1 {
		dfs(grid, i, j-1, m, n)
	}
	if i+1 < m && grid.Get(i+1, j) == 1 {
		dfs(grid, i+1, j, m, n)
	}
	if j+1 < n && grid.Get(i, j+1) == 1 {
		dfs(grid, i, j+1, m, n)
	}
}
