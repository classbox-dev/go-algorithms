package islands_test

import (
	"fmt"
	xislands "hsecode.com/stdlib-tests/v2/internal/islands"
	"hsecode.com/stdlib-tests/v2/internal/utils"
	"hsecode.com/stdlib/v2/graph/islands"
	"hsecode.com/stdlib/v2/matrix"
	"strings"
	"testing"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func gridRepr(grid *matrix.Matrix[int]) string {
	var repr strings.Builder
	for i := 0; i < grid.Rows; i++ {
		for j := 0; j < grid.Cols; j++ {
			if _, err := fmt.Fprintf(&repr, fmt.Sprintf("%v ", grid.Get(i, j))); err != nil {
				panic(err)
			}
		}
		if _, err := fmt.Fprintf(&repr, "\n"); err != nil {
			panic(err)
		}
	}
	return repr.String()
}

func TestUnit__Basic(t *testing.T) {
	for i := 1; i < 128; i += 2 {
		for j := 1; j < 128; j += 2 {
			h, w := utils.Rand.Intn(i), utils.Rand.Intn(j)
			grid := matrix.New[int](h, w)
			gridRef := matrix.New[int](h, w)
			gridRerp := matrix.New[int](h, w)
			for x := 0; x < h; x++ {
				xw := w
				for xw > 0 {
					rr := utils.Rand.Uint64()
					n := min(64, xw)
					for p := 0; p < n; p++ {
						grid.Set(x, xw-1, int(rr&1))
						gridRef.Set(x, xw-1, int(rr&1))
						gridRerp.Set(x, xw-1, int(rr&1))
						rr >>= 1
						xw--
					}
				}
			}
			expected := xislands.Count(gridRef)
			result := islands.Count(grid)
			if expected != result {
				repr := gridRepr(gridRerp)
				t.Fatalf("Expected %d, got %d\nGrid:\n%v", expected, result, repr)
			}
		}
	}
}

func TestUnit__Perf(t *testing.T) {
	for i := 1; i < 128; i += 1 {
		for j := 1; j < 128; j += 1 {
			grid := matrix.New[int](i, j)
			for x := 0; x < i; x++ {
				xw := j
				for xw > 0 {
					rr := utils.Rand.Uint64()
					n := min(64, xw)
					for p := 0; p < n; p++ {
						grid.Set(x, xw-1, int(rr&1))
						rr >>= 1
						xw--
					}
				}
			}
			utils.Use(islands.Count(grid))
		}
	}
}
