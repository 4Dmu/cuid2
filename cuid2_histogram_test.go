package cuid2

import (
	"math"
	"strings"
	"testing"
)

func TestCuid2Histogram(t *testing.T) {
	amount := 100000
	t.Logf("generating %d unique IDS", amount)
	pool := createIdPool(amount, t)
	sampleIds := pool.ids[0:10]
	t.Log(len(pool.ids))

	for id, count := range pool.idsMap {
		if count > 1 {
			t.Errorf("collision detected: value (%s) found %d times", id, count)
		}
	}

	{
		tolerance := 0.1
		idLength := 23
		totalLetters := idLength * amount
		base := 36
		exptectedBinSize := math.Ceil(float64(totalLetters) / float64(base))
		minBinSize := math.Round(exptectedBinSize * (1 - tolerance))
		maxBinSize := math.Round(exptectedBinSize * (1 + tolerance))
		charFrequencies := make(map[string]float64)

		for _, id := range pool.ids {
			for _, char := range strings.Split(id[2:], "") {
				charFrequencies[char] = charFrequencies[char] + 1
			}
		}

		t.Logf("testing character frequency...")
		t.Logf("expectedBinSize: %f", exptectedBinSize)
		t.Logf("minBinSize: %f", minBinSize)
		t.Logf("maxBinSize: %f", maxBinSize)
		t.Logf("charFrequencies: %v", charFrequencies)

		for char, frequency := range charFrequencies {
			if frequency > maxBinSize || frequency < minBinSize {
				t.Errorf("bin size range exceeded for char (%s) min (%f) max (%f) got (%f)", char, minBinSize, maxBinSize, frequency)
			}
		}

		if len(charFrequencies) != base {
			t.Errorf("length of charFrequencies should equal base, got (%d) expected (%d)", len(charFrequencies), base)
		}

	}

	{
		tolerance := 0.1
		exptectedBinSize := math.Ceil(float64(amount) / float64(len(pool.histogram)))
		minBinSize := math.Round(exptectedBinSize * (1 - tolerance))
		maxBinSize := math.Round(exptectedBinSize * (1 + tolerance))

		t.Logf("sample ids:")
		t.Logf("%v", sampleIds)
		t.Logf("expectedBinSize: %f", exptectedBinSize)
		t.Logf("minBinSize: %f", minBinSize)
		t.Logf("maxBinSize: %f", maxBinSize)

		for _, h := range pool.histogram {
			if h < int(minBinSize) || h > int(maxBinSize) {
				t.Errorf("bin size range exceeded for histogram min (%f) max (%f) got (%d)", minBinSize, maxBinSize, h)
			}
		}
	}

}
