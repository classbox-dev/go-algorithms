package rand

import "math/rand"

type Rand struct {
	alphabet string
	length   int
}

func New(alphabet string) *Rand {
	r := new(Rand)
	r.alphabet = alphabet
	r.length = len(alphabet)
	return r
}

func (r *Rand) String(length int) string {
	bytes := make([]byte, length)
	n := r.length
	for i := 0; i < length; i++ {
		bytes[i] = r.alphabet[rand.Intn(n)]
	}
	return string(bytes)
}
