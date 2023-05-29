package levenshtein_test

import (
	ref "hsecode.com/stdlib-tests/internal/levenshtein"
	"hsecode.com/stdlib-tests/internal/rand"
	"hsecode.com/stdlib-tests/internal/utils"
	"hsecode.com/stdlib/strings/levenshtein"
	"runtime"
	"testing"
)

const (
	Insert  = byte('I')
	Replace = byte('R')
	Delete  = byte('D')
	Match   = byte('M')
)

func transcriptDistance(t string) int {
	distance := 0
	for i := 0; i < len(t); i++ {
		if t[i] == Insert || t[i] == Replace || t[i] == Delete {
			distance++
		}
	}
	return distance
}

func applyTranscript(t, src, dst string) string {
	output := make([]byte, 0, len(dst))
	a, b := 0, 0
	for i := 0; i < len(t); i++ {
		switch t[i] {
		case Insert:
			output = append(output, dst[b])
			b++
		case Replace:
			output = append(output, dst[b])
			a++
			b++
		case Delete:
			a++
		case Match:
			output = append(output, src[a])
			a++
			b++
		}
	}
	return string(output)
}

func transcriptValid(t, src, dst string) (ok bool) {
	ok = false
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()
	output := applyTranscript(t, src, dst)
	ok = output == dst
	return
}

func TestUnit__Basic(t *testing.T) {
	inputs := []struct {
		s1, s2 string
	}{
		{"", "aaaaa"},
		{"bbbbb", ""},
		{"a", "b"},
		{"a", "a"},
		{"vintner", "writers"},
		{"923423235", "135312512352"},
	}

	for _, inp := range inputs {

		lsh := levenshtein.New(inp.s1, inp.s2)
		distance := lsh.Distance()
		transcript := lsh.Transcript()

		expected := ref.Distance(inp.s1, inp.s2)

		if distance != expected {
			t.Fatalf(`New("%v", "%v").Distance() == %v, expected %v`, inp.s1, inp.s2, distance, expected)
		}
		if transcriptDistance(transcript) != expected {
			t.Fatalf(`New("%v", "%v").Transcript() is not optimal`, inp.s1, inp.s2)
		}
		if !transcriptValid(transcript, inp.s1, inp.s2) {
			t.Fatalf(`New("%v", "%v").Transcript() is invalid`, inp.s1, inp.s2)
		}
	}
	runtime.GC()
}

func TestUnit__RandomSmall(t *testing.T) {
	r := rand.New("aaaabcdeeefggggh")
	for l := 1; l <= 20; l++ {
		for i := 0; i < 20; i++ {
			s1 := r.String(l)
			s2 := r.String(l)

			lsh := levenshtein.New(s1, s2)
			distance := lsh.Distance()
			transcript := lsh.Transcript()

			expected := ref.Distance(s1, s2)

			if distance != expected {
				t.Fatalf(`New("%v", "%v").Distance() == %v, expected %v`, s1, s2, distance, expected)
			}
			if transcriptDistance(transcript) != expected {
				t.Fatalf(`New("%v", "%v").Transcript() is not optimal`, s1, s2)
			}
			if !transcriptValid(transcript, s1, s2) {
				t.Fatalf(`New("%v", "%v").Transcript() is invalid`, s1, s2)
			}
		}
	}
	runtime.GC()
}

func TestUnit__RandomLarge(t *testing.T) {
	r := rand.New("aaaabcdeeefggggh")
	for l := 2; l <= 500; l++ {
		s1 := r.String(l)
		s2 := r.String(l)
		expected := ref.Distance(s1, s2)

		lsh := levenshtein.New(s1, s2)
		distance := lsh.Distance()
		transcript := lsh.Transcript()

		if distance != expected {
			t.Fatalf("Distance is not optimal")
		}
		if transcriptDistance(transcript) != expected {
			t.Fatalf("Transcript is not optimal")
		}
		if !transcriptValid(transcript, s1, s2) {
			t.Fatalf("Transcript is invalid")
		}
	}
	runtime.GC()
}

func TestPerf__Random(t *testing.T) {
	runtime.GC()
	r := rand.New("aaaabcdeeefggggh")
	for l := 500; l <= 800; l++ {
		s1 := r.String(l)
		s2 := r.String(l)
		v := levenshtein.New(s1, s2)
		utils.Use(v.Distance())
		utils.Use(v.Transcript())
		runtime.GC()
	}
}
