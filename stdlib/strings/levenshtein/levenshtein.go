package levenshtein

// Levenshtein is a data structure to compute Levenshtein distance and transcript between the given strings.
type Levenshtein struct {
	values     []int
	rows, cols int
	src, dst   string
}

// New creates a Levenshtein instance to transform src to dst.
// The function is not Unicode-aware i.e. it treats the arguments as sequences of bytes.
func New(src, dst string) *Levenshtein {
	rows, cols := len(src)+1, len(dst)+1
	values := make([]int, rows*cols)
	_ = values[rows*cols-1]

	for j := 0; j < cols; j++ {
		values[j] = j
	}
	for i := 1; i < rows; i++ {
		values[i*cols] = i
	}
	for i := 1; i < rows; i++ {
		for j := 1; j < cols; j++ {
			eqq := 0
			if src[i-1] != dst[j-1] {
				eqq = 1
			}
			v := min(values[(i-1)*cols+j]+1, min(values[i*cols+j-1]+1, values[(i-1)*cols+j-1]+eqq))
			values[i*cols+j] = v
		}
	}
	return &Levenshtein{
		values: values,
		rows:   rows,
		cols:   cols,
		src:    src,
		dst:    dst,
	}
}

// Distance returns the Levenshtein distance.
func (ls *Levenshtein) Distance() int {
	return ls.values[len(ls.values)-1]
}

// Transcript returns one of the edit transcripts corresponding to the Levenshtein distance.
func (ls *Levenshtein) Transcript() string {
	const (
		Insert  = byte('I')
		Replace = byte('R')
		Delete  = byte('D')
		Match   = byte('M')
	)
	M := ls.values
	_ = M[ls.rows*ls.cols-1]

	transcript := make([]byte, 0, len(ls.src)+len(ls.dst))

	i, j := ls.rows-1, ls.cols-1
	for i > 0 || j > 0 {
		v := M[i*ls.cols+j]
		eqq := 0
		if i > 0 && j > 0 && ls.src[i-1] != ls.dst[j-1] {
			eqq = 1
		}
		switch {
		case j > 0 && v == M[i*ls.cols+j-1]+1:
			transcript = append(transcript, Insert)
			j--
		case i > 0 && v == M[(i-1)*ls.cols+j]+1:
			transcript = append(transcript, Delete)
			i--
		case i > 0 && j > 0 && v == M[(i-1)*ls.cols+j-1]+eqq:
			if eqq > 0 {
				transcript = append(transcript, Replace)
			} else {
				transcript = append(transcript, Match)
			}
			i--
			j--
		}
	}
	for i, j := 0, len(transcript)-1; i < j; i, j = i+1, j-1 {
		transcript[i], transcript[j] = transcript[j], transcript[i]
	}
	return string(transcript)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
