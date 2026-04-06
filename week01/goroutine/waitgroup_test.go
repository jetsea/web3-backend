package goroutine

import (
	"sync"
	"testing"
	"time"
)

func TestNilWorker(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	Worker(1, &wg, nil)
}

func TestRunWorkers_Concurrency(t *testing.T) {
	start := time.Now()
	results := RunWorkers(5) // each worker sleeps 10ms
	elapsed := time.Since(start)

	// 5 concurrent workers each sleeping 10ms should finish well under 50ms.
	if elapsed > 40*time.Millisecond {
		t.Errorf("RunWorkers took %v; workers should run concurrently", elapsed)
	}
	if len(results) != 5 {
		t.Errorf("got %d results, want 5", len(results))
	}
}
