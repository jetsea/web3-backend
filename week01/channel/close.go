package channel

// Push value into a channel unlimitedly until outside closes the channel.
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

// Push all values from the slice into a channel, then close it.
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

// Take reads at most n values from src and returns them as slice.
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

// Take all values from the channel until it's closed, and return them as slice.
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
