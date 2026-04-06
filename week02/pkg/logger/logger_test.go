package logger

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TestNewDevelopment tests creating a development logger
func TestNewDevelopment(t *testing.T) {
	logger, err := NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create development logger: %v", err)
	}
	defer logger.Sync()

	if logger == nil {
		t.Fatal("Expected non-nil logger")
	}
}

// TestNewProduction tests creating a production logger
func TestNewProduction(t *testing.T) {
	tmpFile := "/tmp/test.log"
	logger, err := NewProduction(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create production logger: %v", err)
	}
	defer logger.Sync()

	if logger == nil {
		t.Fatal("Expected non-nil logger")
	}

	logger.Info("production log", String("test", "value"))
}

// TestNewNop tests creating a no-op logger
func TestNewNop(t *testing.T) {
	logger := NewNop()

	if logger == nil {
		t.Fatal("Expected non-nil logger")
	}

	logger.Info("this should not be logged")
	logger.Error("this should also not be logged")
}

// TestLogger_With tests creating child logger with fields
func TestLogger_With(t *testing.T) {
	logger, _ := NewDevelopment()
	defer logger.Sync()

	childLogger := logger.With(String("service", "test"), Int("version", 1))

	if childLogger == nil {
		t.Fatal("Expected non-nil child logger")
	}

	childLogger.Info("child logger message")
}

// TestLoggingLevels tests different log levels
func TestLoggingLevels(t *testing.T) {
	// Create a custom core that captures log entries
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	var buf bytes.Buffer
	writeSyncer := zapcore.AddSync(&buf)
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core)
	wrappedLogger := &Logger{
		logger: logger,
		sugar:  logger.Sugar(),
	}
	defer wrappedLogger.Sync()

	wrappedLogger.Debug("debug message", String("key", "debug"))
	wrappedLogger.Info("info message", String("key", "info"))
	wrappedLogger.Warn("warn message", String("key", "warn"))
	wrappedLogger.Error("error message", String("key", "error"))

	output := buf.String()
	if output == "" {
		t.Error("Expected log output, got empty string")
	}
}

// TestLoggingFields tests logging with different field types
func TestLoggingFields(t *testing.T) {
	logger, _ := NewDevelopment()
	defer logger.Sync()

	logger.Info("test message",
		String("stringField", "testValue"),
		Int("intField", 42),
		Int64("int64Field", 123456789),
		Float64("floatField", 3.14),
		Bool("boolField", true),
		Any("anyField", map[string]string{"key": "value"}),
		Err(errors.New("test error")),
	)
}

// TestNamed tests creating a named logger
func TestNamed(t *testing.T) {
	logger, _ := NewDevelopment()
	defer logger.Sync()

	namedLogger := logger.Named("service.test")
	namedLogger.Info("named logger message")

	if namedLogger == nil {
		t.Fatal("Expected non-nil named logger")
	}
}

// TestWithFields tests With method chaining
func TestWithFields(t *testing.T) {
	logger, _ := NewDevelopment()
	defer logger.Sync()

	fieldLogger := logger.With(
		String("component", "test"),
		Int("iteration", 1),
	)

	fieldLogger.Info("message with fields")
}

// TestFieldConstructors tests field constructors
func TestFieldConstructors(t *testing.T) {
	tests := []struct {
		name  string
		field Field
	}{
		{"String", String("key", "value")},
		{"Int", Int("key", 42)},
		{"Int64", Int64("key", 123456789)},
		{"Float64", Float64("key", 3.14)},
		{"Bool", Bool("key", true)},
		{"Any", Any("key", "anyValue")},
		{"Err", Err(errors.New("test"))},
	}

	logger, _ := NewDevelopment()
	defer logger.Sync()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.Info("test", tt.field)
		})
	}
}

// TestGlobalLogger tests global logger functions
func TestGlobalLogger(t *testing.T) {
	// Reset default logger
	SetDefault(NewNop())

	Info("global info", String("test", "value"))
	Error("global error", Err(errors.New("test")))
	Debug("global debug", Int("count", 1))
	Warn("global warn", Bool("flag", true))
}

// TestSetDefault tests setting default logger
func TestSetDefault(t *testing.T) {
	logger, _ := NewDevelopment()
	defer logger.Sync()

	SetDefault(logger)

	Info("new default logger message")
}

// TestStdLogger tests stdlib-compatible logger
func TestStdLogger(t *testing.T) {
	logger, _ := NewDevelopment()
	defer logger.Sync()

	stdLogger := logger.Std()
	stdLogger.Println("std println")
	stdLogger.Printf("std printf: %s", "formatted")
}

// TestErrorLogging tests error field handling
func TestErrorLogging(t *testing.T) {
	logger, _ := NewDevelopment()
	defer logger.Sync()

	err := errors.New("test error with context")
	logger.Error("operation failed",
		Err(err),
		String("operation", "test"),
		Int("attempts", 3),
	)
}

// TestDurationField tests duration field
func TestDurationField(t *testing.T) {
	logger, _ := NewDevelopment()
	defer logger.Sync()

	duration := 150 * time.Millisecond
	logger.Info("operation completed",
		Duration("elapsed", duration),
		Int("items_processed", 100),
	)
}

// TestTimeField tests time field
func TestTimeField(t *testing.T) {
	logger, _ := NewDevelopment()
	defer logger.Sync()

	now := time.Now()
	logger.Info("timestamp test",
		Time("created_at", now),
	)
}

// TestWithChainedFields tests chaining With calls
func TestWithChainedFields(t *testing.T) {
	logger, _ := NewDevelopment()
	defer logger.Sync()

	l1 := logger.With(String("service", "api"))
	l2 := l1.With(Int("port", 8080))
	l3 := l2.With(Bool("secure", true))

	l3.Info("chained fields")
}

// TestInvalidLogLevel tests creating logger with invalid log level
func TestInvalidLogLevel(t *testing.T) {
	_, err := New(Config{
		Level:      "invalid",
		Production: false,
	})

	if err == nil {
		t.Error("Expected error for invalid log level")
	}
}

// TestLoggerSync tests syncing the logger
func TestLoggerSync(t *testing.T) {
	logger, _ := NewDevelopment()

	logger.Info("message before sync")
	err := logger.Sync()

	if err != nil {
		t.Errorf("Sync failed: %v", err)
	}
}
