package goroutine

import "sync"

// UnsafeCounter is NOT goroutine-safe.
// Concurrent Inc() calls will produce a data race.
type UnsafeCounter struct {
	count int
}

// SafeCounter is goroutine-safe via sync.Mutex.
type SafeCounter struct {
	mu    sync.Mutex
	count int
}

type Counter interface {
	Inc()
	Value() int
}

// Inc increments the counter without any synchronisation.
func (c *UnsafeCounter) Inc() { c.count++ }

// Value returns the current count.
func (c *UnsafeCounter) Value() int { return c.count }

// Inc increments the counter in a thread-safe way.
func (c *SafeCounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

// Value returns the current count in a thread-safe way.
func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

// ConcurrentIncrement launches n goroutines that each call inc() once,
// waits for all of them, then returns the final value.
func ConcurrentIncrement(n int, counter Counter) {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Inc()
		}()
	}
	wg.Wait()
}
