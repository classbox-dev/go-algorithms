package strings

import (
	"errors"
)

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// LCS returns the longest common subsequence of the given strings.
func LCS(s1, s2 string) string {

	m, n := len(s1)+1, len(s2)+1
	M := make([]int, n*m)

	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			if s1[i-1] == s2[j-1] {
				M[n*i+j] = M[n*(i-1)+j-1] + 1
			} else {
				M[n*i+j] = max(M[n*i+j-1], M[n*(i-1)+j])
			}
		}
	}

	b := make([]byte, 0, max(n, m))

	for i, j := m-1, n-1; i > 0 || j > 0; {
		if i > 0 && j > 0 && s1[i-1] == s2[j-1] {
			b = append(b, s1[i-1])
			i--
			j--
		} else if i > 0 && (j == 0 || M[(i-1)*n+j] >= M[i*n+j-1]) {
			i--
		} else if j > 0 && (i == 0 || M[(i-1)*n+j] < M[i*n+j-1]) {
			j--
		}
	}

	N := len(b)
	for i := 0; i < N/2; i++ {
		b[i], b[N-1-i] = b[N-1-i], b[i]
	}
	return string(b)
}

// Segmentation breaks the given string into words.
// A 'word' is a non-empty string for which the given isWord function returns true.
// An error is returned if the segmentation is not possible.
func Segmentation(s string, isWord func(w string) bool) ([]string, error) {
	n := len(s)
	if n == 0 {
		return []string{}, nil
	}
	split := make([]int, n+1)
	split[n] = 1
	for r := n; r > 0; r-- {
		if split[r] > 0 {
			for l := r - 1; l >= 0; l-- {
				if isWord(s[l:r]) {
					split[l] = len(s[l:r])
				}
			}
		}
	}
	if split[0] == 0 {
		return nil, errors.New("no segmentation possible")
	}
	words := make([]string, 0, 4)
	for i := 0; i < n; {
		l := split[i]
		words = append(words, s[i:i+l])
		i += l
	}
	return words, nil
}
