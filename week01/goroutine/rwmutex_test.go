package goroutine

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestSafeMap_SetGet(t *testing.T) {
	m := NewSafeMap()
	m.Set("btc", 60000)

	v, ok := m.Get("btc")
	if !ok || v != 60000 {
		t.Errorf("Get(btc) = %d, %v; want 60000, true", v, ok)
	}
}

func TestSafeMap_MissingKey(t *testing.T) {
	m := NewSafeMap()
	_, ok := m.Get("eth")
	if ok {
		t.Error("expected false for missing key")
	}
}

func TestSafeMap_Delete(t *testing.T) {
	m := NewSafeMap()
	m.Set("sol", 100)
	m.Delete("sol")
	_, ok := m.Get("sol")
	if ok {
		t.Error("key should have been deleted")
	}
}

func TestSafeMap_Len(t *testing.T) {
	m := NewSafeMap()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)
	if m.Len() != 3 {
		t.Errorf("Len() = %d, want 3", m.Len())
	}
}

func TestSafeMap_ConcurrentReadWrite(t *testing.T) {
	m := NewSafeMap()
	var wg sync.WaitGroup
	const n = 200

	// Writers
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(10)))
			key := "key"
			fmt.Printf("Writing %d\n", idx)
			m.Set(key, idx)
			fmt.Printf("Wrote %d\n", idx)
		}(i)
	}

	// Readers — should never panic or deadlock
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(10)))
			x, ok := m.Get("key")
			if !ok {
				t.Errorf("Get(key) = %d, %v; want true", x, ok)
			} else {
				fmt.Printf("Get value: %d\n", x)
			}
		}()
	}

	wg.Wait()
}
