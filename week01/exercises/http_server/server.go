// Package http_server demonstrates an HTTP server that respects
// context cancellation so slow handlers don't hold open connections.
package http_server

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// SlowHandler simulates work that takes operationDuration.
// If the client disconnects (ctx.Done()), it returns early.
func SlowHandler(operationDuration time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		select {
		case <-time.After(operationDuration):
			fmt.Fprintln(w, "operation completed")
		case <-ctx.Done():
			http.Error(w, "request cancelled: "+ctx.Err().Error(), http.StatusServiceUnavailable)
		}
	}
}

// HealthHandler always responds 200 OK immediately.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
}

// EchoHandler echoes the query parameter "msg".
func EchoHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get("msg")
	fmt.Fprintln(w, msg)
}

// NewServer builds and returns a configured *http.Server.
func NewServer(addr string, readTimeout, writeTimeout time.Duration) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", HealthHandler)
	mux.HandleFunc("/echo", EchoHandler)
	mux.HandleFunc("/slow", SlowHandler(5*time.Second))

	return &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}

// RunWithContext starts srv and shuts it down gracefully when ctx is cancelled.
func RunWithContext(ctx context.Context, srv *http.Server) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}
