package primes

import (
	"hsecode.com/algorithms/stdlib/primes"
	"math/big"
	"testing"
)

func TestUnit__IsPrime(t *testing.T) {
	N := 5000
	ps := map[uint32]int{}

	for i := 1; i <= N; i++ {
		p := primes.NthPrime(uint32(i))
		if !big.NewInt(int64(p)).ProbablyPrime(0) {
			t.Fatalf("NthPrime(%d)=%d is not prime!", i, p)
		}
		ii, ok := ps[p]
		if ok {
			t.Fatalf("%d returned for NthPrime(%d) and NthPrime(%d)", p, ii, i)
		}
		ps[p] = i
	}
}

func TestUnit__CheckPoints(t *testing.T) {
	ns := []uint32{793, 4536, 54729, 636450, 2672370, 16456324}
	ps := []uint32{6079, 43541, 675559, 9546007, 44194303, 303986209}
	for i, n := range ns {
		if primes.NthPrime(n) != ps[i] {
			t.Fatal("Checkpoint failed")
		}
	}
}

func TestPerf__CheckPoints(t *testing.T) {
	ns := []uint32{793, 4536, 54729, 636450, 2672370, 5000000}
	for _, n := range ns {
		primes.NthPrime(n)
	}
}
