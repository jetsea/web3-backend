package context

import (
	"context"
	"errors"
	"net/http"
	"time"
)

// FetchWithTimeout makes an HTTP GET to url with the given timeout.
// Returns the HTTP status code, or an error if the request fails or times out.
func FetchWithTimeout(url string, timeout time.Duration) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

// SlowOperation simulates an operation that takes operationDuration.
// It returns an error if the context expires before the operation finishes.
func SlowOperation(ctx context.Context, operationDuration time.Duration) error {
	select {
	case <-time.After(operationDuration):
		return nil
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
