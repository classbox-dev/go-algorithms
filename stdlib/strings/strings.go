package strings

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

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
