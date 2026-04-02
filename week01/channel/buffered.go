package channel

// BufferedFill sends n integers into a buffered channel and returns it.
// Because the channel is buffered, the sends do not block.
func BufferedFill(n int) <-chan int {
	ch := make(chan int, n)
	for i := 0; i < n; i++ {
		ch <- i
	}
	close(ch)
	return ch
}

// Drain reads all values from a channel and returns them as a slice.
func Drain(ch <-chan int) []int {
	var out []int
	for v := range ch {
		out = append(out, v)
	}
	return out
}

// BoundedProducer sends values 0..n-1 into a buffered channel of
// capacity cap, then closes it.  Use a separate goroutine so the
// caller controls when to start consuming.
func BoundedProducer(n, cap int) <-chan int {
	ch := make(chan int, cap)
	go func() {
		for i := 0; i < n; i++ {
			ch <- i
		}
		close(ch)
	}()
	return ch
}
