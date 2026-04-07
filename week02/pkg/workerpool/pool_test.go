package workerpool

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

// mockLogger implements Logger for testing
type mockLogger struct {
	infoMessages  []string
	errorMessages []string
}

func (m *mockLogger) Info(msg string, fields ...interface{}) {
	m.infoMessages = append(m.infoMessages, msg)
}

func (m *mockLogger) Error(msg string, fields ...interface{}) {
	m.errorMessages = append(m.errorMessages, msg)
}

// testJob implements Job for testing
type testJob struct {
	name      string
	duration  time.Duration
	shouldErr bool
	executeFn func() error
}

// implement Execute to satisfy Job interface
func (j *testJob) Execute() error {
	time.Sleep(j.duration)
	if j.executeFn != nil {
		return j.executeFn()
	}
	if j.shouldErr {
		return errors.New("job failed")
	}
	return nil
}

func TestPool_BasicOperation(t *testing.T) {
	logger := &mockLogger{}
	//3 workers, queue size 10
	pool := New(3, 10, logger)
	ctx := context.Background()
	// Start the worker pool
	pool.Start(ctx)

	var jobsProcessed int32

	for i := 0; i < 5; i++ {

		job := &testJob{
			name:     "test-job",
			duration: 10 * time.Millisecond,
			executeFn: func() error {
				atomic.AddInt32(&jobsProcessed, 1)
				return nil
			},
		}
		// Add 5 jobs to the pool
		if err := pool.Add(job); err != nil {
			t.Fatalf("Job chan is full! New job failed: %v", err)
		}
	}

	// Stop the worker pool
	pool.Stop()

	if jobsProcessed != 5 {
		t.Errorf("Expected 5 jobs processed, got %d", jobsProcessed)
	}
}

func TestPool_Errors(t *testing.T) {
	logger := &mockLogger{}
	//2 workers, queue size 5
	pool := New(2, 5, logger)
	ctx := context.Background()

	pool.Start(ctx)

	successJob := &testJob{name: "success", duration: 10 * time.Millisecond}
	errorJob := &testJob{name: "error", duration: 10 * time.Millisecond, shouldErr: true}

	pool.Add(successJob)
	pool.Add(errorJob)

	pool.Stop()

	t.Log(logger.errorMessages)
	if len(logger.errorMessages) == 0 {
		t.Error("Expected error messages, got none")
	}
}

func TestPool_ContextCancellation(t *testing.T) {
	logger := &mockLogger{}
	//2 workers, queue size 10
	pool := New(2, 10, logger)
	ctx, cancel := context.WithCancel(context.Background())

	pool.Start(ctx)

	cancel()
	//sleeping ensures workers are running and select on ctx.Done() instead of jobChan closed
	time.Sleep(50 * time.Millisecond)

	pool.Stop()

	t.Log(logger.infoMessages)
	// Verify workers shut down gracefully
	if len(logger.infoMessages) == 0 {
		t.Error("Expected shutdown messages")
	}
}

func TestPool_QueueFull(t *testing.T) {
	logger := &mockLogger{}
	//1 worker, queue size 2
	pool := New(1, 2, logger)
	ctx := context.Background()

	pool.Start(ctx)

	// Fill queue
	for i := 0; i < 2; i++ {
		job := &testJob{
			name:     "blocking",
			duration: 500 * time.Millisecond,
		}
		pool.Add(job)
	}

	// This should fail because queue is full
	rejectedJob := &testJob{name: "rejected", duration: 10 * time.Millisecond}
	err := pool.Add(rejectedJob)
	if err != ErrQueueFull {
		t.Errorf("Expected ErrQueueFull, got %v", err)
	}

	pool.Stop()
}

func TestPool_Results(t *testing.T) {
	logger := &mockLogger{}
	//2 workers, queue size 10
	pool := New(2, 10, logger)
	ctx := context.Background()

	pool.Start(ctx)

	job := &testJob{name: "test", duration: 10 * time.Millisecond}
	pool.Add(job)

	//sleep for job to be processed before stopping the pool

	pool.Stop()

	//ok=true: chan is open or closed but not drained
	//ok=false: chan is closed and drained
	_, ok := <-pool.Results()
	if !ok {
		t.Error("Expected there to be results, got none")
	}
	_, ok = <-pool.Results()
	if ok {
		t.Error("Expected results channel to be closed")
	}
}

func TestPool_Stats(t *testing.T) {
	logger := &mockLogger{}
	pool := New(3, 10, logger)

	stats := pool.Stats()
	if stats.Workers != 3 && stats.ActiveWorkers != 0 && stats.QueueSize != 0 {
		t.Errorf("Expected 3 workers, got %d workers, %d active workers, %d queue size", stats.Workers, stats.ActiveWorkers, stats.QueueSize)
	}

	// Submit jobs
	pool.Start(context.Background())
	for i := 0; i < 5; i++ {
		job := &testJob{name: "test", duration: 10 * time.Millisecond}
		pool.Add(job)
	}

	stats = pool.Stats()
	//because of the sleep in testJob, all 5 jobs will be active at the same time
	if stats.QueueSize != 5 && stats.ActiveWorkers != 3 {
		t.Errorf("Expected queue size 5, got %d, Expected active workers 3, got %d ", stats.QueueSize, stats.ActiveWorkers)
	}

	pool.Stop()
}
