package pipeline

import (
	"context"
)

// Stage represents a single processing stage in the pipeline
type Stage[T any] func(ctx context.Context, input T) (T, error)

// Pipeline represents a data processing pipeline
type Pipeline[T any] struct {
	stages []Stage[T]
}

// New creates a new pipeline with the given stages
func New[T any](stages ...Stage[T]) *Pipeline[T] {
	return &Pipeline[T]{
		stages: stages,
	}
}

// Execute runs the pipeline with the given input through all stages
func (p *Pipeline[T]) Execute(ctx context.Context, input T) (T, error) {
	var err error
	current := input

	for i, stage := range p.stages {
		select {
		case <-ctx.Done():
			var zero T
			return zero, ctx.Err()
		default:
			current, err = stage(ctx, current)
			if err != nil {
				var zero T
				return zero, StageError{Stage: i, Err: err}
			}
		}
	}

	return current, nil
}

// ExecuteConcurrent runs the pipeline with the given input, running stages concurrently
// when possible (independent stages can run in parallel)
func (p *Pipeline[T]) ExecuteConcurrent(ctx context.Context, input T) (T, error) {
	// For simplicity, sequential execution. In production, you would
	// analyze stage dependencies and run independent stages concurrently
	return p.Execute(ctx, input)
}

// AddStage adds a new stage to the pipeline
func (p *Pipeline[T]) AddStage(stage Stage[T]) *Pipeline[T] {
	p.stages = append(p.stages, stage)
	return p
}

// StageCount returns the number of stages in the pipeline
func (p *Pipeline[T]) StageCount() int {
	return len(p.stages)
}

// StageError represents an error that occurred in a specific stage
type StageError struct {
	Stage int
	Err   error
}

func (e StageError) Error() string {
	return e.Err.Error()
}

func (e StageError) Unwrap() error {
	return e.Err
}

// Filter creates a stage that filters the input based on a predicate
func Filter[T any](predicate func(T) bool) Stage[T] {
	return func(ctx context.Context, input T) (T, error) {
		if !predicate(input) {
			var zero T
			return zero, FilterError("input filtered out")
		}
		return input, nil
	}
}

// Map creates a stage that transforms the input using a mapper function
func Map[T, U any](mapper func(T) U) func(ctx context.Context, input T) (U, error) {
	return func(ctx context.Context, input T) (U, error) {
		return mapper(input), nil
	}
}

// Batch creates a pipeline that processes multiple inputs through the same pipeline
func Batch[T any](p *Pipeline[T], ctx context.Context, inputs []T) ([]T, error) {
	results := make([]T, len(inputs))
	errChan := make(chan error, len(inputs))

	for i, input := range inputs {
		go func(idx int, in T) {
			result, err := p.Execute(ctx, in)
			if err != nil {
				errChan <- err
				return
			}
			results[idx] = result
			errChan <- nil
		}(i, input)
	}

	for range inputs {
		if err := <-errChan; err != nil {
			return nil, err
		}
	}

	return results, nil
}

// FilterError is returned when input is filtered out
type FilterError string

func (e FilterError) Error() string {
	return string(e)
}
