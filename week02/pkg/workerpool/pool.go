package workerpool

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// ErrQueueFull is returned when the job queue is at capacity
var ErrQueueFull = errors.New("job queue is full")

// Job represents a unit of work to be processed
type Job interface {
	Execute() error
}

// Result represents the outcome of a job execution
type Result struct {
	JobID     string
	Success   bool
	Error     error
	Duration  time.Duration
	Timestamp time.Time
}

// Pool implements a worker pool for concurrent job processing
type Pool struct {
	workers        int
	jobChan        chan Job
	resultChan     chan Result
	wg             sync.WaitGroup
	logger         Logger
	runningWorkers int32
}

// Logger defines the logging interface
type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
}

// New creates a new worker pool
func New(workerCount, queueSize int, logger Logger) *Pool {
	return &Pool{
		workers:    workerCount,
		jobChan:    make(chan Job, queueSize),
		resultChan: make(chan Result, queueSize),
		logger:     logger,
	}
}

// Start begins processing jobs with the worker pool
func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker(ctx, i)
	}
}

// worker processes jobs from the queue
func (p *Pool) worker(ctx context.Context, id int) {
	defer p.wg.Done()

	for {
		select {
		case job, ok := <-p.jobChan:
			if !ok {
				return
			}
			p.processJob(ctx, id, job)

		case <-ctx.Done():
			p.logger.Info("Worker shutting down", "workerID", id)
			return
		}
	}
}

// processJob handles the execution of a single job
func (p *Pool) processJob(ctx context.Context, workerID int, job Job) {
	atomic.AddInt32(&p.runningWorkers, 1)
	defer atomic.AddInt32(&p.runningWorkers, -1)

	start := time.Now()
	err := job.Execute()
	duration := time.Since(start)

	result := Result{
		Success:   err == nil,
		Error:     err,
		Duration:  duration,
		Timestamp: time.Now(),
	}

	if err != nil {
		p.logger.Error("Job failed", "workerID", workerID, "error", err, "duration", duration)
	} else {
		p.logger.Info("Job completed", "workerID", workerID, "duration", duration)
	}

	select {
	case p.resultChan <- result:
	case <-ctx.Done():
	}
}

// Add adds a job to the queue
func (p *Pool) Add(job Job) error {
	select {
	case p.jobChan <- job:
		return nil
	default: //non-blocking
		return ErrQueueFull
	}
}

// Results returns a channel for job results
func (p *Pool) Results() <-chan Result {
	return p.resultChan
}

// Stop gracefully shuts down the worker pool
func (p *Pool) Stop() {
	close(p.jobChan)
	p.wg.Wait()
	close(p.resultChan)
}

// Stats returns pool statistics
type Stats struct {
	Workers       int
	QueueSize     int
	ActiveWorkers int
}

// Stats returns current pool statistics
func (p *Pool) Stats() Stats {
	return Stats{
		Workers:       p.workers,
		QueueSize:     len(p.jobChan),
		ActiveWorkers: int(atomic.LoadInt32(&p.runningWorkers)),
	}
}
