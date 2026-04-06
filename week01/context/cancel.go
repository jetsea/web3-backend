// Package context demonstrates Go's context package for cancellation,
// timeouts, and request-scoped value propagation.
package context

import (
	"context"
	"fmt"
	"time"
)

// cancel or finish a context after a specified duration, and see which one wins.
func CancelOrFinish(workDuration, cancelAfter time.Duration) string {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(cancelAfter)
		cancel()
	}()

	select {
	case <-time.After(workDuration):
		return "done"
	case <-ctx.Done():
		return "cancelled: " + ctx.Err().Error()
	}
}

// ContextTreeState represents the state of a context tree.
type ContextTreeState struct {
	ParentDone     bool
	Child1Done     bool
	Child2Done     bool
	GrandchildDone bool
}

// PrintTree demonstrates how a cancellation propagates down a context tree.
// Parent cancellation cancels all children.
func PrintTree(ctx context.Context) (ContextTreeState, func()) {
	child1, cancel1 := context.WithCancel(ctx)

	child2, cancel2 := context.WithCancel(ctx)

	grandchild, cancelGC := context.WithCancel(child1)

	state := ContextTreeState{
		//context.Background() will never be cancelled, so ctx.Done() will always return nil
		ParentDone: ctx.Done() == nil,
		//context.WithCancel() will always return non-nil for Done()
		Child1Done:     child1.Done() == nil,
		Child2Done:     child2.Done() == nil,
		GrandchildDone: grandchild.Done() == nil,
	}

	// Print for debugging
	fmt.Printf("parent done:     %v\n", state.ParentDone)
	fmt.Printf("child1 done:     %v\n", state.Child1Done)
	fmt.Printf("child2 done:     %v\n", state.Child2Done)
	fmt.Printf("grandchild done: %v\n", state.GrandchildDone)
	// Return a cleanup function to cancel all contexts
	cleanup := func() {
		cancel1()
		cancel2()
		cancelGC()
	}

	return state, cleanup
}
