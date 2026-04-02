package goroutine

import "sync"

// SafeMap is a goroutine-safe map using sync.RWMutex.
// RWMutex allows multiple concurrent readers but only one writer,
// making it more efficient than a plain Mutex for read-heavy workloads.
type SafeMap struct {
	mu   sync.RWMutex
	data map[string]int
}

// NewSafeMap creates an initialised SafeMap.
func NewSafeMap() *SafeMap {
	return &SafeMap{data: make(map[string]int)}
}

// Set stores a key-value pair (writer — exclusive lock).
func (m *SafeMap) Set(key string, value int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

// Get retrieves a value by key (reader — shared lock).
func (m *SafeMap) Get(key string) (int, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[key]
	return v, ok
}

// Delete removes a key (writer — exclusive lock).
func (m *SafeMap) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

// Len returns the number of entries (reader — shared lock).
func (m *SafeMap) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}
