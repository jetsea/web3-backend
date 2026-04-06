package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/web3-backend/week02/pkg/errors"
	"github.com/web3-backend/week02/pkg/logger"
	"github.com/web3-backend/week02/pkg/pipeline"
	"github.com/web3-backend/week02/pkg/workerpool"
)

// DepositJob represents a cryptocurrency deposit job
type DepositJob struct {
	DepositID string
	UserID    string
	Amount    float64
	Currency  string
	TxHash    string
}

func (j *DepositJob) Execute() error {
	// Simulate blockchain monitoring
	time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)

	// Simulate occasional failures (10% chance)
	if rand.Float32() < 0.1 {
		return fmt.Errorf("blockchain timeout for tx %s", j.TxHash)
	}

	return nil
}

// Deposit represents a deposit record
type Deposit struct {
	DepositID     string
	UserID        string
	Amount        float64
	Currency      string
	TxHash        string
	Status        string
	Confirmations int
	UpdatedAt     time.Time
}

func main() {
	// Initialize logger
	log, err := logger.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	log.Info("Starting exchange deposit processor")

	ctx := context.Background()

	// Create worker pool
	pool := workerpool.New(5, 20, &zapAdapter{log})
	pool.Start(ctx)

	// Simulate incoming deposits
	deposits := generateDeposits(10)

	log.Info("Processing deposits", logger.Int("count", len(deposits)))

	// Submit jobs
	for _, deposit := range deposits {
		job := &DepositJob{
			DepositID: deposit.DepositID,
			UserID:    deposit.UserID,
			Amount:    deposit.Amount,
			Currency:  deposit.Currency,
			TxHash:    deposit.TxHash,
		}

		if err := pool.Add(job); err != nil {
			log.Error("Failed to submit deposit",
				logger.Err(err),
				logger.String("depositID", deposit.DepositID),
			)
		}
	}

	// Wait for processing
	time.Sleep(2 * time.Second)

	// Stop the pool
	pool.Stop()

	// Process successful deposits through pipeline
	log.Info("Running deposit pipeline")
	processDepositPipeline(ctx, log, deposits)

	log.Info("Exchange deposit processor completed")
}

// generateDeposits generates sample deposit records
func generateDeposits(count int) []Deposit {
	deposits := make([]Deposit, count)
	currencies := []string{"BTC", "ETH", "USDT"}

	for i := 0; i < count; i++ {
		deposits[i] = Deposit{
			DepositID: fmt.Sprintf("DEP-%d", i+1),
			UserID:    fmt.Sprintf("user-%d", rand.Intn(100)),
			Amount:    rand.Float64() * 2,
			Currency:  currencies[rand.Intn(len(currencies))],
			TxHash:    fmt.Sprintf("0x%x", rand.Int63()),
			Status:    "pending",
			UpdatedAt: time.Now(),
		}
	}

	return deposits
}

// processDepositPipeline demonstrates the pipeline pattern for deposit processing
func processDepositPipeline(ctx context.Context, log *logger.Logger, deposits []Deposit) {
	// Define pipeline stages
	p := pipeline.New(
		validateDeposit,
		confirmOnChain,
		updateBalance,
		sendNotification,
	)

	for _, deposit := range deposits {
		result, err := p.Execute(ctx, deposit)
		if err != nil {
			stageErr, ok := err.(pipeline.StageError)
			if ok {
				log.Error("Pipeline stage failed",
					logger.Int("stage", stageErr.Stage),
					logger.Err(stageErr),
					logger.String("depositID", deposit.DepositID),
				)
			}
			continue
		}

		processed := result
		log.Info("Deposit processed successfully",
			logger.String("depositID", processed.DepositID),
			logger.String("status", processed.Status),
		)
	}
}

// Pipeline stages

func validateDeposit(ctx context.Context, d Deposit) (Deposit, error) {
	log := logger.NewNop()

	if d.Amount <= 0 {
		return Deposit{}, errors.InvalidRequest("deposit amount must be positive").
			WithDetails(fmt.Sprintf("Received amount: %.2f", d.Amount))
	}

	log.Info("Deposit validated",
		logger.String("depositID", d.DepositID),
		logger.Float64("amount", d.Amount),
	)

	d.Status = "validated"
	return d, nil
}

func confirmOnChain(ctx context.Context, d Deposit) (Deposit, error) {
	log := logger.NewNop()

	// Simulate blockchain confirmation
	confirmations := rand.Intn(6) + 1
	if confirmations < 3 {
		return Deposit{}, errors.TransactionFailed(d.TxHash, "insufficient confirmations")
	}

	log.Info("Transaction confirmed on chain",
		logger.String("txHash", d.TxHash),
		logger.Int("confirmations", confirmations),
	)

	d.Status = "confirmed"
	d.Confirmations = confirmations
	return d, nil
}

func updateBalance(ctx context.Context, d Deposit) (Deposit, error) {
	log := logger.NewNop()

	// Simulate balance update
	// In real implementation, this would update the database
	log.Info("User balance updated",
		logger.String("userID", d.UserID),
		logger.Float64("amount", d.Amount),
		logger.String("currency", d.Currency),
	)

	d.Status = "completed"
	return d, nil
}

func sendNotification(ctx context.Context, d Deposit) (Deposit, error) {
	log := logger.NewNop()

	// Simulate notification sending
	log.Info("Notification sent to user",
		logger.String("userID", d.UserID),
		logger.String("type", "deposit_completed"),
	)

	return d, nil
}

// zapAdapter adapts logger.Logger to workerpool.Logger interface
type zapAdapter struct {
	logger *logger.Logger
}

func (a *zapAdapter) Info(msg string, fields ...interface{}) {
	logFields := convertFields(fields...)
	a.logger.Info(msg, logFields...)
}

func (a *zapAdapter) Error(msg string, fields ...interface{}) {
	logFields := convertFields(fields...)
	a.logger.Error(msg, logFields...)
}

func convertFields(fields ...interface{}) []logger.Field {
	result := make([]logger.Field, 0, len(fields)/2)

	for i := 0; i < len(fields)-1; i += 2 {
		if key, ok := fields[i].(string); ok {
			switch v := fields[i+1].(type) {
			case string:
				result = append(result, logger.String(key, v))
			case int:
				result = append(result, logger.Int(key, v))
			case float64:
				result = append(result, logger.Float64(key, v))
			case bool:
				result = append(result, logger.Bool(key, v))
			case error:
				result = append(result, logger.Err(v))
			default:
				result = append(result, logger.Any(key, v))
			}
		}
	}

	return result
}
