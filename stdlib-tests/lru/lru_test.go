package lru_test

import (
	"fmt"
	"hsecode.com/stdlib/lru"
	"reflect"
	"strings"
	"testing"
)

const (
	getOp = 1
	putOp = 2
)

type Action struct {
	op   int
	args []int
}

type Case struct {
	cap      int
	actions  []Action
	expected []int
}

func sprintResult(result []int) string {
	ps := make([]string, 0, len(result))
	for _, x := range result {
		if x == -1 {
			ps = append(ps, "notfound")
		} else {
			ps = append(ps, fmt.Sprintf("%v", x))
		}
	}
	return fmt.Sprintf("[%v]", strings.Join(ps, ", "))
}

func reprCase(cs Case) (string, string) {
	var b strings.Builder
	_, err := fmt.Fprintf(&b, "New(%d); ", cs.cap)
	if err != nil {
		panic(err)
	}

	n := len(cs.actions)
	ie := 0

	for i, ac := range cs.actions {
		switch ac.op {
		case getOp:
			_, err := fmt.Fprintf(&b, "Get(%d)", ac.args[0])
			if err != nil {
				panic(err)
			}
			ie++
		case putOp:
			_, err := fmt.Fprintf(&b, "Put(%d, %d)", ac.args[0], ac.args[1])
			if err != nil {
				panic(err)
			}
		}
		if i != (n - 1) {
			b.WriteString("; ")
		}
	}
	return b.String(), sprintResult(cs.expected)
}

func runCase(cs Case) []int {
	cache := lru.New(cs.cap)
	results := make([]int, 0)

	for _, ac := range cs.actions {
		switch ac.op {
		case getOp:
			v, ok := cache.Get(ac.args[0])
			if !ok {
				v = -1
			}
			results = append(results, v)
		case putOp:
			cache.Put(ac.args[0], ac.args[1])
		}
	}
	return results
}

var cases = []Case{
	{
		1, []Action{
			{getOp, []int{0}},
		}, []int{-1},
	},
	{
		2, []Action{
			{putOp, []int{1, 2}},
			{putOp, []int{2, 3}},
			{getOp, []int{1}},
			{getOp, []int{2}},
		}, []int{2, 3},
	},
	{
		2, []Action{
			{putOp, []int{1, 1}},
			{putOp, []int{2, 2}},
			{putOp, []int{3, 3}},
			{getOp, []int{1}},
		}, []int{-1},
	},
	{
		2, []Action{
			{putOp, []int{1, 1}},
			{putOp, []int{2, 2}},
			{getOp, []int{1}},
			{putOp, []int{3, 3}},
			{getOp, []int{1}},
			{getOp, []int{2}},
			{getOp, []int{3}},
		}, []int{1, 1, -1, 3},
	},
	{
		3, []Action{
			{putOp, []int{1, 1}},
			{putOp, []int{2, 2}},
			{putOp, []int{1, 10}},
			{getOp, []int{1}},
		}, []int{10},
	},
	{
		3, []Action{
			{putOp, []int{1, 1}},
			{putOp, []int{2, 2}},
			{putOp, []int{3, 3}},
			{putOp, []int{1, 10}},
			{putOp, []int{4, 4}},
			{getOp, []int{1}},
			{getOp, []int{2}},
			{getOp, []int{3}},
			{getOp, []int{4}},
		}, []int{10, -1, 3, 4},
	},
	{
		3, []Action{
			{putOp, []int{1, 1}},
			{putOp, []int{2, 2}},
			{putOp, []int{3, 3}},
			{putOp, []int{1, 1}},
			{putOp, []int{4, 4}},
			{getOp, []int{1}},
			{getOp, []int{2}},
			{getOp, []int{3}},
			{getOp, []int{4}},
		}, []int{1, -1, 3, 4},
	},
}

func TestUnit__Cases(t *testing.T) {

	for _, cs := range cases {
		results := runCase(cs)
		if !reflect.DeepEqual(results, cs.expected) {
			input, output := reprCase(cs)
			t.Fatalf("Test failed:\nInput: %v\nExpected output: %v\nActual output: %v", input, output, sprintResult(results))
		}
	}
}

func TestUnit__LRUMany(t *testing.T) {
	N := 100000
	cache := lru.New(N)
	for i := 0; i < N; i++ {
		cache.Put(i, i)
		if i%10 == 0 {
			for j := 0; j < 100; j++ {
				if j%2 == 0 {
					cache.Get(j)
				} else {
					cache.Put(j, j+10000)
				}
			}
		}
	}
	for j := 0; j < 100; j++ {
		if _, ok := cache.Get(j); !ok {
			t.Fatalf("Get or Put have not prevented cache element from removal")
		}
	}
}

func TestPerf__Get(t *testing.T) {
	N := 100000
	cache := lru.New(N)
	for i := 0; i < N; i++ {
		cache.Put(i, i)
	}
	for i := 0; i < N; i++ {
		if _, ok := cache.Get(i); !ok {
			t.Fatalf("Failed to get an inserted value")
		}
		if _, ok := cache.Get(N - i - 1); !ok {
			t.Fatalf("Failed to get an inserted value")
		}
	}
}

func TestPerf__Put(t *testing.T) {
	N := 100000
	for try := 0; try < 10; try++ {
		cache := lru.New(N)
		for i := 0; i < N; i++ {
			cache.Put(i, i)
		}
		for i := 0; i < N/2; i++ {
			cache.Put(i, 0)
		}
		for i := N / 2; i < N; i++ {
			cache.Put(i, i)
		}
		for i := 0; i < N; i++ {
			if _, ok := cache.Get(i); !ok {
				t.Fatalf("Failed to get an inserted value")
			}
			if _, ok := cache.Get(N - i - 1); !ok {
				t.Fatalf("Failed to get an inserted value")
			}
		}
	}
}
