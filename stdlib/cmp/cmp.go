// Package cmp defines Min and Max functions for a variable number of ValueType arguments.
package cmp

import "github.com/cheekybits/genny/generic"

//go:generate genny -in=$GOFILE -out=int/dont_edit.go gen "ValueType=int"

// ValueType is a generic type for Min and Max (imported from github.com/cheekybits/genny/generic). Concrete types have to support comparison via operators < and >.
//
// The package contains subpackages where ValueType is automatically replaced with concrete types.
//
// The int subpackage is required for tests. Make sure to include the following comment in your source code:
//
//	//go:generate genny -in=$GOFILE -out=int/dont_edit.go gen "ValueType=int"
//
type ValueType generic.Number

// Min returns the minimum of a variable number of ValueType arguments. Panics if called with no arguments.
func Min(values ...ValueType) ValueType {
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

// Max returns the maximum of a variable number of ValueType arguments. Panics if called with no arguments.
func Max(values ...ValueType) ValueType {
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
