package context

import (
	"context"
	"testing"
	"time"
)

func TestDoWork_Completes(t *testing.T) {
	ctx := context.Background()
	result := DoWork(ctx, 10*time.Millisecond)
	if result != "done" {
		t.Errorf("DoWork = %q, want %q", result, "done")
	}
}

func TestDoWork_Cancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	result := DoWork(ctx, time.Second)
	if result == "done" {
		t.Error("expected cancellation, got done")
	}
}

func TestCancelAfter_WorkFinishesFirst(t *testing.T) {
	// Work takes 10ms, cancel fires after 100ms → work should win.
	result := CancelAfter(10*time.Millisecond, 100*time.Millisecond)
	if result != "done" {
		t.Errorf("CancelAfter = %q, want %q", result, "done")
	}
}

func TestCancelAfter_CancelWins(t *testing.T) {
	// Work takes 200ms, cancel fires after 20ms → cancel should win.
	result := CancelAfter(200*time.Millisecond, 20*time.Millisecond)
	if result == "done" {
		t.Error("expected cancellation result, got done")
	}
}
