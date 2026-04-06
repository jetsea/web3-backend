package channel

// Fill the channel completely in one go.
// When data is smaller than the chan capacity, use this!
func BufferedFill(n int) <-chan int {
	ch := make(chan int, n)
	for i := 0; i < n; i++ {
		ch <- i
	}
	close(ch)
	return ch
}

// continuously sends values 0..n-1 to a buffered channel with capacity cap.
// When data is bigger than the chan capacity, use this!
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

// Drain reads all values from a channel and returns them as a slice.
func Drain(ch <-chan int) []int {
	var out []int
	for v := range ch {
		out = append(out, v)
	}
	return out
}
