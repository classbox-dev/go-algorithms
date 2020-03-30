package lcs

func Len(s1, s2 string) int {
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
	return M[n*m-1]
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
