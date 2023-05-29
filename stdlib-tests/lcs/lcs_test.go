package lcs_test

import (
	xlcs "hsecode.com/stdlib-tests/internal/lcs"
	"hsecode.com/stdlib-tests/internal/rand"
	"hsecode.com/stdlib-tests/internal/utils"
	"hsecode.com/stdlib/strings"
	"runtime"
	"testing"
)

const charset = "abcdefghijklmnopqrstuvwxyz"

func TestUnit__Basic(t *testing.T) {
	r := rand.New(charset)
	for i := 2; i < 1300; i += 10 {
		s1, s2 := r.String(i/2+utils.Rand.Intn(i/2)), r.String(i/2+utils.Rand.Intn(i/2))
		result := strings.LCS(s1, s2)
		exLen := xlcs.Len(s1, s2)
		if exLen != len(result) {
			t.Fatalf(`Expected LCS("%v", "%v") of length %d, got %d`, s1, s2, exLen, len(result))
		}
		if !(isSubseq(result, s1) && isSubseq(result, s2)) {
			t.Fatalf(`"%s" is not a common subsequence of "%s" and "%s"`, result, s1, s2)
		}
		runtime.GC()
	}
}

func TestPerf__Basic(t *testing.T) {
	r := rand.New(charset)
	for i := 2; i < 1300; i += 5 {
		for j := 0; j < 2; j++ {
			s1, s2 := r.String(i/2+utils.Rand.Intn(i/2)), r.String(i/2+utils.Rand.Intn(i/2))
			result := strings.LCS(s1, s2)
			utils.Use(result)
			runtime.GC()
		}
	}
}

func isSubseq(t, s string) bool {
	tn, sn := len(t), len(s)
	c := 0
	for i := 0; i < sn && c < tn; i++ {
		if s[i] == t[c] {
			c++
		}
	}
	return c == tn
}
