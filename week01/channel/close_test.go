package channel

import "testing"

func TestProducerConsumer(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	ch := Producer(values)
	got := Consumer(ch)

	if len(got) != len(values) {
		t.Fatalf("len = %d, want %d", len(got), len(values))
	}
	for i, v := range got {
		if v != values[i] {
			t.Errorf("got[%d] = %d, want %d", i, v, values[i])
		}
	}
}

func TestGenerate_Take(t *testing.T) {
	done := make(chan struct{})
	nums := Generate(done)

	got := Take(nums, 5)
	close(done)

	if len(got) != 5 {
		t.Fatalf("Take returned %d items, want 5", len(got))
	}
	for i, v := range got {
		if v != i {
			t.Errorf("got[%d] = %d, want %d", i, v, i)
		}
	}
}

func TestCheckClosed_OpenChannel(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 42

	v, ok := CheckClosed(ch)
	if !ok || v != 42 {
		t.Errorf("CheckClosed = (%d, %v), want (42, true)", v, ok)
	}
}

func TestCheckClosed_ClosedChannel(t *testing.T) {
	ch := make(chan int)
	close(ch)

	v, ok := CheckClosed(ch)
	if ok || v != 0 {
		t.Errorf("CheckClosed = (%d, %v), want (0, false) for closed channel", v, ok)
	}
}
