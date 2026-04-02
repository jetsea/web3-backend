// Package context demonstrates Go's context package for cancellation,
// timeouts, and request-scoped value propagation.
package context

import (
	"fmt"
	"net/http"
	"time"
)

// HTTPServer demonstrates how to use context with HTTP servers
func HTTPServer() {
	// Create a handler that uses context
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get context from request
		ctx := r.Context()

		// Simulate a long-running operation
		select {
		case <-time.After(2 * time.Second):
			fmt.Fprint(w, "Request processed successfully")
		case <-ctx.Done():
			fmt.Println("Request cancelled:", ctx.Err())
			http.Error(w, "Request cancelled", http.StatusRequestTimeout)
		}
	})

	// Create server with timeout
	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Server starting on :8080")
	fmt.Println("To test, run: curl -m 1 http://localhost:8080")

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	// Wait a bit to show the server is running
	time.Sleep(1 * time.Second)
	fmt.Println("Server demonstration complete")
}
