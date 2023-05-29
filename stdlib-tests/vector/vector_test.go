package vector_test

import (
	"hsecode.com/stdlib-tests/v2/internal/utils"
	"hsecode.com/stdlib/v2/vector"
	"reflect"
	"runtime"
	"runtime/debug"
	"testing"
	"unsafe"
)

type Element = int

func toSlice[E any](arr *vector.Vector[E]) []E {
	sl := make([]E, 0)
	for i := 0; i < arr.Len; i++ {
		sl = append(sl, arr.Get(i))
	}
	return sl
}

func TestUnit__GetOutOfRange(t *testing.T) {
	h := vector.New[int](3)
	for _, idx := range []int{-1, 10} {
		msg := "did not get panic on out-of-range Get"
		utils.ExpectedPanic(t, msg, func() {
			h.Get(idx)
		})
	}
}

func TestUnit__SetOutOfRange(t *testing.T) {
	h := vector.New[int](3)
	for _, idx := range []int{-1, 10} {
		msg := "did not get panic on out-of-range Set"
		utils.ExpectedPanic(t, msg, func() {
			h.Set(idx, Element(100))
		})
	}
}

func TestUnit__InsertOutOfRange(t *testing.T) {
	h := vector.New[int](3)
	for _, idx := range []int{-1, 10} {
		msg := "did not get panic on out-of-range Insert"
		utils.ExpectedPanic(t, msg, func() {
			h.Insert(idx, 42)
		})
	}
}

func TestUnit__DeleteOutOfRange(t *testing.T) {
	h := vector.New[int](3)
	for _, idx := range []int{-1, 10} {
		msg := "did not get panic on out-of-range Delete"
		utils.ExpectedPanic(t, msg, func() {
			h.Delete(idx)
		})
	}
}

func TestUnit__PopFromEmpty(t *testing.T) {
	h := vector.New[int](0)
	msg := "did not panic on Pop from empty vector"
	utils.ExpectedPanic(t, msg, func() {
		h.Pop()
	})
}

func TestUnit__Get(t *testing.T) {
	h := vector.New[int](0)
	for i := 0; i < 30; i++ {
		e := Element(i * i)
		h.Push(e)
	}
	for i := 30 - 1; i >= 0; i-- {
		if h.Get(i) != i*i {
			t.Fatal("Get(i) returned unexpected value")
		}
	}
}

func TestUnit__Set(t *testing.T) {
	h := vector.New[int](0)
	for i := 0; i < 30; i++ {
		e := Element(i * i)
		h.Push(e)
	}
	for i := 30 - 1; i >= 0; i-- {
		newValue := Element(utils.Rand.Uint64() % 1024)
		h.Set(i, newValue)
		if h.Get(i) != newValue {
			t.Fatal("Get(i) after Set(i, x) did not return `x`")
		}
	}
}

func TestUnit__InsertMiddle(t *testing.T) {
	expected := make([]Element, 0, 30)
	h := vector.New[int](0)
	for i := 0; i < 30; i++ {
		e := Element(i * i)
		h.Push(e)
		expected = append(expected, e)
	}
	for i := 0; i < 30; i++ {
		newValue := Element(utils.Rand.Uint64() % 1024)
		h.Insert(i, newValue)
		expected = append(expected[:i], append([]Element{newValue}, expected[i:]...)...)
		result := toSlice(h)
		if !reflect.DeepEqual(result, expected) {
			t.Fatal("unexpected vector content after a series of middle inserts")
		}
	}
}

func TestUnit__Pop(t *testing.T) {
	expected := make([]Element, 0, 30)
	h := vector.New[int](0)
	for i := 0; i < 30; i++ {
		e := Element(i * i)
		h.Push(e)
		expected = append(expected, e)
	}
	for i := 30 - 1; i >= 0; i-- {
		if h.Pop() != expected[i] {
			t.Fatal("Pop returned unexpected value")
		}
	}
}

func TestUnit__Push(t *testing.T) {
	expected := make([]Element, 0, 30)
	h := vector.New[int](0)
	for i := 0; i < 30; i++ {
		e := Element(i * i)
		h.Push(e)
		expected = append(expected, e)
		result := toSlice(h)
		if !reflect.DeepEqual(result, expected) {
			t.Fatalf("unexpected vector content after a series of %d pushes", i+1)
		}
	}
}

func TestUnit__Len(t *testing.T) {
	h := vector.New[int](0)
	if h.Len != 0 {
		t.Error("Len() of empty vector is not 0")
	}
	h = vector.New[int](10)
	if h.Len != 0 {
		t.Error("Len() of empty vector is not 0")
	}
	h.Push(42)
	if h.Len != 1 {
		t.Error("Len() of one-element vector is not 1")
	}
	h.Delete(0)
	if h.Len != 0 {
		t.Error("Len() of empty vector is not 0")
	}
}

func TestUnit__LenGrowing(t *testing.T) {
	h := vector.New[int](0)
	for i := 0; i < 1000; i++ {
		l := h.Len
		h.Push(Element(i * i))
		if h.Len != l+1 {
			t.Fatal("value of Len() before and after Push must differ by 1")
		}
	}
}

func TestUnit__LenShrinking(t *testing.T) {
	h := vector.New[int](0)
	for i := 0; i < 1000; i++ {
		h.Push(Element(i * i))
	}
	for i := 0; i < 1000; i++ {
		l := h.Len
		h.Pop()
		if h.Len != l-1 {
			t.Fatal("value of Len() before and after Pop must differ by 1", i)
		}
	}
}

func TestUnit__DeleteLeft(t *testing.T) {
	h := vector.New[int](0)
	for i := 0; i < 10; i++ {
		h.Push(Element(i * i))
	}
	for i := 0; i < 10; i++ {
		h.Delete(0)
	}
	if h.Len != 0 {
		t.Fatal("vector is not empty after deleting all elements")
	}
}

func TestPerf__PushAllPopAll(t *testing.T) {
	n := int((utils.Rand.Uint32() % 1000000) + 1000000)
	leakMargin := int64(4096) // bytes

	arr := vector.New[int](0)

	leak := utils.MemoryLeak(func() {
		for i := 0; i < n; i++ {
			arr.Push(10000)
		}
		for i := 0; i < n; i++ {
			arr.Pop()
		}
		if arr.Len != 0 {
			t.Fatalf("Len() is not 0 after inserting and removing %d elements", n)
		}
	})
	if leak > leakMargin {
		t.Fatalf("Memory leak: %d bytes has left after inserting and removing %d elements. Expected at most %d", leak, n, leakMargin)
	}

	arr.Push(42) // prevents GC from optimising the vector out
}

func TestPerf__InsertAllDeleteAll(t *testing.T) {
	n := int((utils.Rand.Uint32() % 10000) + 10000)
	leakMargin := int64(4096) // bytes

	arr := vector.New[int](0)

	leak := utils.MemoryLeak(func() {
		for i := 0; i < n; i++ {
			arr.Insert(0, 10000)
		}
		for i := 0; i < n; i++ {
			arr.Delete(arr.Len / 2)
		}
		if arr.Len != 0 {
			t.Fatalf("Len() is not 0 after inserting and removing %d elements", n)
		}
	})
	if leak > leakMargin {
		t.Fatalf("Memory leak: %d bytes has left after inserting and removing %d elements. Expected at most %d", leak, n, leakMargin)
	}

	arr.Push(42) // prevents GC from optimising the vector out
}

func TestPerf__InsertOnly(t *testing.T) {
	n := int((utils.Rand.Uint32() % 10000) + 10000)

	elemSize := int64(unsafe.Sizeof(Element(0)))
	leakMargin := int64(n) * elemSize * 4 // bytes

	arr := vector.New[int](0)

	leak := utils.MemoryLeak(func() {
		for i := 0; i < n; i++ {
			arr.Push(10000)
		}
		if arr.Len != n {
			t.Fatalf("Len() is %d after pushing %d elements", arr.Len, n)
		}
	})
	if leak > leakMargin {
		t.Fatalf("%d bytes is used after pushing %d elements. Expected at most %d", leak, n, leakMargin)
	}

	arr.Push(42) // prevents GC from optimising the vector out
}

func TestPerf__RandomPushPop(t *testing.T) {
	debug.SetGCPercent(-1) // Disable GC
	runtime.GC()

	arr := vector.New[int](0)

	n := 1000000
	for i := 0; i < n/2; i++ {
		arr.Push(Element(i))
	}
	for i := 0; i < n; i++ {
		rnd := utils.Rand.NormFloat64() + 0.3
		if rnd > 0. {
			arr.Push(Element(i))
		} else if arr.Len > n/4 {
			arr.Pop()
		}
	}
}

func TestPerf__RandomInsertDelete(t *testing.T) {
	debug.SetGCPercent(-1) // Disable GC
	runtime.GC()

	arr := vector.New[int](0)

	n := 50000
	for i := 0; i < n/2; i++ {
		arr.Push(Element(i))
	}
	for i := 0; i < n; i++ {
		rnd := utils.Rand.NormFloat64() + 0.3
		idx := utils.Rand.Intn(arr.Len)
		if rnd > 0. {
			arr.Insert(idx, Element(i))
		} else if arr.Len > n/4 {
			arr.Delete(idx)
		}
	}
}
