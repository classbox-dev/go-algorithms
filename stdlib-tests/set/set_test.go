package radix_test

import (
	"fmt"
	"hsecode.com/stdlib-tests/internal/utils"
	"hsecode.com/stdlib/set"
	"reflect"
	"runtime"
	"sort"
	"testing"
)

type Direction uint8

const (
	hardLimit = 200000
)

const (
	Forward Direction = iota
	Backward
)

type I int

func (v I) Less(other set.Element) bool {
	return int(v) < int(other.(I))
}
func (v I) Equal(other set.Element) bool {
	return int(v) == int(other.(I))
}

func consume(it *set.Iterator, dir Direction, f func(v int) (stop bool)) {
	var iter func() bool
	iterFunc := ""
	if dir == Forward {
		iter = func() bool { return it.Next() }
		iterFunc = "Next()"
	} else {
		iter = func() bool { return it.Prev() }
		iterFunc = "Prev()"
	}
	i := 0
	for iter() && i < hardLimit {
		v := it.Value()
		if v == nil {
			e := fmt.Sprintf("Iterator.Value() returned <nil> after .%s==true", iterFunc)
			panic(e)
		}
		if f(int(v.(I))) {
			break
		}
		i++
	}
	if i >= hardLimit {
		panic("Infinite loop detected on consuming an iterator")
	}
}

func TestUnit__Insert(t *testing.T) {
	ss := set.New()

	a, b := 57, 1001
	elems := utils.RangeShuffled(a, b)

	n := 0
	for _, e := range elems {
		ss.Insert(I(e))

		if n+1 != ss.Len() {
			t.Fatal("Len() was not increased by 1 on inserting an element")
		}
		n = ss.Len()

		found := false
		consume(ss.Begin(), Forward, func(v int) (stop bool) {
			if v == e {
				found = true
				stop = true
			}
			return
		})

		if !found {
			t.Fatal("Begin() did not yield an inserted element")
		}

		if v, ok := ss.Find(I(e)); !ok || int(v.(I)) != e {
			t.Fatal("Find() could not find an inserted element")
		}
	}
}

func TestUnit__Delete(t *testing.T) {
	ss := set.New()

	a, b := 57, 1001
	elems := utils.RangeShuffled(a, b)

	for _, e := range elems {
		ss.Insert(I(e))
	}

	n := ss.Len()
	for _, e := range elems {
		if ok := ss.Delete(I(e)); !ok {
			t.Fatal("Delete() returned false for an existing element")
		}

		if _, ok := ss.Find(I(e)); ok {
			t.Fatal("Find() returned a deleted element")
		}
		if ss.Len() != n-1 {
			t.Fatal("Len() was not descreased by 1 on deleting an element")
		}
		n = ss.Len()

		consume(ss.Begin(), Forward, func(v int) (stop bool) {
			if v == e {
				t.Fatal("Begin()+Next() yielded a deleted element")
			}
			return
		})

		consume(ss.End(), Backward, func(v int) (stop bool) {
			if v == e {
				t.Fatal("End()+Prev() yielded a deleted element")
			}
			return
		})

		if ok := ss.Delete(I(e)); ok {
			t.Fatal("Delete() returned true for a missing element")
		}
		if n != ss.Len() {
			t.Fatal("Len() changed after deleting a missing element")
		}
	}
}

func TestUnit__Find(t *testing.T) {
	ss := set.New()

	a, b := 57, 1001
	ss.Insert(I(b))
	elems := utils.RangeReversed(a, b)

	for _, e := range elems {
		ss.Insert(I(e))
		for needle := e; needle <= e+1; needle++ {
			v, ok := ss.Find(I(needle))
			if !ok {
				t.Fatal("Find() could not find an existing element")
			}
			if int(v.(I)) != needle {
				t.Fatal("Find() returned an invalid element")
			}
		}
		if _, ok := ss.Find(I(e - 1)); ok {
			t.Fatal("Find() returned success for a missing element")
		}
	}
}

func TestUnit__IterateAll(t *testing.T) {
	ss := set.New()

	a, b := 57, 1001
	elems := utils.RangeShuffled(a, b)
	for _, e := range elems {
		ss.Insert(I(e))
	}

	sort.Ints(elems)

	result := make([]int, 0)
	consume(ss.Begin(), Forward, func(v int) (stop bool) {
		result = append(result, v)
		return
	})
	if !reflect.DeepEqual(result, elems) {
		t.Fatalf("Begin()+Next(): yielded elements do not match expected sequence")
	}

	sort.Sort(sort.Reverse(sort.IntSlice(elems)))
	result = result[:0]
	consume(ss.End(), Backward, func(v int) (stop bool) {
		result = append(result, v)
		return
	})
	if !reflect.DeepEqual(result, elems) {
		t.Fatalf("End()+Prev(): yielded elements do not match expected sequence")
	}
}

func TestUnit__EmptyIterators(t *testing.T) {
	ss := set.New()

	a, b := 57, 1001
	elems := utils.RangeShuffled(a, b)
	for _, e := range elems {
		ss.Insert(I(e))
	}

	result := make([]int, 0)
	consume(ss.Begin(), Backward, func(v int) (stop bool) {
		result = append(result, v)
		return
	})
	if !reflect.DeepEqual(result, []int{}) {
		t.Fatalf("Begin()+Prev(): yielded elements do not match expected sequence")
	}

	result = result[:0]
	consume(ss.End(), Forward, func(v int) (stop bool) {
		result = append(result, v)
		return
	})
	if !reflect.DeepEqual(result, []int{}) {
		t.Fatalf("End()+Next(): yielded elements do not match expected sequence")
	}
}

func TestUnit__DeleteSomeAndIterate(t *testing.T) {
	ss := set.New()

	a, b := 57, 1001
	elems := utils.RangeShuffled(a, b)

	for _, e := range elems {
		ss.Insert(I(e))
	}
	for _, e := range elems {
		if e%2 == 0 {
			ss.Delete(I(e))
		}
	}

	expected := make([]int, 0, len(elems)/2+1)
	for _, e := range elems {
		if e%2 == 1 {
			expected = append(expected, e)
		}
	}
	sort.Ints(expected)

	result := make([]int, 0, len(elems)/2+1)

	// Forward
	consume(ss.Begin(), Forward, func(v int) (stop bool) {
		result = append(result, v)
		return
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Begin()+Next(): yielded elements do not match expected sequence")
	}

	sort.Sort(sort.Reverse(sort.IntSlice(expected)))
	result = result[:0]

	// Backward
	consume(ss.End(), Backward, func(v int) (stop bool) {
		result = append(result, v)
		return
	})
	if !reflect.DeepEqual(result, expected) {
		fmt.Println(result)
		fmt.Println(expected)
		t.Fatalf("End()+Prev(): yielded elements do not match expected sequence")
	}
}

func TestUnit__InsertDuplicatesAndIterate(t *testing.T) {
	ss := set.New()

	a, b := 57, 1001
	elems := utils.RangeShuffled(a, b)

	for i := 0; i < 6; i++ {
		for _, e := range elems {
			if ok := ss.Insert(I(e)); ok && i > 0 {
				t.Fatal("Insert() returned true for a duplicated element")
			}
		}
		utils.Rand.Shuffle(len(elems), func(i, j int) {
			elems[i], elems[j] = elems[j], elems[i]
		})
	}

	if ss.Len() != len(elems) {
		t.Fatal("Len() does not match the number of unique elements")
	}

	sort.Ints(elems)
	result := make([]int, 0)
	consume(ss.Begin(), Forward, func(v int) (stop bool) {
		result = append(result, v)
		return
	})

	if !reflect.DeepEqual(result, elems) {
		t.Fatalf("Begin()+Next(): yielded elements do not match expected sequence")
	}
}

type testCase struct {
	start     int
	direction Direction
	expected  []int
}

func TestUnit__UpperBound(t *testing.T) {
	ss := set.New()

	elems := utils.RangeShuffled(1, 100)
	for _, v := range elems {
		ss.Insert(I(v * 10))
	}

	cases := []testCase{
		{-1, Forward, []int{10, 20, 30, 40, 50, 60}},
		{9, Forward, []int{10, 20, 30, 40, 50, 60}},
		{599, Forward, []int{600, 610, 620, 630, 640, 650}},
		{501, Forward, []int{510, 520, 530, 540, 550, 560}},
		{100, Forward, []int{110, 120, 130, 140, 150, 160, 170}},
		{970, Forward, []int{980, 990}},
		{2000, Forward, []int{}},
	}
	for _, cs := range cases {
		i := 0
		result := make([]int, 0)
		consume(ss.UpperBound(I(cs.start)), cs.direction, func(v int) (stop bool) {
			if i >= len(cs.expected) {
				stop = true
				return
			}
			result = append(result, v)
			i++
			return
		})
		if !reflect.DeepEqual(result, cs.expected) {
			t.Fatalf("UpperBound() yielded unexpected sequence of elements")
		}
	}
}

func TestUnit__LowerBound(t *testing.T) {
	ss := set.New()

	elems := utils.RangeShuffled(1, 100)
	for _, v := range elems {
		ss.Insert(I(v * 10))
	}

	cases := []testCase{
		{-1, Forward, []int{10, 20, 30, 40, 50, 60}},
		{9, Forward, []int{10, 20, 30, 40, 50, 60}},
		{599, Forward, []int{600, 610, 620, 630, 640, 650}},
		{501, Forward, []int{510, 520, 530, 540, 550, 560}},
		{100, Forward, []int{100, 110, 120, 130, 140, 150, 160, 170}},
		{970, Forward, []int{970, 980, 990}},
		{2000, Forward, []int{}},
	}
	for _, cs := range cases {
		i := 0
		result := make([]int, 0)
		consume(ss.LowerBound(I(cs.start)), cs.direction, func(v int) (stop bool) {
			if i >= len(cs.expected) {
				stop = true
				return
			}
			result = append(result, v)
			i++
			return
		})
		if !reflect.DeepEqual(result, cs.expected) {
			t.Fatalf("LowerBound() yielded unexpected sequence of elements")
		}
	}
}

func TestUnit__UninitialisedIterator(t *testing.T) {
	ss := set.New()
	elems := utils.RangeShuffled(1, 100)
	for _, v := range elems {
		ss.Insert(I(v * 10))
	}

	its := []*set.Iterator{
		ss.Begin(),
		ss.End(),
		ss.LowerBound(I(100)),
		ss.UpperBound(I(323)),
	}
	msg := "Expected panic when calling .Value() on an iterator before .Next() or .Prev()"
	for _, it := range its {
		utils.ExpectedPanic(t, msg, func() {
			it.Value()
		})
	}
}

func rangeRandomUnique(a, b, length int) []int {
	h := make(map[int]struct{})
	output := make([]int, 0, length)
	for i := 0; i < length; {
		v := a + utils.Rand.Intn(b-a)
		if _, ok := h[v]; ok {
			continue
		}
		output = append(output, v)
		h[v] = struct{}{}
		i++
	}
	return output
}

func TestUnit__Random(t *testing.T) {

	runtime.GC()

	ss := set.New()
	magnituge := 30000

	limits := [2]int{50000000, 50000100}
	l := 0

	sign := -1
	for i := 0; i < 10; i++ {
		a := limits[l]
		b := a + sign*magnituge
		limits[l] = b + sign*100

		if a > b {
			a, b = b, a
		}

		values := rangeRandomUnique(a, b, magnituge)
		for _, v := range values {
			ss.Insert(I(v))
			ss.Insert(I(v))
		}

		for j := 0; j < 100; j++ {
			e := I(values[len(values)-1-j])
			ss.Delete(e)
			ss.Delete(e)
		}

		values = values[:len(values)-100]
		sort.Ints(values)

		if sign > 0 {
			realMax := values[len(values)-1]
			consume(ss.End(), Backward, func(v int) (stop bool) {
				if v != realMax {
					t.Fatal("End() does not point to the maximum value")
				}
				stop = true
				return
			})
		} else {
			realMin := values[0]
			consume(ss.Begin(), Forward, func(v int) (stop bool) {
				if v != realMin {
					t.Fatal("Begin() does not point to the minimum value")
				}
				stop = true
				return
			})
		}

		runLength := 5

		for c := 0; c < 8; c++ {
			idx := utils.Rand.Intn(len(values) - 100)

			if _, ok := ss.Find(I(values[idx])); !ok {
				t.Fatal("Find() could not find an existing element")
			}

			cnt := 0
			result := make([]int, 0)
			consume(ss.LowerBound(I(values[idx])), Forward, func(v int) (stop bool) {
				if cnt >= runLength {
					stop = true
					return
				}
				result = append(result, v)
				cnt++
				return
			})
			if !reflect.DeepEqual(result, values[idx:idx+runLength]) {
				t.Fatal("LowerBound() from a random existing value yielded unexpected sequence")
			}
		}

		for c := 0; c < 8; c++ {
			idx := utils.Rand.Intn(len(values) - 100)
			cnt := 0
			result := make([]int, 0)
			consume(ss.UpperBound(I(values[idx])), Forward, func(v int) (stop bool) {
				if cnt >= runLength {
					stop = true
					return
				}
				result = append(result, v)
				cnt++
				return
			})

			if !reflect.DeepEqual(result, values[idx+1:idx+1+runLength]) {
				t.Fatal("UpperBound() from a random existing value yielded unexpected sequence")
			}
		}

		l = (l + 1) % 2
		sign *= -1
	}
}

func TestPerf__Random(t *testing.T) {

	runtime.GC()

	ss := set.New()
	magnituge := 30000

	limits := [2]int{50000000, 50000100}
	l := 0

	sign := -1
	for i := 0; i < 14; i++ {
		a := limits[l]
		b := a + sign*magnituge*10
		limits[l] = b + sign*100

		if a > b {
			a, b = b, a
		}

		values := rangeRandomUnique(a, b, magnituge)
		for i, v := range values {
			ss.Insert(I(v))
			ss.Insert(I(v))
			ss.End()
			ss.Begin()

			if i > 10 {
				idx := utils.Rand.Intn(i)
				ss.Find(I(values[idx]))
				ss.LowerBound(I(values[idx]))
				ss.UpperBound(I(values[idx]))
			}
		}

		for j := 0; j < 100; j++ {
			e := I(values[len(values)-1-j])
			ss.Delete(e)
			ss.Delete(e)
			ss.End()
			ss.Begin()
		}

		l = (l + 1) % 2
		sign *= -1
	}
}
