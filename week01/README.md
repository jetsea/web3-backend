# Week 01: Go Concurrency Fundamentals

> **Goal**: Master Goroutine, Channel, and Context
> **Schedule**: Monday / Wednesday / Friday, 4 hours each
> **Deliverables**: Concurrent crawler + Producer-consumer model + HTTP server with timeout

## Project Structure

```
week01/
├── goroutine/                  # Monday: Goroutine basics
│   ├── basic.go                # Goroutine introduction
│   ├── basic_test.go
│   ├── waitgroup.go            # sync.WaitGroup usage
│   ├── waitgroup_test.go
│   ├── mutex.go                # sync.Mutex for race-free counter
│   ├── mutex_test.go
│   ├── rwmutex.go              # sync.RWMutex for read-heavy maps
│   └── rwmutex_test.go
├── channel/                    # Wednesday: Channel patterns
│   ├── basic.go                # Unbuffered channel
│   ├── basic_test.go
│   ├── buffered.go             # Buffered channel
│   ├── buffered_test.go
│   ├── select.go               # Select statement
│   ├── select_test.go
│   ├── close.go                # Closing channels, range
│   └── close_test.go
├── context/                    # Friday: Context package
│   ├── cancel.go               # context.WithCancel
│   ├── cancel_test.go
│   ├── timeout.go              # context.WithTimeout
│   ├── timeout_test.go
│   ├── value.go                # context.WithValue
│   └── value_test.go
├── exercises/
│   ├── crawler/                # Monday exercise: concurrent URL fetcher
│   │   ├── crawler.go
│   │   └── crawler_test.go
│   ├── producer_consumer/      # Wednesday exercise: producer-consumer model
│   │   ├── producer_consumer.go
│   │   └── producer_consumer_test.go
│   └── http_server/            # Friday exercise: HTTP server with context timeout
│       ├── server.go
│       └── server_test.go
└── leetcode/
    ├── two_sum.go              # LeetCode #1
    ├── two_sum_test.go
    ├── merge_sorted_lists.go   # LeetCode #21
    └── merge_sorted_lists_test.go
```

## Running

```bash
# Run all tests
go test ./...

# Run a specific package
go test ./goroutine/...
go test ./channel/...
go test ./context/...

# Run exercises
go run exercises/crawler/crawler.go
go run exercises/http_server/server.go

# Get code coverage
go test ./... -coverprofile=coverage
# View detail coverage(every function)
go tool cover -func coverage
```

## Key Concepts

| Topic | What You Learn |
|-------|---------------|
| Goroutine | Lightweight threads, `go` keyword |
| WaitGroup | Wait for multiple goroutines to finish |
| Mutex | Prevent race conditions |
| RWMutex | Optimized for read-heavy 4 |
| Channel | Communication between goroutines |
| Select | Listen on multiple channels |
| Context | Cancellation, timeout, value propagation |
