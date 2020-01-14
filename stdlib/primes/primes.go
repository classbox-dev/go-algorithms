package primes

import (
	"math"
)

var primesTo100 = []uint32{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97}

// NthPrime returns nth (base-1) prime number
func NthPrime(n uint32) uint32 {
	if n == 0 {
		panic("`n` must be greater than 0")
	}
	if n <= uint32(len(primesTo100)) {
		return primesTo100[n-1]
	}

	fN := float64(n)
	lnN := math.Log(fN)
	lnLnN := math.Log(lnN)

	fupper := math.Ceil(fN*lnN + fN*lnLnN)
	upper := uint32(fupper)
	sieveLimit := int(math.Sqrt(fupper + 1))

	upper = upper | uint32(1)

	sieve := make([]uint8, (upper-3)>>1+1)

	for i := 3; i <= sieveLimit; i += 2 {
		ii := (i - 3) >> 1
		if sieve[ii] == 1 {
			continue
		}
		for j := i * i; j < int(upper); j += 2 * i {
			sieve[(j-3)>>1] = 1
		}
	}
	for i, cnt := uint32(0), 1; i < uint32(len(sieve)); i++ {
		cnt += int(sieve[i] ^ uint8(1))
		if cnt == int(n) {
			return i<<1 + 3
		}
	}
	panic("This point should not be reached.")
}
