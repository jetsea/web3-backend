package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/web3-backend/week02/pkg/errors"
	"github.com/web3-backend/week02/pkg/fanout"
	"github.com/web3-backend/week02/pkg/logger"
)

// PriceSource represents a DEX price source
type PriceSource struct {
	Name      string
	Available bool
}

// SwapRequest represents a token swap request
type SwapRequest struct {
	TokenIn           string
	TokenOut          string
	AmountIn          float64
	SlippageTolerance float64
	UserID            string
}

// PriceResponse represents a price query response
type PriceResponse struct {
	Source    string
	Price     float64
	AmountOut float64
	Available bool
	Error     error
}

func main() {
	// Initialize logger
	log, err := logger.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	log.Info("Starting DeFi swap aggregator")

	ctx := context.Background()

	// Define DEX sources
	sources := []PriceSource{
		{Name: "uniswap", Available: true},
		{Name: "sushiswap", Available: true},
		{Name: "curve", Available: true},
		{Name: "balancer", Available: true},
		{Name: "pancakeswap", Available: false}, // Simulate unavailable DEX
	}

	// Simulate swap requests
	swaps := generateSwapRequests(3)

	for _, swap := range swaps {
		processSwap(ctx, log, sources, swap)
	}

	log.Info("DeFi swap aggregator completed")
}

// generateSwapRequests generates sample swap requests
func generateSwapRequests(count int) []SwapRequest {
	swaps := make([]SwapRequest, count)
	pairs := []struct{ in, out string }{
		{"ETH", "USDT"},
		{"BTC", "USDT"},
		{"USDT", "ETH"},
	}

	for i := 0; i < count; i++ {
		pair := pairs[i%len(pairs)]
		swaps[i] = SwapRequest{
			TokenIn:           pair.in,
			TokenOut:          pair.out,
			AmountIn:          rand.Float64() * 5,
			SlippageTolerance: 0.05,
			UserID:            fmt.Sprintf("user-%d", rand.Intn(100)),
		}
	}

	return swaps
}

// processSwap demonstrates fan-out/fan-in pattern for price aggregation
func processSwap(ctx context.Context, log *logger.Logger, sources []PriceSource, swap SwapRequest) {
	log.Info("Processing swap request",
		logger.String("tokenIn", swap.TokenIn),
		logger.String("tokenOut", swap.TokenOut),
		logger.Float64("amountIn", swap.AmountIn),
	)

	// Fan-Out: Query multiple DEXes concurrently
	priceChannels := make([]<-chan PriceResponse, 0, len(sources))

	for _, source := range sources {
		ch := queryDEX(ctx, source, swap)
		priceChannels = append(priceChannels, ch)
	}

	// Fan-In: Aggregate results
	aggregatedPrices := aggregatePrices(ctx, priceChannels)

	log.Info("Price aggregation completed",
		logger.Int("total_sources", len(sources)),
		logger.Int("valid_responses", len(aggregatedPrices)),
	)

	// Find best price
	bestPrice, err := findBestPrice(aggregatedPrices, swap.SlippageTolerance)
	if err != nil {
		log.Error("Failed to find best price",
			logger.Err(err),
			logger.String("userID", swap.UserID),
		)
		return
	}

	// Execute swap on best DEX
	executeSwap(log, bestPrice, swap)
}

// queryDEX simulates querying a DEX for token price
func queryDEX(ctx context.Context, source PriceSource, swap SwapRequest) <-chan PriceResponse {
	ch := make(chan PriceResponse, 1)

	go func() {
		defer close(ch)

		// Simulate network delay
		time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)

		// Simulate source being unavailable
		if !source.Available {
			ch <- PriceResponse{
				Source:    source.Name,
				Available: false,
				Error:     fmt.Errorf("DEX %s unavailable", source.Name),
			}
			return
		}

		// Simulate price calculation
		basePrice := rand.Float64() * 2000
		priceVariation := basePrice * (rand.Float64()*0.1 - 0.05) // +/- 5%
		price := basePrice + priceVariation

		ch <- PriceResponse{
			Source:    source.Name,
			Price:     price,
			AmountOut: swap.AmountIn * price,
			Available: true,
		}
	}()

	return ch
}

// aggregatePrice aggregates prices from multiple channels
func aggregatePrices(ctx context.Context, channels []<-chan PriceResponse) []PriceResponse {
	results := make([]PriceResponse, 0, len(channels))

	for _, ch := range channels {
		select {
		case result := <-ch:
			if result.Available {
				results = append(results, result)
			}
		case <-ctx.Done():
			return results
		}
	}

	return results
}

// findBestPrice finds the best price among aggregated results
func findBestPrice(prices []PriceResponse, slippageTolerance float64) (PriceResponse, error) {
	if len(prices) == 0 {
		return PriceResponse{}, errors.NotFound("available price sources")
	}

	// Find the best price (maximum amount out)
	var best PriceResponse
	maxAmount := 0.0

	for _, price := range prices {
		if price.AmountOut > maxAmount {
			maxAmount = price.AmountOut
			best = price
		}
	}

	// Check slippage
	if len(prices) > 1 {
		avgAmount := 0.0
		for _, p := range prices {
			avgAmount += p.AmountOut
		}
		avgAmount /= float64(len(prices))

		slippage := (avgAmount - best.AmountOut) / avgAmount
		if slippage > slippageTolerance {
			return PriceResponse{}, errors.InvalidRequest("price deviation too high").
				WithDetails(fmt.Sprintf("Slippage: %.2f%%, Tolerance: %.2f%%", slippage*100, slippageTolerance*100))
		}
	}

	return best, nil
}

// executeSwap simulates executing a swap on the best DEX
func executeSwap(log *logger.Logger, price PriceResponse, swap SwapRequest) {
	log.Info("Executing swap",
		logger.String("source", price.Source),
		logger.String("tokenIn", swap.TokenIn),
		logger.String("tokenOut", swap.TokenOut),
		logger.Float64("amountIn", swap.AmountIn),
		logger.Float64("amountOut", price.AmountOut),
		logger.Float64("price", price.Price),
		logger.String("userID", swap.UserID),
	)

	// Simulate transaction delay
	time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)

	// Simulate occasional failures (5% chance)
	if rand.Float32() < 0.05 {
		log.Error("Swap failed",
			logger.String("source", price.Source),
			logger.String("userID", swap.UserID),
		)
		return
	}

	log.Info("Swap completed successfully",
		logger.String("source", price.Source),
		logger.String("userID", swap.UserID),
	)
}

// demonstrateConcurrentQuery demonstrates the ConcurrentQuery helper
func demonstrateConcurrentQuery(ctx context.Context, log *logger.Logger) {
	tokens := []string{"ETH", "BTC", "USDT"}

	prices, err := fanout.ConcurrentQuery(ctx, tokens, func(ctx context.Context, token string) (float64, error) {
		// Simulate API call
		time.Sleep(50 * time.Millisecond)
		return rand.Float64() * 3000, nil
	})

	if err != nil {
		log.Error("Concurrent query failed", logger.Err(err))
		return
	}

	log.Info("Fetched token prices", logger.Int("count", len(prices)))
}
