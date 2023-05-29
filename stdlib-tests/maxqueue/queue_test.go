package maxqueue_test

import (
	"hsecode.com/stdlib-tests/v2/internal/utils"
	"hsecode.com/stdlib/v2/maxqueue"
	"testing"
)

func TestUnit__QueueOrder(t *testing.T) {
	q := maxqueue.New[int]()

	n := 100
	input := make([]int, 0, n)
	output := make([]int, 0, n)

	for i := 0; i < n; i++ {
		input = append(input, utils.Rand.Int()%1000)
	}
	for i := 0; i < n; i++ {
		q.Push(input[i])
	}
	for i := 0; i < n; i++ {
		v, _ := q.Pop()
		output = append(output, v)
	}
	for i := 0; i < n; i++ {
		if output[i] != input[i] {
			t.Fatalf("order of Pop()'d elements differ from the order they were Push()'d into the queue")
		}
	}
}

func TestUnit__PopEmpty(t *testing.T) {
	q := maxqueue.New[int]()
	_, err := q.Pop()
	if err == nil {
		t.Fatal("expected an error when calling Pop() on an empty queue")
	}
}

func TestUnit__PopSingle(t *testing.T) {
	q := maxqueue.New[int]()
	q.Push(100)

	v, err := q.Pop()
	if err != nil || v != 100 {
		t.Fatal("failed to Pop() on a queue with a single element")
	}
}

func TestUnit__PushAndPop(t *testing.T) {
	q := maxqueue.New[int]()

	n := 100
	for i := 0; i < n; i++ {
		q.Push(i * i)
	}

	for i := 0; i < n; i++ {
		if _, err := q.Pop(); err != nil {
			t.Fatal("unexpected error on Pop()")
		}
	}

	_, err := q.Pop()
	if err == nil {
		t.Fatalf("queue is not empty after %d Push'es and %d Pop's", n, n)
	}
}

func TestUnit__MaxEmpty(t *testing.T) {
	q := maxqueue.New[int]()
	_, err := q.Max()
	if err == nil {
		t.Fatal("expected an error when calling Max() on an empty queue")
	}
}

func TestUnit__MaxSingle(t *testing.T) {
	q := maxqueue.New[int]()
	q.Push(100)

	v, err := q.Max()
	if err != nil || v != 100 {
		t.Fatal("failed to Max() on a queue with a single element")
	}
}

func TestUnit__Max(t *testing.T) {
	q := maxqueue.New[int]()

	q.Push(100)
	q.Push(10)

	v, err := q.Max()
	if err != nil || v != 100 {
		t.Fatalf("incorrect Max() value")
	}

	if _, err := q.Pop(); err != nil {
		t.Fatal("unexpected error on Pop")
	}

	v, err = q.Max()
	if err != nil || v != 10 {
		t.Fatalf("incorrect Max() value")
	}

	q.Push(1)

	v, err = q.Max()
	if err != nil || v != 10 {
		t.Fatalf("incorrect Max() value")
	}

	q.Push(1000)

	v, err = q.Max()
	if err != nil || v != 1000 {
		t.Fatalf("incorrect Max() value")
	}
}

func TestUnit__MaxAll(t *testing.T) {
	q := maxqueue.New[int]()
	n := 100

	for i := 0; i < n; i++ {
		q.Push(i * i)
	}

	v, err := q.Max()
	if err != nil || v != (n-1)*(n-1) {
		t.Fatalf("incorrect Max() value")
	}

	for i := 0; i < n; i++ {
		q.Push(i*i + 1)
	}

	v, err = q.Max()
	if err != nil || v != (n-1)*(n-1)+1 {
		t.Fatalf("incorrect Max() value")
	}

	for i := 0; i < 2*n-1; i++ {
		if _, err := q.Pop(); err != nil {
			t.Fatal("unexpected error on Pop")
		}
	}

	v, err = q.Max()
	if err != nil || v != (n-1)*(n-1)+1 {
		t.Fatalf("incorrect Max() value")
	}

	q.Push(0)

	v, err = q.Max()
	if err != nil || v != (n-1)*(n-1)+1 {
		t.Fatalf("incorrect Max() value")
	}

	q.Push(n*n + 10)

	v, err = q.Max()
	if err != nil || v != n*n+10 {
		t.Fatalf("incorrect Max() value")
	}
}

func TestPerf__RandomOperations(t *testing.T) {
	q := maxqueue.New[int]()
	for i := 0; i < 2000000; i++ {
		q.Push(utils.Rand.Int() % 10000)
		_, err := q.Max()
		if err == nil && (utils.Rand.Int()%1000) > 800 {
			//noinspection GoUnhandledErrorResult
			q.Pop()
			//noinspection GoUnhandledErrorResult
			q.Max()
		}
	}
}
