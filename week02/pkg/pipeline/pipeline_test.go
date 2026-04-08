package pipeline

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestPipeline_BasicExecution(t *testing.T) {
	// Define stages: multiply by 2, add 1, multiply by 3
	stages := []Stage[int]{
		func(ctx context.Context, input int) (int, error) {
			return input * 2, nil
		},
		func(ctx context.Context, input int) (int, error) {
			return input + 1, nil
		},
		func(ctx context.Context, input int) (int, error) {
			return input * 3, nil
		},
	}

	p := New(stages...)
	ctx := context.Background()

	result, err := p.Execute(ctx, 5)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// (5 * 2 + 1) * 3 = 33
	expected := 33
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

func TestPipeline_EmptyPipeline(t *testing.T) {
	p := New[int]()
	ctx := context.Background()

	result, err := p.Execute(ctx, 42)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result != 42 {
		t.Errorf("Expected input to be returned unchanged, got %d", result)
	}
}

func TestPipeline_StageError(t *testing.T) {
	customErr := errors.New("stage error")

	stages := []Stage[int]{
		func(ctx context.Context, input int) (int, error) {
			return input * 2, nil
		},
		func(ctx context.Context, input int) (int, error) {
			return 0, customErr
		},
	}

	p := New(stages...)
	ctx := context.Background()

	_, err := p.Execute(ctx, 5)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var stageErr StageError
	if !errors.As(err, &stageErr) {
		t.Fatalf("Expected StageError, got %T", err)
	}

	if stageErr.Stage != 1 {
		t.Errorf("Expected stage 1, got %d", stageErr.Stage)
	}
}

func TestPipeline_ContextCancellation(t *testing.T) {
	stages := []Stage[int]{
		func(ctx context.Context, input int) (int, error) {
			time.Sleep(10 * time.Millisecond)
			return input * 2, nil
		},
		func(ctx context.Context, input int) (int, error) {
			time.Sleep(10 * time.Millisecond)
			return input + 1, nil
		},
	}

	p := New(stages...)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(5 * time.Millisecond)
		cancel()
	}()

	_, err := p.Execute(ctx, 5)
	if err == nil {
		t.Fatal("Expected context cancellation error")
	}

	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}

func TestPipeline_AddStage(t *testing.T) {
	p := New[int]()
	ctx := context.Background()

	p.AddStage(func(ctx context.Context, input int) (int, error) {
		return input * 2, nil
	})

	p.AddStage(func(ctx context.Context, input int) (int, error) {
		return input + 10, nil
	})

	result, err := p.Execute(ctx, 5)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// (5 * 2) + 10 = 20
	expected := 20
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}

	if p.StageCount() != 2 {
		t.Errorf("Expected 2 stages, got %d", p.StageCount())
	}
}

func TestPipeline_Batch(t *testing.T) {
	stages := []Stage[int]{
		func(ctx context.Context, input int) (int, error) {
			return input * 2, nil
		},
		func(ctx context.Context, input int) (int, error) {
			return input + 10, nil
		},
	}

	p := New(stages...)
	ctx := context.Background()

	inputs := []int{1, 2, 3, 4, 5}
	results, err := Batch(p, ctx, inputs)
	if err != nil {
		t.Fatalf("Batch failed: %v", err)
	}

	expected := []int{12, 14, 16, 18, 20}
	for i, result := range results {
		if result != expected[i] {
			t.Errorf("Index %d: expected %d, got %d", i, expected[i], result)
		}
	}
}

func TestPipeline_Filter(t *testing.T) {
	// Create a pipeline that has 2 stages:
	// a filter that only allows positive numbers,
	// and a stage that multiplies the input by 2
	p := New[int](
		Filter(func(n int) bool { return n > 0 }),
		func(ctx context.Context, input int) (int, error) {
			return input * 2, nil
		},
	)
	ctx := context.Background()

	// Positive number should pass
	result, err := p.Execute(ctx, 5)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if result != 10 {
		t.Errorf("Expected 10, got %d", result)
	}

	// Negative number should be filtered
	_, err = p.Execute(ctx, -5)
	t.Logf("Filter error: %v", err)
	if err == nil {
		t.Fatal("Expected filter error, got nil")
	}

	var filterErr FilterError
	if !errors.As(err, &filterErr) {
		t.Fatalf("Expected FilterError, got %T", err)
	}
}

func TestPipeline_StringPipeline(t *testing.T) {
	stages := []Stage[string]{
		func(ctx context.Context, input string) (string, error) {
			return input + " world", nil
		},
		func(ctx context.Context, input string) (string, error) {
			return "Hello " + input, nil
		},
	}

	p := New(stages...)
	ctx := context.Background()

	result, err := p.Execute(ctx, "")
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	expected := "Hello  world"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestPipeline_StructProcessing(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	stages := []Stage[Person]{
		func(ctx context.Context, p Person) (Person, error) {
			p.Age += 1
			return p, nil
		},
		func(ctx context.Context, p Person) (Person, error) {
			p.Name = "Mr. " + p.Name
			return p, nil
		},
	}

	p := New(stages...)
	ctx := context.Background()

	input := Person{Name: "John", Age: 30}
	result, err := p.Execute(ctx, input)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.Age != 31 {
		t.Errorf("Expected age 31, got %d", result.Age)
	}
	if result.Name != "Mr. John" {
		t.Errorf("Expected name %q, got %q", "Mr. John", result.Name)
	}
}

func TestMap(t *testing.T) {
	// Test Map function with pipeline integration
	ctx := context.Background()

	// Create a pipeline that processes integers
	intPipeline := New(
		func(ctx context.Context, input int) (int, error) {
			return input * 2, nil
		},
		func(ctx context.Context, input int) (int, error) {
			return input + 10, nil
		},
	)

	// Use Map to convert int result to string
	intToString := Map(func(i int) string {
		return fmt.Sprintf("Processed: %d", i)
	})

	// Execute pipeline and then use Map for type conversion
	intResult, err := intPipeline.Execute(ctx, 5)
	if err != nil {
		t.Fatalf("Pipeline execution failed: %v", err)
	}

	// Convert pipeline result to string using Map
	strResult, err := intToString(ctx, intResult)
	if err != nil {
		t.Errorf("Map failed: %v", err)
	}

	expected := "Processed: 20" // (5 * 2) + 10 = 20
	if strResult != expected {
		t.Errorf("Expected %q, got %q", expected, strResult)
	}

	// Test Map with struct transformation in pipeline workflow

	type Person struct {
		Name string
		Age  int
	}

	// Create a pipeline that processes Person structs
	personPipeline := New(
		func(ctx context.Context, p Person) (Person, error) {
			p.Age += 1
			return p, nil
		},
	)

	// Use Map to convert Person to string representation
	personToString := Map(func(p Person) string {
		return fmt.Sprintf("%s is %d years old", p.Name, p.Age)
	})

	// Execute pipeline and then use Map for type conversion
	inputPerson := Person{Name: "Alice", Age: 30}
	personResult, err := personPipeline.Execute(ctx, inputPerson)
	if err != nil {
		t.Fatalf("Pipeline execution failed: %v", err)
	}

	// Convert pipeline result to string using Map
	personStrResult, err := personToString(ctx, personResult)
	if err != nil {
		t.Errorf("Map failed: %v", err)
	}

	expectedPersonStr := "Alice is 31 years old"
	if personStrResult != expectedPersonStr {
		t.Errorf("Expected %q, got %q", expectedPersonStr, personStrResult)
	}
}
