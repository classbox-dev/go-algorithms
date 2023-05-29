package rand

import (
	"hsecode.com/stdlib-tests/v2/internal/utils"
)

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
		bytes[i] = r.alphabet[utils.Rand.Intn(n)]
	}
	return string(bytes)
}
