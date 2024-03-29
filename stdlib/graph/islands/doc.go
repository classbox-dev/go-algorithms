// Package islands implements counting the number of non-zero regions on an integer grid.
/*
Definitions

The grid is represented by an int matrix. The following import is required:

  import "hsecode.com/stdlib/matrix"

The matrix values stand for land (non-zero) and water (zero).
An island is surrounded by water and is formed by connecting adjacent lands horizontally or vertically.
You may assume that all four edges of the grid are surrounded by water.

Hints

You do not have to use "hsecode.com/stdlib/graph" to implement this package.
The function can (and probably should) mutate the grid.
*/
package islands
