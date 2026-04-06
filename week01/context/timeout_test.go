package context

import (
	"context"
	"testing"
	"time"
)

func TestSlowOperation_Completes(t *testing.T) {
	ctx := context.Background()
	err := SlowOperation(ctx, 10*time.Millisecond)
	if err != nil {
		t.Errorf("SlowOperation returned error %v, want nil", err)
	}
}

func TestSlowOperation_Timeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := SlowOperation(ctx, 200*time.Millisecond)
	if err == nil {
		t.Error("expected timeout error, got nil")
	}
	if !IsTimeout(err) {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}

func TestTimeoutWrapper_Completes(t *testing.T) {
	err := TimeoutWrapper(20*time.Millisecond, func(ctx context.Context) error {
		return SlowOperation(ctx, 10*time.Millisecond)
	})
	if err != nil {
		t.Errorf("TimeoutWrapper returned %v, want nil", err)
	}
}

func TestTimeoutWrapper_Exceeded(t *testing.T) {
	err := TimeoutWrapper(10*time.Millisecond, func(ctx context.Context) error {
		return SlowOperation(ctx, 20*time.Millisecond)
	})
	if !IsTimeout(err) {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}

func TestIsTimeout(t *testing.T) {
	if !IsTimeout(context.DeadlineExceeded) {
		t.Error("IsTimeout(DeadlineExceeded) should be true")
	}
	if IsTimeout(context.Canceled) {
		t.Error("IsTimeout(Canceled) should be false")
	}
}
