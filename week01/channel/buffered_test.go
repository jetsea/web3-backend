package channel

import "testing"

func TestBufferedFill(t *testing.T) {
	ch := BufferedFill(5)
	got := Drain(ch)
	if len(got) != 5 {
		t.Fatalf("len = %d, want 5", len(got))
	}
	for i, v := range got {
		if v != i {
			t.Errorf("got[%d] = %d, want %d", i, v, i)
		}
	}
}

func TestBoundedProducer(t *testing.T) {
	ch := BoundedProducer(10, 3)
	got := Drain(ch)
	if len(got) != 10 {
		t.Fatalf("len = %d, want 10", len(got))
	}
}

func TestBuffered_NeverBlocks(t *testing.T) {
	// Sends to a full channel should block, but filling up to capacity should not.
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	// At this point the buffer is full.  Reading one slot should unblock a send.
	<-ch
	ch <- 4 // should not block
}
