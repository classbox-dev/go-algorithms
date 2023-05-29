package primes_test

import (
	"fmt"
	"hsecode.com/stdlib-tests/v2/internal/utils"
	"hsecode.com/stdlib/v2/math"
	"math/big"
	"testing"
)

func TestUnit__Panic(t *testing.T) {
	for _, i := range []int{0, -1, -10} {
		msg := fmt.Sprintf("NthPrime(%v) did not panic", i)
		utils.ExpectedPanic(t, msg, func() {
			math.NthPrime(i)
		})
	}
}

func TestUnit__IsPrime(t *testing.T) {
	N := 10000
	ps := map[int]int{}

	for i := 1; i <= N; i += 2 {
		p := math.NthPrime(i)
		if !big.NewInt(int64(p)).ProbablyPrime(0) {
			t.Fatalf("NthPrime(%d)=%d is not prime", i, p)
		}
		ii, ok := ps[p]
		if ok {
			t.Fatalf("%d returned for NthPrime(%d) and NthPrime(%d)", p, ii, i)
		}
		ps[p] = i
	}
}

func TestUnit__CheckPoints(t *testing.T) {
	ns := []int{793, 4536, 54729, 636450, 10056324}
	ps := []int{6079, 43541, 675559, 9546007, 180494837}
	for i, n := range ns {
		if math.NthPrime(n) != ps[i] {
			t.Fatal("NthPrime() returned unexpected value for a checkpoint")
		}
	}
}

func TestPerf__CheckPoints(t *testing.T) {
	ns := []int{793, 4536, 54729, 636450, 2672370, 10056325}
	for _, n := range ns {
		math.NthPrime(n)
	}
}
