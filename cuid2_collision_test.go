package cuid2

import (
	"fmt"
	"sync"
	"testing"
)

type CollisionTester struct {
	mu  sync.Mutex
	ids map[string]int
}

func (ct *CollisionTester) add(id string) {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	ct.ids[id] = ct.ids[id] + 1
}

func test(wg *sync.WaitGroup, tester *CollisionTester, gen *Gen) {
	defer wg.Done()
	for i := 0; i < 500; i++ {
		tester.add(gen.Cuid2())
	}
}

func TestCuid2Collisions(t *testing.T) {
	var wg sync.WaitGroup
	tester := CollisionTester{ids: make(map[string]int)}
	gen := New(GenOpts{})
	for i := 0; i < 500; i++ {
		wg.Add(1)
		go test(&wg, &tester, &gen)
	}
	wg.Wait()

	failed := false
	for key, count := range tester.ids {
		if count > 1 {
			failed = true
			fmt.Printf("(%s) is duplicated %d times\n", key, count)
		}
	}

	fmt.Printf("%d ids generated\n", len(tester.ids))

	if failed {
		t.Fatal("duplicated ids found")
	}

	fmt.Println("no duplicates found")

}
