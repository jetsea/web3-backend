// Package producer_consumer demonstrates the classic producer-consumer pattern
// using goroutines and channels.
package producer_consumer

import (
	"fmt"
	"sync"
	"time"
)

// Task represents a unit of work.
type Task struct {
	ID   int
	Data string
}

// Result holds the processed output of a Task.
type Result struct {
	TaskID int
	Output string
}

// Produce generates count tasks per producer and sends them to tasks.
// Each producer closes its own goroutine; the caller must close tasks
// when ALL producers are done (see Run).
func Produce(id, count int, tasks chan<- Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < count; i++ {
		task := Task{
			ID:   id*100 + i,
			Data: fmt.Sprintf("task-%d-%d", id, i),
		}
		tasks <- task
	}
}

// Consume reads tasks from the tasks channel until it is closed,
// processes each one, and sends the Result to results.
func Consume(id int, tasks <-chan Task, results chan<- Result, wg *sync.WaitGroup, delay time.Duration) {
	defer wg.Done()
	for task := range tasks {
		if delay > 0 {
			time.Sleep(delay)
		}
		results <- Result{
			TaskID: task.ID,
			Output: fmt.Sprintf("processed %s by consumer-%d", task.Data, id),
		}
	}
}

// Run wires producers and consumers together.
//   - numProducers: number of producers, each generates tasksPerProducer tasks
//   - numConsumers: number of consumers
//   - bufSize:      channel buffer size
//
// Returns all Result values.
func Run(numProducers, tasksPerProducer, numConsumers, bufSize int, consumerDelay time.Duration) []Result {
	tasks := make(chan Task, bufSize)
	results := make(chan Result, numProducers*tasksPerProducer)

	var producerWg sync.WaitGroup
	var consumerWg sync.WaitGroup

	// Start producers.
	for i := 0; i < numProducers; i++ {
		producerWg.Add(1)
		go Produce(i, tasksPerProducer, tasks, &producerWg)
	}

	// Start consumers.
	for i := 0; i < numConsumers; i++ {
		consumerWg.Add(1)
		go Consume(i, tasks, results, &consumerWg, consumerDelay)
	}

	// Close tasks once all producers finish.
	go func() {
		producerWg.Wait()
		close(tasks)
	}()

	// Close results once all consumers finish.
	go func() {
		consumerWg.Wait()
		close(results)
	}()

	var out []Result
	for r := range results {
		out = append(out, r)
	}
	return out
}
