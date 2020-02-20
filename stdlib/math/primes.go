package math

import (
	"math"
)

var primesTo100 = []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97}

// NthPrime returns n-th (1-based) prime number.
// Panics if n is less than 1.
func NthPrime(n int) int {
	if n < 1 {
		panic("`n` must be greater than 0")
	}
	if n <= len(primesTo100) {
		return primesTo100[n-1]
	}

	fN := float64(n)
	lnN := math.Log(fN)
	lnLnN := math.Log(lnN)

	fupper := math.Ceil(fN*lnN + fN*lnLnN)
	upper := int(fupper) | 1
	sieveLimit := int(math.Sqrt(fupper + 1))

	sieve := make([]uint8, (upper-3)>>1+1)

	for i := 3; i <= sieveLimit; i += 2 {
		ii := (i - 3) >> 1
		if sieve[ii] == 1 {
			continue
		}
		for j := i * i; j < upper; j += 2 * i {
			sieve[(j-3)>>1] = 1
		}
	}
	for i, cnt := 0, 1; i < len(sieve); i++ {
		cnt += int(sieve[i] ^ uint8(1))
		if cnt == n {
			return i<<1 + 3
		}
	}
	panic("This point should not be reached.")
}
