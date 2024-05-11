package cuid2

import (
	"maps"
	"math"
	"slices"
	"sync"
	"testing"
)

type parrallelPools struct {
	mu    sync.Mutex
	pools []idPool
}

func (p *parrallelPools) addPool(pool idPool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.pools = append(p.pools, pool)
}

func createIdPoolsInParralel(numRountines, max int, t *testing.T) []idPool {
	wg := sync.WaitGroup{}
	pools := parrallelPools{pools: make([]idPool, 0)}

	for i := 0; i < numRountines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pool := createIdPool(max, t)
			pools.addPool(pool)
		}()
	}
	wg.Wait()
	return pools.pools
}

func TestCuid2Collisions(t *testing.T) {
	amount := int(math.Pow(7, 8) * 2)
	numPools := 7
	pools := createIdPoolsInParralel(numPools, amount/numPools, t)
	ids := make([]string, 0)
	set := make(map[string]int)

	for _, pool := range pools {
		ids = slices.Concat(pool.ids)
		maps.Copy(set, pool.idsMap)
	}

	sampleIds := ids[:10]
	t.Logf("sample ids: %v", sampleIds)
	histogram := pools[0].histogram
	expectedBinSize := math.Ceil(float64(amount) / float64(numPools) / float64(len(histogram)))
	tolerance := 0.05
	minBinSize := math.Round(expectedBinSize * (1 - tolerance))
	maxBinSize := math.Round(expectedBinSize * (1 + tolerance))

	if len(set) != amount {
		t.Fatalf("expected: %d, got:%d", amount, len(set))
	}

	failed := false
	for key, count := range set {
		if count > 1 {
			failed = true
			t.Logf("(%s) is duplicated %d times", key, count)
		}
	}

	if failed {
		t.Fatal("duplicated ids found")
	}

	for _, h := range histogram {
		if h < int(minBinSize) || h > int(maxBinSize) {
			t.Fatalf("Histogram is outside tollerance min (%f) max (%f) got (%d)", minBinSize, maxBinSize, h)
		}
	}

	t.Logf("%d ids generated", len(set))

	t.Logf("no duplicates found")
}
