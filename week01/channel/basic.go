// Package channel demonstrates Go channel patterns.
package channel

// SendAndReceive shows the most basic unbuffered channel usage:
// the sender blocks until the receiver is ready, and vice-versa.
func SendAndReceive(value int) int {
	ch := make(chan int) // unbuffered

	go func() {
		ch <- value // blocks until someone receives
	}()

	return <-ch // blocks until the goroutine sends
}

// Pipeline passes a value through two transformation stages
// connected by unbuffered channels.
//
//	input → double → addOne → result
func Pipeline(input int) int {
	double := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			out <- (<-in) * 2
			close(out)
		}()
		return out
	}

	addOne := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			out <- <-in + 1
			close(out)
		}()
		return out
	}

	src := make(chan int, 1)
	src <- input
	close(src)

	return <-addOne(double(src))
}
