package cuid2

import (
	"math"
	"math/big"
	"strconv"
	"strings"
	"testing"
)

func parseInt(char string, radix int) int64 {
	val, _ := strconv.ParseInt(char, radix, 8)
	return val
}

func idToBigFloat(id string) *big.Float {
	chars := strings.Split(id, "")
	sum := big.NewFloat(0)
	radix := 36

	for _, char := range chars {
		bigRadix := big.NewFloat(float64(radix))
		bigChar := big.NewFloat(float64(parseInt(char, radix)))
		sumTimesRadix := sum.Mul(sum, bigRadix)
		sum = sumTimesRadix.Add(sumTimesRadix, bigChar)
	}

	return sum
}

func buildHistogram(numbers []*big.Float, bucketCount int) []int {
	buckets := make([]int, bucketCount)
	counter := 1
	bucketLength := math.Ceil(math.Pow(36, 23) / float64(bucketCount))

	for _, number := range numbers {

		devided := number.Quo(number, big.NewFloat(bucketLength))
		float, _ := devided.Float64()
		bucket := int64(math.Floor(float))
		buckets[bucket] = buckets[bucket] + 1
		counter++
	}
	return buckets
}

type idPool struct {
	ids       []string
	idsMap    map[string]int
	numbers   []*big.Float
	histogram []int
}

func createIdPool(max int, testing *testing.T) idPool {
	set := make(map[string]int)
	gen := New(GenOpts{})

	for i := 0; i < max; i++ {
		id := gen.Cuid2()
		set[id] = set[id] + 1

		// if i%10000 == 0 {
		// 	testing.Logf("%f", math.Floor((float64(i)/float64(max))*100))
		// }

		if set[id] > 1 {
			testing.Logf("collision at: %d", i)
			break
		}
	}

	testing.Log("no collisions detected")

	ids := keys(set)
	numbers := []*big.Float{}

	for _, id := range ids {
		numbers = append(numbers, idToBigFloat(id[1:]))
	}

	histogram := buildHistogram(numbers, 20)

	return idPool{
		ids:       ids,
		idsMap:    set,
		numbers:   numbers,
		histogram: histogram,
	}
}
