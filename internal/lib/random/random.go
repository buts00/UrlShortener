package random

import (
	"math/rand"
	"time"
)

func NewRandomAlias(aliasLength int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := make([]rune, 0)
	res := make([]rune, aliasLength)
	for r := 'a'; r <= 'z'; r++ {
		chars = append(chars, r)
	}
	for r := 'A'; r <= 'Z'; r++ {
		chars = append(chars, r)
	}
	for r := '0'; r <= '9'; r++ {
		chars = append(chars, r)
	}

	for i := range res {
		res[i] = chars[rnd.Intn(len(chars))]
	}
	return string(res)

}
