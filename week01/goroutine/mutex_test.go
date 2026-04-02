package goroutine

import "testing"

// TestUnsafeCounter demonstrates that UnsafeCounter is not goroutine-safe.
// It may fail intermittently due to data races.
func TestUnsafeCounter(t *testing.T) {
	const n = 1000
	c := &UnsafeCounter{}

	ConcurrentIncrement(n, c)

	// UnsafeCounter is not goroutine-safe, so the result may be less than n
	// This test demonstrates the data race issue
	if got := c.Value(); got != n {
		t.Logf("UnsafeCounter.Value() = %d, want %d (demonstrating data race)", got, n)
	}
}
func TestSafeCounter(t *testing.T) {
	const n = 1000
	c := &SafeCounter{}

	ConcurrentIncrement(n, c)

	if got := c.Value(); got != n {
		t.Errorf("SafeCounter.Value() = %d, want %d", got, n)
	}
}

// TestUnsafeCounter_Sequential checks sequential increment correctness.
func TestUnsafeCounter_Sequential(t *testing.T) {
	c := &UnsafeCounter{}
	for i := 0; i < 10; i++ {
		c.Inc()
	}
	if c.Value() != 10 {
		t.Errorf("expected 10, got %d", c.Value())
	}
}

// TestSafeCounter_Sequential checks sequential increment correctness.
func TestSafeCounter_Sequential(t *testing.T) {
	c := &SafeCounter{}
	for i := 0; i < 10; i++ {
		c.Inc()
	}
	if c.Value() != 10 {
		t.Errorf("expected 10, got %d", c.Value())
	}
}
