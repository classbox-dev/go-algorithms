package treeserialise_test

import (
	xtree "hsecode.com/stdlib-tests/internal/tree"
	"hsecode.com/stdlib-tests/internal/utils"
	"hsecode.com/stdlib/tree"
	"reflect"
	"strconv"
	"testing"
)

func TestUnit__DecodeInvalidValue(t *testing.T) {
	_, err := tree.Decode([]string{"sdagfasgd"})
	if err == nil {
		t.Error("error expected for non-integer input")
	}
}

func TestUnit__DecodeInvalidStructure(t *testing.T) {
	inputs := [][]string{
		{"10", "1", "1", "nil", "3", "3", "nil", "nil", "nil", "100", "nil", "nil", "100", "234", "23"},
		{"1", "nil", "nil", "2", "4", "32", "5"},
		{"nil", "234"},
	}
	for _, data := range inputs {
		original := make([]string, len(data))
		copy(original, data)
		_, err := tree.Decode(data)
		if err == nil {
			t.Fatalf("error expected for non-nil node being children of nil: %v", original)
		}
	}
}

func TestUnit__Decode(t *testing.T) {

	inputs := [][]string{
		{"10", "1", "1", "nil", "3", "3", "nil", "nil", "nil", "100", "nil", "nil", "100", "nil", "nil"},
		{"1", "2", "3", "nil", "4", "32", "5", "nil", "nil", "4", "nil", "5", "nil", "nil"},
		{},
		{"nil"},
		{"35"},
	}

	for i := 10; i < 100; i++ {
		data := make([]string, 0, i)
		for j := 0; j < i; j++ {
			data = append(data, strconv.Itoa(utils.Rand.Int()%20000))
			data = append(data, strconv.Itoa(utils.Rand.Int()%20000))
		}
		inputs = append(inputs, data)
	}

	for _, data := range inputs {
		data = xtree.Normalise(data)
		original := make([]string, len(data))
		copy(original, data)

		T, err := tree.Decode(data)
		if err != nil {
			t.Fatal(err)
		}
		out := xtree.Encode(T)
		if !reflect.DeepEqual(original, out) {
			t.Fatalf("decoding failed:\nExpected: %v\nOutput: %v", data, out)
		}
	}
}

func TestUnit__Encode(t *testing.T) {

	inputs := [][]string{
		{"10", "1", "1", "nil", "3", "3", "nil", "nil", "nil", "100", "nil", "nil", "100", "nil", "nil"},
		{"1", "2", "3", "nil", "4", "32", "5", "nil", "nil", "4", "nil", "5", "nil", "nil"},
		{},
		{"nil"},
		{"35"},
	}

	for i := 10; i < 100; i++ {
		data := make([]string, 0, i)
		for j := 0; j < i; j++ {
			data = append(data, strconv.Itoa(utils.Rand.Int()%20000))
		}
		inputs = append(inputs, data)
	}

	for _, data := range inputs {
		data = xtree.Normalise(data)

		original := make([]string, len(data))
		copy(original, data)

		T, err := xtree.Decode(data)
		if err != nil {
			t.Fatal(err)
		}
		out := T.Encode()
		if !reflect.DeepEqual(original, out) {
			t.Fatalf("encoding failed:\nTree: %v\nOutput: %v", data, out)
		}
	}
}

func TestUnit__RandomEncodeDecodeWithNil(t *testing.T) {
	for i := 1; i < 100; i += 5 {
		var T *tree.Tree
		for j := 0; j < i; j++ {
			xtree.Insert(&T, utils.Rand.Int())
		}
		expected := xtree.Encode(T)
		result := T.Encode()
		if !reflect.DeepEqual(expected, result) {
			t.Fatalf("Encode produced an invalid result on a random tree with %d nodes", i)
		}
		xT, err := tree.Decode(expected)
		if err != nil {
			t.Fatalf("error on decoding a valid tree: %v", err)
		}
		if !reflect.DeepEqual(xtree.Encode(xT), expected) {
			t.Fatalf("Decode produced an invalid result on a random tree with %d nodes", i)
		}
	}
}

func TestUnit__RandomEncodeDecode(t *testing.T) {
	for i := 1; i < 2000; i += 5 {
		data := make([]string, 0, i)
		for j := 0; j < i; j++ {
			data = append(data, strconv.Itoa(utils.Rand.Int()))
		}
		original := make([]string, len(data))
		copy(original, data)

		T, err := tree.Decode(data)
		if err != nil {
			t.Fatalf("encoding failed with slice of length %d", len(data))
		}
		repr := T.Encode()
		if !reflect.DeepEqual(original, repr) {
			t.Fatalf("encode/decode test failed")
		}
	}
}

func TestPerf__RandomEncodeDecode(t *testing.T) {
	for i := 1; i < 2000; i++ {
		data := make([]string, 0, i)
		for j := 0; j < i; j++ {
			data = append(data, strconv.Itoa(utils.Rand.Int()))
		}
		T, err := tree.Decode(data)
		if err != nil {
			t.Fatalf("encoding failed with slice of length %d", len(data))
		}
		T.Encode()
	}
}
