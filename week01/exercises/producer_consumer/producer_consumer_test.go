package producer_consumer

import (
	"testing"
	"time"
)

func TestRun_TotalResults(t *testing.T) {
	// 2 producers × 5 tasks = 10 total results
	results := Run(2, 5, 3, 10, 0)
	if len(results) != 10 {
		t.Errorf("expected 10 results, got %d", len(results))
	}
}

func TestRun_SingleProducerConsumer(t *testing.T) {
	results := Run(1, 3, 1, 5, 0)
	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}
}

func TestRun_ManyProducersFewConsumers(t *testing.T) {
	results := Run(5, 4, 1, 20, 0)
	if len(results) != 20 {
		t.Errorf("expected 20 results, got %d", len(results))
	}
}

func TestRun_Concurrency(t *testing.T) {
	// 3 consumers each with 10ms delay, 30 tasks total.
	// Serial: 30 * 10ms = 300ms.  Concurrent (3 workers): ~100ms.
	start := time.Now()
	results := Run(1, 30, 3, 30, 10*time.Millisecond)
	elapsed := time.Since(start)

	if len(results) != 30 {
		t.Fatalf("expected 30 results, got %d", len(results))
	}
	if elapsed > 200*time.Millisecond {
		t.Errorf("Run took %v; 3 concurrent consumers should finish in ~100ms", elapsed)
	}
}

func TestRun_NoTasks(t *testing.T) {
	results := Run(1, 0, 1, 1, 0)
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
