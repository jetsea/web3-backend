// Package context demonstrates Go's context package for cancellation,
// timeouts, and request-scoped value propagation.
package context

import (
	"context"
	"fmt"
	"time"
)

// DoWork simulates a long-running task that respects context cancellation.
// It returns "done" if work finishes, or "cancelled" if the context is cancelled first.
func DoWork(ctx context.Context, duration time.Duration) string {
	select {
	case <-time.After(duration):
		return "done"
	case <-ctx.Done():
		return "cancelled: " + ctx.Err().Error()
	}
}

// CancelAfter launches DoWork in a goroutine and cancels the context
// after cancelAfter duration.  Returns the result string.
func CancelAfter(workDuration, cancelAfter time.Duration) string {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(cancelAfter)
		cancel()
	}()

	return DoWork(ctx, workDuration)
}

// PrintTree demonstrates how a cancellation propagates down a context tree.
// Parent cancellation cancels all children.
func PrintTree(ctx context.Context) {
	child1, cancel1 := context.WithCancel(ctx)
	defer cancel1()

	child2, cancel2 := context.WithCancel(ctx)
	defer cancel2()

	grandchild, cancelGC := context.WithCancel(child1)
	defer cancelGC()

	fmt.Printf("parent done:     %v\n", ctx.Done() == nil)
	fmt.Printf("child1 done:     %v\n", child1.Done() == nil)
	fmt.Printf("child2 done:     %v\n", child2.Done() == nil)
	fmt.Printf("grandchild done: %v\n", grandchild.Done() == nil)
}
