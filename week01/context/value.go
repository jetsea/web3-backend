package context

import (
	"context"
	"fmt"
)

// contextKey is an unexported type for context keys in this package.
// Using a named type prevents key collisions with other packages.
type contextKey string

const (
	KeyUserID    contextKey = "userID"
	KeyRequestID contextKey = "requestID"
	KeyTraceID   contextKey = "traceID"
)

// WithUserID returns a new context carrying the given user ID.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, KeyUserID, userID)
}

// UserIDFrom extracts the user ID from the context.
// Returns ("", false) if not present.
func UserIDFrom(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(KeyUserID).(string)
	return v, ok
}

// WithRequestID attaches a request ID to the context.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, KeyRequestID, requestID)
}

// RequestIDFrom extracts the request ID from the context.
func RequestIDFrom(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(KeyRequestID).(string)
	return v, ok
}

// LogFields returns a formatted string of all known context values.
// Useful for structured logging in middleware.
func LogFields(ctx context.Context) string {
	userID, _ := UserIDFrom(ctx)
	requestID, _ := RequestIDFrom(ctx)
	return fmt.Sprintf("userID=%s requestID=%s", userID, requestID)
}

// BuildRequestContext builds a context with all standard request metadata.
func BuildRequestContext(ctx context.Context, userID, requestID, traceID string) context.Context {
	ctx = WithUserID(ctx, userID)
	ctx = WithRequestID(ctx, requestID)
	ctx = context.WithValue(ctx, KeyTraceID, traceID)
	return ctx
}
