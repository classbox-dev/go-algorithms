package levenshtein_test

import (
	"fmt"
	"hsecode.com/stdlib/v2/strings/levenshtein"
)

func ExampleLevenshtein_Distance() {
	L := levenshtein.New("vintner", "writers")
	fmt.Println(L.Distance())
	// Output: 5
}

func ExampleLevenshtein_Transcript() {
	L := levenshtein.New("vintner", "writers")
	fmt.Println(L.Transcript())
	// Output: RIMDMDMMI
}
