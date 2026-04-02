package channel

import (
	"context"
	"testing"
	"time"
)

func TestRace_FirstWins(t *testing.T) {
	a := make(chan int, 1)
	b := make(chan int, 1)

	a <- 10
	b <- 20

	got := Race(a, b)
	// Either channel may win; just ensure the value is one of the two.
	if got != 10 && got != 20 {
		t.Errorf("Race returned %d, want 10 or 20", got)
	}
}

func TestRace_SlowerLoses(t *testing.T) {
	a := make(chan int, 1)
	b := make(chan int, 1)

	// Only send on b; a is empty.
	b <- 99
	got := Race(a, b)
	if got != 99 {
		t.Errorf("Race = %d, want 99", got)
	}
}

func TestNonBlocking_ValuePresent(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 7

	v, ok := NonBlocking(ch)
	if !ok || v != 7 {
		t.Errorf("NonBlocking = (%d, %v), want (7, true)", v, ok)
	}
}

func TestNonBlocking_Empty(t *testing.T) {
	ch := make(chan int, 1)

	v, ok := NonBlocking(ch)
	if ok || v != 0 {
		t.Errorf("NonBlocking = (%d, %v), want (0, false)", v, ok)
	}
}

func TestWithDeadline_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	work := make(chan int, 5)
	work <- 1
	work <- 2
	cancel() // cancel immediately

	results := WithDeadline(ctx, work, time.Second)
	// May have captured 0-2 items; important is it returns without hanging.
	_ = results
}

func TestWithDeadline_Timeout(t *testing.T) {
	ctx := context.Background()
	work := make(chan int) // nothing will be sent

	start := time.Now()
	results := WithDeadline(ctx, work, 20*time.Millisecond)
	elapsed := time.Since(start)

	if len(results) != 0 {
		t.Errorf("expected empty results, got %v", results)
	}
	if elapsed < 20*time.Millisecond {
		t.Error("should have waited for the per-item timeout")
	}
}
