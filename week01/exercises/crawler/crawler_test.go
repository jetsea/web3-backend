package crawler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// newTestServer returns a test HTTP server that sleeps for delay then responds.
func newTestServer(t *testing.T, delay time.Duration) *httptest.Server {
	t.Helper() //purpose of this line is to make error message  outside of the function
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if delay > 0 {
			time.Sleep(delay)
		}
		w.Write([]byte("hello"))
	}))
}

func TestFetchAll_Success(t *testing.T) {
	srv1 := newTestServer(t, 0)
	defer srv1.Close()
	srv2 := newTestServer(t, 0)
	defer srv2.Close()
	t.Logf("srv1.URL: %s, srv2.URL: %s", srv1.URL, srv2.URL)
	results := FetchAll(context.Background(), []string{srv1.URL, srv2.URL})

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	for _, r := range results {
		if r.Err != nil {
			t.Errorf("unexpected error for %s: %v", r.URL, r.Err)
		}
		if r.Bytes != 5 { // "hello" = 5 bytes
			t.Errorf("expected 5 bytes, got %d", r.Bytes)
		}
	}
}

func TestFetchAll_ContextCancellation(t *testing.T) {
	srv := newTestServer(t, 200*time.Millisecond)
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	results := FetchAll(ctx, []string{srv.URL})
	t.Logf("results: %v", results)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Err == nil {
		t.Error("expected timeout error, got nil")
	}
}

func TestFetchAll_InvalidURL(t *testing.T) {
	results := FetchAll(context.Background(), []string{"not-a-url"})
	t.Logf("results: %v", results)
	if len(results) != 1 {
		t.Fatalf("expected 1 result")
	}
	if results[0].Err == nil {
		t.Error("expected error for invalid URL")
	}
}

func TestResult_String(t *testing.T) {
	r := Result{URL: "http://example.com", Bytes: 100, Duration: 10 * time.Millisecond}
	s := r.String()
	t.Logf("result: %s", s)
	if s == "" {
		t.Error("String() should not be empty")
	}
}
