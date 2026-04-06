package context

import (
	"context"
	"errors"
	"time"
)

// SlowOperation simulates an operation that takes operationDuration.
// It returns an error if the context expires before the operation finishes.
func SlowOperation(ctx context.Context, operationDuration time.Duration) error {
	select {
	//if success, return nil. Simulate the operation takes long time.
	case <-time.After(operationDuration):
		return nil
	//if cancel, return error. if ctx.Done() is nil, this case will never be selected
	case <-ctx.Done():
		return ctx.Err()
	}
}

// TimeoutWrapper runs fn inside a context with the given timeout.
// Returns context.DeadlineExceeded if fn does not finish in time.
func TimeoutWrapper(timeout time.Duration, fn func(ctx context.Context) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return fn(ctx)
}

// IsTimeout reports whether err is a context deadline exceeded error.
func IsTimeout(err error) bool {
	return errors.Is(err, context.DeadlineExceeded)
}
