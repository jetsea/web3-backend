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

// NewServer builds and returns a configured *http.Server.
func NewServer(addr string, readTimeout, writeTimeout time.Duration) *http.Server {
	mux := http.NewServeMux()
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
		//when srv shuts down, or port is already in use, or other error occurs, ListenAndServe will return a non-nil error, which we send to errCh
		errCh <- srv.ListenAndServe()
	}()

	select {
	//1.port is already in use 2. invalid addr 3. permission denied 4. server crashes 5. memory leak 6. srv.Close() is called 7. srv.Shutdown() is called outside.
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}
