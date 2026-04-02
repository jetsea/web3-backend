package goroutine

import (
	"testing"
	"time"
)

func TestRunWorkers_Count(t *testing.T) {
	results := RunWorkers(5)
	if len(results) != 5 {
		t.Errorf("got %d results, want 5", len(results))
	}
}

func TestRunWorkers_Concurrency(t *testing.T) {
	start := time.Now()
	RunWorkers(5) // each worker sleeps 10ms
	elapsed := time.Since(start)

	// 5 concurrent workers each sleeping 10ms should finish well under 50ms.
	if elapsed > 40*time.Millisecond {
		t.Errorf("RunWorkers took %v; workers should run concurrently", elapsed)
	}
}
