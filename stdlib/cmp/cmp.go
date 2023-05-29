package cmp

import "golang.org/x/exp/constraints"

// Min returns the minimum of a variable number of arguments. Panics if called with no arguments.
func Min[T constraints.Ordered](values ...T) T {
	if len(values) < 1 {
		panic("Min requires at least one argument")
	}
	m := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] < m {
			m = values[i]
		}
	}
	return m
}

// Max returns the maximum of a variable number of arguments. Panics if called with no arguments.
func Max[T constraints.Ordered](values ...T) T {
	if len(values) < 1 {
		panic("Min requires at least one argument")
	}
	m := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] > m {
			m = values[i]
		}
	}
	return m
}
