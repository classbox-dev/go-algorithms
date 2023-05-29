package strings_test

import (
	"fmt"
	"hsecode.com/stdlib/v2/strings"
)

func ExampleLCS() {
	fmt.Println(strings.LCS("vintner", "writers"))
	fmt.Println(strings.LCS("ABCD", "ACBAD"))
	// Output:
	// iter
	// ABD
}
