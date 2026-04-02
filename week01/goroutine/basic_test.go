package goroutine

import (
	"testing"
	"time"
)

func TestRunSequential(t *testing.T) {
	// Should complete without panic
	RunSequential([]string{"Alice", "Bob", "Charlie"})
}

func TestRunConcurrent(t *testing.T) {
	names := []string{"Alice", "Bob", "Charlie"}
	start := time.Now()
	RunConcurrent(names)
	elapsed := time.Since(start)

	// All goroutines run with 10ms delay concurrently,
	// so total should be much less than 4*10ms = 40ms.
	if elapsed > 20*time.Millisecond {
		t.Errorf("RunConcurrent took %v, expected < 20ms (goroutines should run concurrently)", elapsed)
	}
}

// func TestCount_AllFinish(t *testing.T) {
// 	got := Count(5, 200*time.Millisecond)
// 	if got != 5 {
// 		t.Errorf("Count = %d, want 5", got)
// 	}
// }

// func TestCount_Deadline(t *testing.T) {
// 	// Goroutines sleep 5ms, but deadline is 1ms — none should finish.
// 	got := Count(10, 1*time.Millisecond)
// 	if got != 0 {
// 		t.Errorf("Count = %d, want 0 (deadline too short)", got)
// 	}
// }
