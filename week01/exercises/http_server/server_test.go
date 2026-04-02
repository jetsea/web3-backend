package http_server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	HealthHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if w.Body.String() != "ok\n" {
		t.Errorf("body = %q, want %q", w.Body.String(), "ok\n")
	}
}

func TestEchoHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/echo?msg=hello", nil)
	w := httptest.NewRecorder()

	EchoHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if w.Body.String() != "hello\n" {
		t.Errorf("body = %q, want %q", w.Body.String(), "hello\n")
	}
}

func TestEchoHandler_EmptyMsg(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/echo", nil)
	w := httptest.NewRecorder()
	EchoHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
}

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

func TestNewServer_Config(t *testing.T) {
	srv := NewServer(":0", 5*time.Second, 10*time.Second)
	if srv == nil {
		t.Fatal("expected non-nil server")
	}
	if srv.ReadTimeout != 5*time.Second {
		t.Errorf("ReadTimeout = %v, want 5s", srv.ReadTimeout)
	}
	if srv.WriteTimeout != 10*time.Second {
		t.Errorf("WriteTimeout = %v, want 10s", srv.WriteTimeout)
	}
}

func TestNewServer_Routes(t *testing.T) {
	srv := NewServer(":0", 5*time.Second, 10*time.Second)
	ts := httptest.NewServer(srv.Handler)
	defer ts.Close()

	resp, err := ts.Client().Get(ts.URL + "/health")
	if err != nil {
		t.Fatalf("GET /health: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
}
