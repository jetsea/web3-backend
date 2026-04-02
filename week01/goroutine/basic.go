// Package goroutine demonstrates basic goroutine concepts.
// A goroutine is a lightweight thread managed by the Go runtime.
// Unlike OS threads, goroutines cost only a few KB of stack memory.
// While OS threads cost a few MB of memory.
package goroutine

import (
	"fmt"
	"sync"
	"time"
)

// Greet prints a greeting message after an optional delay.
func Greet(name string, delay time.Duration) {
	if delay > 0 {
		time.Sleep(delay)
	}
	fmt.Printf("Hello, %s!\n", name)
}

// RunSequential calls Greet synchronously — each call blocks until done.
func RunSequential(names []string) {
	for _, name := range names {
		Greet(name, 0)
	}
}

// RunConcurrent launches each Greet in a separate goroutine.
// It uses a WaitGroup to wait for all of them to finish.
func RunConcurrent(names []string) {
	var wg sync.WaitGroup
	for _, name := range names {
		wg.Add(1)
		go func(n string) {
			defer wg.Done()
			Greet(n, 10*time.Millisecond)
		}(name)
	}
	wg.Wait()
}

// Count returns how many goroutines finished within the deadline.
// This demonstrates that goroutines are killed when the program exits
// unless we explicitly wait for them.
// func Count(n int, deadline time.Duration) int {
// 	done := make(chan struct{}, n)

// 	for i := 0; i < n; i++ {
// 		go func() {
// 			time.Sleep(5 * time.Millisecond)
// 			done <- struct{}{}
// 		}()
// 	}

// 	var finished int
// 	timeout := time.After(deadline)
// 	for {
// 		select {
// 		case <-done:
// 			finished++
// 			if finished == n {
// 				return finished
// 			}
// 		case <-timeout:
// 			return finished
// 		}
// 	}
// }
