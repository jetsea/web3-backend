package context

import (
	"context"
	"testing"
	"time"
)

func TestCancelOrFinish_FinishesFirst(t *testing.T) {
	// Work takes 10ms, cancel fires after 100ms → work should win.
	result := CancelOrFinish(10*time.Millisecond, 100*time.Millisecond)
	if result != "done" {
		t.Errorf("CancelOrFinish = %q, want %q", result, "done")
	}
}

func TestCancelOrFinish_CancelWins(t *testing.T) {
	// Work takes 200ms, cancel fires after 20ms → cancel should win.
	result := CancelOrFinish(200*time.Millisecond, 20*time.Millisecond)
	if result == "done" {
		t.Error("expected cancellation result, got done")
	}
}

func TestPrintTree_ParentNotCancelled(t *testing.T) {
	// Test with a non-cancelled parent context
	ctx := context.Background()
	state, cleanup := PrintTree(ctx)
	defer cleanup()

	// the done state of context.Background() is always nil
	if !state.ParentDone {
		t.Error("expected parent to be nil")
	}
	// the done state of context.WithCancel() is always non-nil
	if state.Child1Done {
		t.Error("expected child1 to not be nil")
	}
	if state.Child2Done {
		t.Error("expected child2 to not be nil")
	}
	if state.GrandchildDone {
		t.Error("expected grandchild to not be nil")
	}
}

func TestPrintTree_ParentCancelled(t *testing.T) {
	// Test with a cancelled parent context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the parent
	state, _ := PrintTree(ctx)

	// the done state of context.WithCancel() is always non-nil
	if state.ParentDone {
		t.Error("expected parent to be done")
	}
	if state.Child1Done {
		t.Error("expected child1 to be done")
	}
	if state.Child2Done {
		t.Error("expected child2 to be done")
	}
	if state.GrandchildDone {
		t.Error("expected grandchild to be done")
	}
}
