package channel

import (
	"fmt"
	"sync"
	"time"
)

// Task 表示一个任务
type Task struct {
	ID   int
	Data string
}

// Result 表示处理结果
type Result struct {
	TaskID int
	Output string
}

// TaskProducer 生产任务
func TaskProducer(id int, tasks chan<- Task, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 5; i++ {
		task := Task{
			ID:   id*100 + i,
			Data: fmt.Sprintf("Task-%d-%d", id, i),
		}
		tasks <- task
		fmt.Printf("Producer %d produced: %v\n", id, task)
		time.Sleep(100 * time.Millisecond)
	}
}

// TaskConsumer 消费任务
func TaskConsumer(id int, tasks <-chan Task, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		// 模拟处理
		time.Sleep(200 * time.Millisecond)
		result := Result{
			TaskID: task.ID,
			Output: fmt.Sprintf("Processed %s by Consumer %d", task.Data, id),
		}
		results <- result
		fmt.Printf("Consumer %d processed: %s\n", id, task.Data)
	}
}

func main() {
	const (
		numProducers = 2
		numConsumers = 3
		bufferSize   = 10
	)

	tasks := make(chan Task, bufferSize)
	results := make(chan Result, bufferSize)

	var producerWg sync.WaitGroup
	var consumerWg sync.WaitGroup

	// 启动生产者
	for i := 0; i < numProducers; i++ {
		producerWg.Add(1)
		go TaskProducer(i, tasks, &producerWg)
	}

	// 启动消费者
	for i := 0; i < numConsumers; i++ {
		consumerWg.Add(1)
		go TaskConsumer(i, tasks, results, &consumerWg)
	}

	// 等待生产者完成，然后关闭 tasks channel
	go func() {
		producerWg.Wait()
		close(tasks)
	}()

	// 等待消费者完成，然后关闭 results channel
	go func() {
		consumerWg.Wait()
		close(results)
	}()

	// 收集结果
	var resultCount int
	for range results {
		resultCount++
	}

	fmt.Printf("\nTotal tasks processed: %d\n", resultCount)
}
