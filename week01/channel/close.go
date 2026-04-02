package channel

// Generate sends integers 0, 1, 2, … into the returned channel
// until the done channel is closed, then closes the output channel.
func Generate(done <-chan struct{}) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 0; ; i++ {
			select {
			case out <- i:
			case <-done:
				return
			}
		}
	}()
	return out
}

// Take reads at most n values from src and returns them.
func Take(src <-chan int, n int) []int {
	result := make([]int, 0, n)
	for v := range src {
		result = append(result, v)
		if len(result) == n {
			break
		}
	}
	return result
}

// Producer sends values into a channel, then closes it.
// Demonstrates the "only the sender closes" convention.
func Producer(values []int) <-chan int {
	ch := make(chan int, len(values))
	go func() {
		for _, v := range values {
			ch <- v
		}
		close(ch) // only the sender closes the channel
	}()
	return ch
}

// Consumer reads all values from ch using a range loop.
// range automatically exits when the channel is closed.
func Consumer(ch <-chan int) []int {
	var result []int
	for v := range ch {
		result = append(result, v)
	}
	return result
}

// CheckClosed demonstrates the two-value receive idiom to detect closure.
func CheckClosed(ch <-chan int) (int, bool) {
	v, ok := <-ch
	return v, ok // ok == false means the channel was closed
}
