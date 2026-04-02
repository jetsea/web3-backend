// Package crawler demonstrates a concurrent URL fetcher using goroutines,
// WaitGroup, and buffered channels.
package crawler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// Result holds the outcome of a single fetch.
type Result struct {
	URL      string
	Bytes    int
	Duration time.Duration
	Err      error
}

func (r Result) String() string {
	if r.Err != nil {
		return fmt.Sprintf("ERROR %s: %v", r.URL, r.Err)
	}
	return fmt.Sprintf("OK    %s  %d bytes  %v", r.URL, r.Bytes, r.Duration)
}

// fetch performs a single HTTP GET and sends the Result to out.
func fetch(ctx context.Context, url string, out chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		out <- Result{URL: url, Err: err}
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		out <- Result{URL: url, Duration: time.Since(start), Err: err}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	out <- Result{
		URL:      url,
		Bytes:    len(body),
		Duration: time.Since(start),
		Err:      err,
	}
}

// FetchAll concurrently fetches all URLs and returns the results.
// The overall operation is bounded by ctx.
func FetchAll(ctx context.Context, urls []string) []Result {
	results := make(chan Result, len(urls))
	var wg sync.WaitGroup

	for _, u := range urls {
		wg.Add(1)
		go fetch(ctx, u, results, &wg)
	}

	// Close results once all fetches are done.
	go func() {
		wg.Wait()
		close(results)
	}()

	out := make([]Result, 0, len(urls))
	for r := range results {
		out = append(out, r)
	}
	return out
}
