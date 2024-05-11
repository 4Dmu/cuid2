package cuid2

import (
	"fmt"
	"math/big"
	"math/rand"

	"github.com/martinlindhe/base36"
	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

type Counter = func() int

func base36Encode[K Number](num K) string {
	return base36.Encode(uint64(num))
}

func base36EncodeBigint(bigInt *big.Int) string {
	return base36.EncodeBytes([]byte(fmt.Sprintf("%d", bigInt)))
}

func createAlphabet() []rune {
	letters := []rune{}
	for i := 0; i < 26; i++ {
		letters = append(letters, rune(i+97))
	}
	return letters
}

func randomLetter(alphabet []rune) rune {
	return alphabet[rand.Intn(len(alphabet))]
}

func createCounter(count int) Counter {
	return func() int {
		return count + 1
	}
}

func keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for v := range m {
		r = append(r, v)
	}
	return r
}
