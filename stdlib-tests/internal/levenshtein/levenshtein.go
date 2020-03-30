package levenshtein

func Distance(s1, s2 string) int {
	s1len := len(s1)
	s2len := len(s2)
	column := make([]int, len(s1)+1)

	for y := 1; y <= s1len; y++ {
		column[y] = y
	}
	for x := 1; x <= s2len; x++ {
		column[0] = x
		lastdiag := x - 1
		for y := 1; y <= s1len; y++ {
			olddiag := column[y]
			var incr int
			if s1[y-1] != s2[x-1] {
				incr = 1
			}
			column[y] = min3(column[y]+1, column[y-1]+1, lastdiag+incr)
			lastdiag = olddiag
		}
	}
	return column[s1len]
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
	} else {
		if b < c {
			return b
		}
	}
	return c
}
