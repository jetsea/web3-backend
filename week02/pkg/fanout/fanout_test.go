package fanout

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestFanOut_Basic(t *testing.T) {
	ctx := context.Background()
	input := 65

	results := FanOut(ctx, input, 3, func(ctx context.Context, n int) (string, error) {
		time.Sleep(10 * time.Millisecond)
		return "result-" + string(rune(n)), nil
	})

	if len(results) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(results))
	}

	for i, result := range results {

		if result.Error != nil {
			t.Errorf("Result %d: unexpected error %v", i, result.Error)
		}
		if result.Index != i || result.Data != "result-"+string(rune(input)) {
			t.Errorf("Result %d: expected index %d and data 'result-%c', got index %d and data '%s'", i, i, rune(input), result.Index, result.Data)
		}
	}
}

func TestFanOut_WithError(t *testing.T) {
	ctx := context.Background()
	input := "test"

	results := FanOut(ctx, input, 4, func(ctx context.Context, s string) (int, error) {
		if s == "test" {
			return 0, errors.New("test error")
		}
		return len(s), nil
	})

	errorCount := 0
	for _, result := range results {
		if result.Error != nil {
			errorCount++
		}
	}

	if errorCount != 4 {
		t.Errorf("Expected 4 errors, got %d", errorCount)
	}
}

func TestFanOut_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	input := 100

	cancel() // Cancel immediately

	results := FanOut(ctx, input, 2, func(ctx context.Context, n int) (string, error) {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
			return "done", nil
		}
	})

	// Results should still be returned, but might be incomplete
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	for i, result := range results {
		if result.Error != context.Canceled {
			t.Errorf("Result %d: expected context.Canceled error, got %v", i, result.Error)
		}
	}
}

func TestFanIn_Basic(t *testing.T) {
	ctx := context.Background()

	// Create channels with results
	ch1 := make(chan Result[int], 1)
	ch2 := make(chan Result[int], 1)
	ch3 := make(chan Result[int], 1)

	ch1 <- Result[int]{Index: 0, Data: 10, Error: nil}
	ch2 <- Result[int]{Index: 1, Data: 20, Error: nil}
	ch3 <- Result[int]{Index: 2, Data: 30, Error: nil}

	close(ch1)
	close(ch2)
	close(ch3)

	out := FanIn(ctx, ch1, ch2, ch3)

	var results []int
	for result := range out {
		results = append(results, result.Data)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}
}

func TestConcurrentQuery_Basic(t *testing.T) {
	ctx := context.Background()
	inputs := []int{1, 2, 3, 4, 5}

	results, err := ConcurrentQuery(ctx, inputs, func(ctx context.Context, n int) (int, error) {
		return n * 2, nil
	})

	if err != nil {
		t.Fatalf("ConcurrentQuery failed: %v", err)
	}

	if len(results) != 5 {
		t.Fatalf("Expected 5 results, got %d", len(results))
	}

	for i, result := range results {
		expected := inputs[i] * 2
		if result != expected {
			t.Errorf("Result %d: expected %d, got %d", i, expected, result)
		}
	}
}

func TestConcurrentQuery_WithError(t *testing.T) {
	ctx := context.Background()
	inputs := []int{1, 2, 3}

	results, err := ConcurrentQuery(ctx, inputs, func(ctx context.Context, n int) (int, error) {
		if n == 2 {
			return 0, errors.New("error on 2")
		}
		return n * 2, nil
	})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if results != nil {
		t.Error("Expected nil results on error")
	}
}

func TestConcurrentQuery_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	inputs := []int{1, 2, 3}

	go func() {
		time.Sleep(1 * time.Millisecond)
		cancel()
	}()

	_, err := ConcurrentQuery(ctx, inputs, func(ctx context.Context, n int) (int, error) {
		time.Sleep(2 * time.Millisecond)
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			return n * 2, nil
		}

	})

	// Should complete with or without error depending on timing
	if err != context.Canceled {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestAggregatePrices_Basic(t *testing.T) {
	ctx := context.Background()

	ch1 := make(chan float64, 1)
	ch2 := make(chan float64, 1)
	ch3 := make(chan float64, 1)

	ch1 <- 100.0
	ch2 <- 102.0
	ch3 <- 98.0

	close(ch1)
	close(ch2)
	close(ch3)

	out := AggregatePrices(ctx, ch1, ch2, ch3)
	agg := <-out

	if !agg.Valid {
		t.Fatal("Expected valid aggregation")
	}

	if agg.Average != 100.0 {
		t.Errorf("Expected average 100.0, got %.2f", agg.Average)
	}

	if agg.Min != 98.0 {
		t.Errorf("Expected min 98.0, got %.2f", agg.Min)
	}

	if agg.Max != 102.0 {
		t.Errorf("Expected max 102.0, got %.2f", agg.Max)
	}
}

func TestAggregatePrices_Empty(t *testing.T) {
	ctx := context.Background()

	out := AggregatePrices(ctx)
	agg := <-out

	if agg.Valid {
		t.Error("Expected invalid aggregation for empty input")
	}
}

func TestFanOut_Concurrency(t *testing.T) {
	ctx := context.Background()
	var count atomic.Int32

	results := FanOut(ctx, nil, 100, func(ctx context.Context, _ interface{}) (bool, error) {
		count.Add(1)
		time.Sleep(10 * time.Millisecond)
		return true, nil
	})

	// Give time for all goroutines to start
	time.Sleep(20 * time.Millisecond)

	if results == nil || len(results) != 100 {
		t.Errorf("Expected 100 results, got %d", len(results))
	}

	finalCount := count.Load()
	if finalCount != 100 {
		t.Errorf("Expected 100 completions, got %d", finalCount)
	}
}

func TestFanIn_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ch1 := make(chan Result[string], 1)
	ch2 := make(chan Result[string], 1)

	ch1 <- Result[string]{Index: 0, Data: "done"}

	go func() {
		time.Sleep(1 * time.Millisecond)
		cancel()
	}()

	out := FanIn(ctx, ch1, ch2)

	count := 0
	for range out {
		count++
	}

	// Should close without waiting for ch2
	if count == 0 {
		t.Error("Expected at least one result")
	}
}
