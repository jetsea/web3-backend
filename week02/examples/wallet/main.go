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

// Blockchain represents a blockchain network
type Blockchain struct {
	Name     string
	ChainID  int
	Explorer string
}

// SigningJob represents a transaction signing job
type SigningJob struct {
	TxID       string
	ChainID    int
	UnsignedTx string
}

func (j *SigningJob) Execute() error {
	// Simulate signing delay
	time.Sleep(time.Duration(30+rand.Intn(70)) * time.Millisecond)

	// Simulate occasional signing failures (5% chance)
	if rand.Float32() < 0.05 {
		return fmt.Errorf("signing failed for tx %s on chain %d", j.TxID, j.ChainID)
	}

	return nil
}

// Transfer represents a multi-chain transfer request
type Transfer struct {
	TransferID    string
	UserID        string
	FromChain     Blockchain
	ToChain       Blockchain
	FromAddress   string
	ToAddress     string
	Amount        float64
	Currency      string
	Status        string
	TxHash        string
	SignedTx      string
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

	log.Info("Starting multi-chain wallet service")

	ctx := context.Background()

	// Define supported blockchains
	chains := []Blockchain{
		{Name: "Ethereum", ChainID: 1, Explorer: "etherscan.io"},
		{Name: "Polygon", ChainID: 137, Explorer: "polygonscan.com"},
		{Name: "BSC", ChainID: 56, Explorer: "bscscan.com"},
		{Name: "Arbitrum", ChainID: 42161, Explorer: "arbiscan.io"},
	}

	// Generate transfer requests
	transfers := generateTransfers(5, chains)

	log.Info("Processing multi-chain transfers", logger.Int("count", len(transfers)))

	// Process each transfer
	for _, transfer := range transfers {
		processTransfer(ctx, log, transfer)
	}

	log.Info("Multi-chain wallet service completed")
}

// generateTransfers generates sample transfer requests
func generateTransfers(count int, chains []Blockchain) []Transfer {
	transfers := make([]Transfer, count)
	currencies := []string{"ETH", "USDT", "DAI"}

	for i := 0; i < count; i++ {
		fromChain := chains[rand.Intn(len(chains))]
		toChain := chains[rand.Intn(len(chains))]

		transfers[i] = Transfer{
			TransferID:  fmt.Sprintf("TXF-%d", i+1),
			UserID:      fmt.Sprintf("user-%d", rand.Intn(100)),
			FromChain:   fromChain,
			ToChain:     toChain,
			FromAddress: fmt.Sprintf("0x%x", rand.Int63()),
			ToAddress:   fmt.Sprintf("0x%x", rand.Int63()),
			Amount:      rand.Float64() * 10,
			Currency:    currencies[rand.Intn(len(currencies))],
			Status:      "pending",
			UpdatedAt:   time.Now(),
		}
	}

	return transfers
}

// processTransfer handles the complete transfer lifecycle
func processTransfer(ctx context.Context, log *logger.Logger, transfer Transfer) {
	log.Info("Processing transfer",
		logger.String("transferID", transfer.TransferID),
		logger.String("fromChain", transfer.FromChain.Name),
		logger.String("toChain", transfer.ToChain.Name),
		logger.Float64("amount", transfer.Amount),
	)

	// Validate transfer
	if err := validateTransfer(transfer); err != nil {
		log.Error("Transfer validation failed",
			logger.Err(err),
			logger.String("transferID", transfer.TransferID),
		)
		return
	}

	// Sign transaction using worker pool
	signedTx, err := signTransaction(ctx, transfer)
	if err != nil {
		log.Error("Transaction signing failed",
			logger.Err(err),
			logger.String("transferID", transfer.TransferID),
		)
		return
	}

	transfer.SignedTx = signedTx

	// Process through pipeline: sign → broadcast → confirm
	p := pipeline.New(
		broadcastToChain,
		waitForConfirmation,
		updateTransferStatus,
	)

	result, err := p.Execute(ctx, transfer)
	if err != nil {
		stageErr, ok := err.(pipeline.StageError)
		if ok {
			log.Error("Pipeline stage failed",
				logger.Int("stage", stageErr.Stage),
				logger.Err(stageErr),
				logger.String("transferID", transfer.TransferID),
			)
		}
		return
	}

	processed := result
	log.Info("Transfer completed successfully",
		logger.String("transferID", processed.TransferID),
		logger.String("status", processed.Status),
		logger.String("txHash", processed.TxHash),
		logger.Int("confirmations", processed.Confirmations),
	)
}

// validateTransfer validates transfer parameters
func validateTransfer(t Transfer) error {
	if t.Amount <= 0 {
		return errors.InvalidRequest("transfer amount must be positive").
			WithDetails(fmt.Sprintf("Amount: %.2f", t.Amount))
	}

	if t.FromAddress == "" || t.ToAddress == "" {
		return errors.InvalidRequest("both from and to addresses are required")
	}

	return nil
}

// signTransaction signs the transaction using worker pool
func signTransaction(ctx context.Context, transfer Transfer) (string, error) {
	pool := workerpool.New(2, 5, &zapAdapter{logger.NewNop()})
	defer pool.Stop()

	pool.Start(ctx)

	job := &SigningJob{
		TxID:       transfer.TransferID,
		ChainID:    transfer.FromChain.ChainID,
		UnsignedTx: "unsigned_tx_data",
	}

	if err := pool.Add(job); err != nil {
		return "", err
	}

	// Wait for signing to complete
	time.Sleep(200 * time.Millisecond)

	// Simulate signed transaction
	signedTx := fmt.Sprintf("0xsigned%x", rand.Int63())

	return signedTx, nil
}

// Pipeline stages

func broadcastToChain(ctx context.Context, t Transfer) (Transfer, error) {
	log := logger.NewNop()

	// Simulate broadcasting to chain
	time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)

	// Simulate occasional broadcast failures (10% chance)
	if rand.Float32() < 0.1 {
		return Transfer{}, errors.TransactionFailed(t.TransferID, "broadcast timeout").
			WithDetails(fmt.Sprintf("Chain: %s", t.FromChain.Name))
	}

	// Generate transaction hash
	txHash := fmt.Sprintf("0x%x", rand.Int63())
	t.TxHash = txHash
	t.Status = "broadcasted"

	log.Info("Transaction broadcasted",
		logger.String("transferID", t.TransferID),
		logger.String("chain", t.FromChain.Name),
		logger.String("txHash", txHash),
	)

	return t, nil
}

func waitForConfirmation(ctx context.Context, t Transfer) (Transfer, error) {
	log := logger.NewNop()

	// Simulate waiting for confirmations
	time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)

	// Simulate occasional confirmation failures (5% chance)
	if rand.Float32() < 0.05 {
		return Transfer{}, errors.TransactionFailed(t.TxHash, "insufficient confirmations").
			WithDetails(fmt.Sprintf("Required: 3, Received: %d", t.Confirmations))
	}

	t.Confirmations = rand.Intn(6) + 3
	t.Status = "confirmed"

	log.Info("Transaction confirmed",
		logger.String("transferID", t.TransferID),
		logger.String("txHash", t.TxHash),
		logger.Int("confirmations", t.Confirmations),
	)

	return t, nil
}

func updateTransferStatus(ctx context.Context, t Transfer) (Transfer, error) {
	log := logger.NewNop()

	// Simulate updating database
	t.Status = "completed"
	t.UpdatedAt = time.Now()

	log.Info("Transfer status updated",
		logger.String("transferID", t.TransferID),
		logger.String("status", t.Status),
	)

	return t, nil
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
