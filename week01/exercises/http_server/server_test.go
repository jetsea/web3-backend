package http_server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSlowHandler_Completes(t *testing.T) {
	handler := SlowHandler(10 * time.Millisecond)
	req := httptest.NewRequest(http.MethodGet, "/slow", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
}

func TestSlowHandler_ClientDisconnects(t *testing.T) {
	handler := SlowHandler(500 * time.Millisecond)

	// Cancel the request context before the operation completes.
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	req := httptest.NewRequest(http.MethodGet, "/slow", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler(w, req)

	// Handler should respond with a non-200 status when context is cancelled.
	if w.Code == http.StatusOK {
		t.Error("expected non-200 status when client disconnects")
	}
}

func TestRunWithContext_NormalShutdown(t *testing.T) {
	// Create server with random port
	srv := NewServer(":0", 5*time.Second, 10*time.Second)

	// Create context that will be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() //reasure that cancel is called to avoid context leak in case of test failure

	// Run server in a goroutine
	errCh := make(chan error, 1)
	go func() {
		errCh <- RunWithContext(ctx, srv)
	}()

	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	// Cancel context to trigger shutdown
	cancel()

	// Wait for server to shutdown
	err := <-errCh
	t.Log(err)
	if err != nil && err != http.ErrServerClosed {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunWithContext_Error(t *testing.T) {
	// Create server with an invalid address
	srv := NewServer("invalid-address", 5*time.Second, 10*time.Second)

	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Run server
	err := RunWithContext(ctx, srv)
	t.Log(err)
	if err == nil {
		t.Error("expected error, got nil")
	}
}
