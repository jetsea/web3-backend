package fanout

import (
	"context"
	"sync"
)

// FanOut distributes input to multiple workers concurrently
// e.g. querying multiple databases or exchanges
func FanOut[T any, R any](
	ctx context.Context,
	input T,
	workers int,
	fn func(context.Context, T) (R, error),
) []Result[R] {
	results := make([]Result[R], workers)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			result, err := fn(ctx, input)
			results[idx] = Result[R]{
				Index: idx,
				Data:  result,
				Error: err,
			}
		}(i)
	}

	wg.Wait()
	return results
}

// FanIn aggregates results from multiple channels
// e.g. collecting results from multiple exchanges or databases
func FanIn[T any](ctx context.Context, channels ...<-chan Result[T]) <-chan Result[T] {
	out := make(chan Result[T], len(channels))
	var wg sync.WaitGroup

	wg.Add(len(channels))

	for _, ch := range channels {
		go func(c <-chan Result[T]) {
			defer wg.Done()
			for {
				select {
				case result, ok := <-c:
					if !ok {
						return
					}
					out <- result
				case <-ctx.Done():
					return
				}
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// Result represents the result of a fan-out operation
type Result[T any] struct {
	Index int
	Data  T
	Error error
}

// ConcurrentQuery executes multiple queries concurrently and aggregates results
func ConcurrentQuery[T any, R any](
	ctx context.Context,
	inputs []T,
	fn func(context.Context, T) (R, error),
) ([]R, error) {
	results := make([]R, len(inputs))
	errChan := make(chan error, len(inputs))
	var wg sync.WaitGroup

	for i, input := range inputs {
		wg.Add(1)
		go func(idx int, in T) {
			defer wg.Done()
			result, err := fn(ctx, in)
			if err != nil {
				errChan <- err
				return
			}
			results[idx] = result
		}(i, input)
	}

	wg.Wait()
	close(errChan)

	// Return first error if any
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

// AggregatePrices aggregates prices from multiple sources
func AggregatePrices(ctx context.Context, priceChannels ...<-chan float64) <-chan AggregatedPrice {
	out := make(chan AggregatedPrice, 1)
	go func() {
		defer close(out)

		var wg sync.WaitGroup
		var mu sync.Mutex
		var prices []float64

		for _, ch := range priceChannels {
			wg.Add(1)
			go func(c <-chan float64) {
				defer wg.Done()
				select {
				case price, ok := <-c:
					if ok {
						mu.Lock()
						prices = append(prices, price)
						mu.Unlock()
					}
				case <-ctx.Done():
					return
				}
			}(ch)
		}

		wg.Wait()

		agg := AggregatedPrice{
			Prices: prices,
			Valid:  len(prices) > 0,
		}
		if agg.Valid {
			agg.Average = calculateAverage(prices)
			agg.Min = min(prices)
			agg.Max = max(prices)
		}

		out <- agg
	}()

	return out
}

// AggregatedPrice represents aggregated price data
type AggregatedPrice struct {
	Prices  []float64
	Average float64
	Min     float64
	Max     float64
	Valid   bool
}

func calculateAverage(prices []float64) float64 {
	if len(prices) == 0 {
		return 0
	}
	sum := 0.0
	for _, p := range prices {
		sum += p
	}
	return sum / float64(len(prices))
}

func min(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	m := values[0]
	for _, v := range values[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

func max(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	m := values[0]
	for _, v := range values[1:] {
		if v > m {
			m = v
		}
	}
	return m
}
