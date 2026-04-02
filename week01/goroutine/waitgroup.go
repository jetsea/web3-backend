package goroutine

import (
	"fmt"
	"sync"
	"time"
)

// Worker simulates a unit of work that takes some time.
func Worker(id int, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(10 * time.Millisecond) // simulate work
	result := fmt.Sprintf("Worker %d done", id)
	fmt.Println(result)

	if results != nil {
		results <- result
	}
}

// RunWorkers launches n workers concurrently and waits for all to finish.
// It returns the collected result strings.
func RunWorkers(n int) []string {
	var wg sync.WaitGroup
	results := make(chan string, n)

	for i := 1; i <= n; i++ {
		wg.Add(1)
		go Worker(i, &wg, results)
	}

	// Close results channel once all workers are done.
	// To avoid blocking, we use a goroutine.
	go func() {
		wg.Wait()
		close(results)
	}()

	var out []string
	for r := range results {
		out = append(out, r)
	}
	return out
}
