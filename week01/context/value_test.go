package context

import (
	"context"
	"testing"
)

func TestWithUserID(t *testing.T) {
	ctx := WithUserID(context.Background(), "user-42")
	got, ok := UserIDFrom(ctx)
	if !ok || got != "user-42" {
		t.Errorf("UserIDFrom = (%q, %v), want (\"user-42\", true)", got, ok)
	}
}

func TestUserIDFrom_Missing(t *testing.T) {
	_, ok := UserIDFrom(context.Background())
	if ok {
		t.Error("expected ok=false for empty context")
	}
}

func TestWithRequestID(t *testing.T) {
	ctx := WithRequestID(context.Background(), "req-123")
	got, ok := RequestIDFrom(ctx)
	if !ok || got != "req-123" {
		t.Errorf("RequestIDFrom = (%q, %v), want (\"req-123\", true)", got, ok)
	}
}

func TestBuildRequestContext(t *testing.T) {
	ctx := BuildRequestContext(context.Background(), "user-1", "req-1", "trace-1")

	userID, ok := UserIDFrom(ctx)
	if !ok || userID != "user-1" {
		t.Errorf("userID = %q, want %q", userID, "user-1")
	}

	requestID, ok := RequestIDFrom(ctx)
	if !ok || requestID != "req-1" {
		t.Errorf("requestID = %q, want %q", requestID, "req-1")
	}
}

func TestLogFields(t *testing.T) {
	ctx := BuildRequestContext(context.Background(), "alice", "r-99", "t-1")
	fields := LogFields(ctx)
	if fields == "" {
		t.Error("LogFields should not be empty")
	}
}

func TestContextValueIsolation(t *testing.T) {
	// Values set in a child context must not be visible in the parent.
	parent := context.Background()
	child := WithUserID(parent, "child-user")

	if _, ok := UserIDFrom(parent); ok {
		t.Error("parent context should NOT see child's userID")
	}
	if id, ok := UserIDFrom(child); !ok || id != "child-user" {
		t.Error("child context should see its own userID")
	}
}
