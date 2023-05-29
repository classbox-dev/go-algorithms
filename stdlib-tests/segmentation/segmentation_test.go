package segmentation_test

import (
	"fmt"
	"hsecode.com/stdlib-tests/v2/internal/rand"
	"hsecode.com/stdlib-tests/v2/internal/utils"
	xstrings "hsecode.com/stdlib/v2/strings"
	"runtime"
	"strings"
	"testing"
)

const charset = "abcdefghijklmnopqrstuvwxyABCDEFGHIJKLMNOPQRSTUVWXY0123456789"

type sentinel struct{}

func TestUnit__Empty(t *testing.T) {
	words, err := xstrings.Segmentation("", func(w string) bool {
		return true
	})
	if err != nil {
		t.Fatal("Segmentation() returned an error for a splittable string")
	}
	if len(words) != 0 {
		t.Fatal("Segmentation() returned non-empty slice for empty string")
	}
}

func TestUnit__Random(t *testing.T) {
	allWords := newDict(10000, 6)
	wordsSet := make(map[string]sentinel, len(allWords))
	for _, w := range allWords {
		wordsSet[w] = sentinel{}
	}
	isWord := func(w string) bool {
		_, ok := wordsSet[w]
		return ok
	}
	for size := 1; size < 20; size++ {
		for i := 0; i < 1000; i++ {
			valid := newString(allWords, size)
			if err := checkSingle(valid, false, isWord); err != nil {
				t.Fatal(err.Error())
			}
		}
		runtime.GC()
		for i := 0; i < 20; i++ {
			valid := newString(allWords, size)
			invalid := []string{valid + "z", "z" + valid, valid + "z" + valid}
			for _, s := range invalid {
				if err := checkSingle(s, true, isWord); err != nil {
					t.Fatal(err.Error())
				}
			}
		}
		runtime.GC()
	}
}

func TestUnit__RandomLong(t *testing.T) {
	allWords := newDict(30, 10000)
	wordsSet := make(map[string]sentinel, len(allWords))
	for _, w := range allWords {
		wordsSet[w] = sentinel{}
	}
	isWord := func(w string) bool {
		_, ok := wordsSet[w]
		return ok
	}
	for size := 1; size < 3; size++ {
		for i := 0; i < 10; i++ {
			valid := newString(allWords, size)
			if err := checkSingle(valid, false, isWord); err != nil {
				t.Fatal(err.Error())
			}
		}
		runtime.GC()
		for i := 0; i < 3; i++ {
			valid := newString(allWords, size)
			invalid := []string{valid + "z", "z" + valid, valid + "z" + valid}
			for _, s := range invalid {
				if err := checkSingle(s, true, isWord); err != nil {
					t.Fatal(err.Error())
				}
			}
		}
		runtime.GC()
	}
}

func TestPerf__Random(t *testing.T) {
	allWords := newDict(10000, 6)
	wordsSet := make(map[string]sentinel, len(allWords))
	for _, w := range allWords {
		wordsSet[w] = sentinel{}
	}
	isWord := func(w string) bool {
		_, ok := wordsSet[w]
		return ok
	}
	for size := 1; size < 20; size++ {
		for i := 0; i < 1000; i++ {
			s := newString(allWords, size)
			_, _ = xstrings.Segmentation(s, isWord)
			_, _ = xstrings.Segmentation("z"+s, isWord)
		}
		runtime.GC()
	}
}

func checkSingle(s string, expectError bool, isWord func(w string) bool) error {
	words, err := xstrings.Segmentation(s, isWord)
	if err != nil {
		if expectError {
			return nil
		}
		return fmt.Errorf("Segmentation() returned unexpected error")
	}
	if strings.Join(words, "") != s {
		return fmt.Errorf("Segmentation() returned slice that cannot be joined into the original string")
	}
	for _, word := range words {
		if !isWord(word) {
			return fmt.Errorf("Segmentation() returned an invalid word")
		}
	}
	return nil
}

func newDict(size int, maxLength int) []string {
	r := rand.New(charset)
	dict := make([]string, 0, size)
	for i := 0; i < size; i++ {
		s := r.String(1 + utils.Rand.Intn(maxLength))
		dict = append(dict, s)
	}
	return dict
}

func newString(dict []string, size int) string {
	var sb strings.Builder
	n := len(dict)
	for i := 0; i < size; i++ {
		sb.WriteString(dict[utils.Rand.Intn(n)])
	}
	return sb.String()
}
