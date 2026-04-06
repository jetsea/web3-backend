package channel

import (
	"context"
	"time"
)

// Race returns whichever of the two channels delivers a value first.
// It demonstrates select with two ready cases.
func Race(a, b <-chan int) int {
	select {
	case v := <-a:
		return v
	case v := <-b:
		return v
	}
}

// NonBlocking attempts to receive from ch without blocking.
// Returns (value, true) if a value is available, (0, false) otherwise.
func NonBlocking(ch <-chan int) (int, bool) {
	select {
	case v := <-ch:
		return v, true
	default:
		return 0, false
	}
}

// WithDeadline listens on work and returns values until ctx is done.
// It also applies a per-receive timeout via time.After.
func WithDeadline(ctx context.Context, work <-chan int, perItemTimeout time.Duration) []int {
	var results []int
	for {
		select {
		case <-ctx.Done():
			return results
		case v, ok := <-work:
			if !ok {
				return results
			}
			results = append(results, v)
		case <-time.After(perItemTimeout):
			// No item arrived within the timeout window.
			return results
		}
	}
}
