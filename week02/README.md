# Week 02: Enterprise-Grade Go Patterns

This project demonstrates production-ready Go patterns used in enterprise environments, specifically for blockchain and fintech applications.

## 📚 Topics Covered

1. **Worker Pool** - Concurrent task processing with controlled concurrency
2. **Pipeline** - Data processing pipeline pattern
3. **Fan-Out/Fan-In** - Concurrent aggregation pattern
4. **Error Handling** - Structured error types and handling
5. **Zap Logger** - Structured logging with zap
6. **Enterprise Examples** - Real-world use cases from exchanges, DeFi protocols, and wallet services

## 🏗️ Project Structure

```
week02/
├── pkg/
│   ├── workerpool/    # Worker pool implementation
│   ├── pipeline/      # Data pipeline implementation
│   ├── fanout/        # Fan-Out/Fan-In patterns
│   ├── errors/        # Custom error types
│   └── logger/        # Zap logger wrapper
└── examples/
    ├── exchange/      # Exchange deposit processing
    ├── defi/          # DeFi swap execution
    └── wallet/        # Multi-chain wallet transfer
```

## 🚀 Running Examples

```bash
# Initialize Go module
go mod tidy

# Run tests
go test ./...

# Run examples
go run examples/exchange/main.go
go run examples/defi/main.go
go run examples/wallet/main.go
```

## 📖 Enterprise Use Cases

### Exchange Deposit Processing
Demonstrates worker pool for processing cryptocurrency deposits with blockchain monitoring.

### DeFi Swap Execution
Shows fan-out/fan-in pattern for querying multiple DEXs and aggregating prices.

### Multi-Chain Wallet Transfer
Illustrates pipeline pattern for signing, broadcasting, and confirming transactions across multiple blockchains.
