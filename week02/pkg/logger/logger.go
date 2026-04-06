package logger

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger with structured logging capabilities
type Logger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

// Config represents logger configuration
type Config struct {
	Level      string
	Production bool
	OutputPath string
}

// New creates a new logger instance
func New(config Config) (*Logger, error) {
	var zapConfig zap.Config

	if config.Production {
		zapConfig = zap.NewProductionConfig()
		zapConfig.OutputPaths = []string{config.OutputPath, "stdout"}
		zapConfig.EncoderConfig.TimeKey = "timestamp"
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// Set log level
	if config.Level != "" {
		if err := zapConfig.Level.UnmarshalText([]byte(config.Level)); err != nil {
			return nil, fmt.Errorf("invalid log level: %w", err)
		}
	}

	logger, err := zapConfig.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	return &Logger{
		logger: logger,
		sugar:  logger.Sugar(),
	}, nil
}

// NewDevelopment creates a development logger with pretty output
func NewDevelopment() (*Logger, error) {
	return New(Config{
		Level:      "debug",
		Production: false,
	})
}

// NewProduction creates a production logger with JSON output
func NewProduction(outputPath string) (*Logger, error) {
	return New(Config{
		Level:      "info",
		Production: true,
		OutputPath: outputPath,
	})
}

// NewNop creates a no-op logger for testing
func NewNop() *Logger {
	logger := zap.NewNop()
	return &Logger{
		logger: logger,
		sugar:  logger.Sugar(),
	}
}

// With creates a child logger with additional fields
func (l *Logger) With(fields ...Field) *Logger {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = f.zapField
	}
	return &Logger{
		logger: l.logger.With(zapFields...),
		sugar:  l.sugar,
	}
}

// Debug logs a message at debug level
func (l *Logger) Debug(msg string, fields ...Field) {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = f.zapField
	}
	l.logger.Debug(msg, zapFields...)
}

// Info logs a message at info level
func (l *Logger) Info(msg string, fields ...Field) {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = f.zapField
	}
	l.logger.Info(msg, zapFields...)
}

// Warn logs a message at warning level
func (l *Logger) Warn(msg string, fields ...Field) {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = f.zapField
	}
	l.logger.Warn(msg, zapFields...)
}

// Error logs a message at error level
func (l *Logger) Error(msg string, fields ...Field) {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = f.zapField
	}
	l.logger.Error(msg, zapFields...)
}

// Fatal logs a message at fatal level and exits
func (l *Logger) Fatal(msg string, fields ...Field) {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = f.zapField
	}
	l.logger.Fatal(msg, zapFields...)
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.logger.Sync()
}

// Named creates a named logger
func (l *Logger) Named(name string) *Logger {
	return &Logger{
		logger: l.logger.Named(name),
		sugar:  l.sugar.Named(name),
	}
}

// Field represents a log field
type Field struct {
	zapField zap.Field
}

// String creates a string field
func String(key, value string) Field {
	return Field{zapField: zap.String(key, value)}
}

// Int creates an int field
func Int(key string, value int) Field {
	return Field{zapField: zap.Int(key, value)}
}

// Int64 creates an int64 field
func Int64(key string, value int64) Field {
	return Field{zapField: zap.Int64(key, value)}
}

// Float64 creates a float64 field
func Float64(key string, value float64) Field {
	return Field{zapField: zap.Float64(key, value)}
}

// Bool creates a bool field
func Bool(key string, value bool) Field {
	return Field{zapField: zap.Bool(key, value)}
}

// Any creates a field with any value
func Any(key string, value interface{}) Field {
	return Field{zapField: zap.Any(key, value)}
}

// Err creates an error field
func Err(err error) Field {
	return Field{zapField: zap.Error(err)}
}

// Duration creates a duration field
func Duration(key string, value time.Duration) Field {
	return Field{zapField: zap.Duration(key, value)}
}

// Time creates a time field
func Time(key string, value time.Time) Field {
	return Field{zapField: zap.Time(key, value)}
}

// Std returns a logger compatible with stdlib log interface
func (l *Logger) Std() *stdLogger {
	return &stdLogger{logger: l}
}

// stdLogger implements stdlib log interface
type stdLogger struct {
	logger *Logger
}

func (l *stdLogger) Println(args ...interface{}) {
	l.logger.Info(fmt.Sprint(args...))
}

func (l *stdLogger) Printf(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l *stdLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(fmt.Sprint(args...))
}

func (l *stdLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatal(fmt.Sprintf(format, args...))
}

// Initialize default global logger
var defaultLogger *Logger

func init() {
	var err error
	defaultLogger, err = NewDevelopment()
	if err != nil {
		panic(err)
	}
}

// SetDefault sets the default global logger
func SetDefault(l *Logger) {
	defaultLogger = l
}

// Info logs to the default logger
func Info(msg string, fields ...Field) {
	defaultLogger.Info(msg, fields...)
}

// Error logs to the default logger
func Error(msg string, fields ...Field) {
	defaultLogger.Error(msg, fields...)
}

// Debug logs to the default logger
func Debug(msg string, fields ...Field) {
	defaultLogger.Debug(msg, fields...)
}

// Warn logs to the default logger
func Warn(msg string, fields ...Field) {
	defaultLogger.Warn(msg, fields...)
}
