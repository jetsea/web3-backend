package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode represents the type of error
type ErrorCode string

const (
	// ErrorCodeInvalidRequest represents a bad request error
	ErrorCodeInvalidRequest ErrorCode = "INVALID_REQUEST"

	// ErrorCodeUnauthorized represents an unauthorized error
	ErrorCodeUnauthorized ErrorCode = "UNAUTHORIZED"

	// ErrorCodeForbidden represents a forbidden error
	ErrorCodeForbidden ErrorCode = "FORBIDDEN"

	// ErrorCodeNotFound represents a not found error
	ErrorCodeNotFound ErrorCode = "NOT_FOUND"

	// ErrorCodeConflict represents a conflict error
	ErrorCodeConflict ErrorCode = "CONFLICT"

	// ErrorCodeInternal represents an internal server error
	ErrorCodeInternal ErrorCode = "INTERNAL_ERROR"

	// ErrorCodeServiceUnavailable represents a service unavailable error
	ErrorCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"

	// ErrorCodeInsufficientBalance represents insufficient balance
	ErrorCodeInsufficientBalance ErrorCode = "INSUFFICIENT_BALANCE"

	// ErrorCodeInvalidAddress represents invalid address
	ErrorCodeInvalidAddress ErrorCode = "INVALID_ADDRESS"

	// ErrorCodeTransactionFailed represents transaction failure
	ErrorCodeTransactionFailed ErrorCode = "TRANSACTION_FAILED"
)

// AppError represents an application error with structured information
type AppError struct {
	Code       ErrorCode
	Message    string
	Details    string
	HTTPStatus int
	Err        error
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: statusCodeForCode(code),
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: statusCodeForCode(code),
		Err:        err,
	}
}

// WithDetails adds details to an error
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// InvalidRequest creates a bad request error
func InvalidRequest(message string) *AppError {
	return New(ErrorCodeInvalidRequest, message)
}

// Unauthorized creates an unauthorized error
func Unauthorized(message string) *AppError {
	return New(ErrorCodeUnauthorized, message)
}

// NotFound creates a not found error
func NotFound(resource string) *AppError {
	return New(ErrorCodeNotFound, fmt.Sprintf("%s not found", resource))
}

// InsufficientBalance creates an insufficient balance error
func InsufficientBalance(required, available float64, currency string) *AppError {
	return New(
		ErrorCodeInsufficientBalance,
		fmt.Sprintf("Insufficient %s balance", currency),
	).WithDetails(fmt.Sprintf("Required: %.2f %s, Available: %.2f %s", required, currency, available, currency))
}

// InvalidAddress creates an invalid address error
func InvalidAddress(address string) *AppError {
	return New(
		ErrorCodeInvalidAddress,
		"Invalid address format",
	).WithDetails(fmt.Sprintf("Address: %s", address))
}

// TransactionFailed creates a transaction failed error
func TransactionFailed(txHash string, reason string) *AppError {
	return New(
		ErrorCodeTransactionFailed,
		"Transaction failed",
	).WithDetails(fmt.Sprintf("Transaction Hash: %s, Reason: %s", txHash, reason))
}

// Internal creates an internal server error
func Internal(message string) *AppError {
	return New(ErrorCodeInternal, message)
}

// statusCodeForCode maps error codes to HTTP status codes
func statusCodeForCode(code ErrorCode) int {
	switch code {
	case ErrorCodeInvalidRequest:
		return http.StatusBadRequest
	case ErrorCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrorCodeForbidden:
		return http.StatusForbidden
	case ErrorCodeNotFound:
		return http.StatusNotFound
	case ErrorCodeConflict:
		return http.StatusConflict
	case ErrorCodeInsufficientBalance, ErrorCodeInvalidAddress:
		return http.StatusBadRequest
	case ErrorCodeTransactionFailed:
		return http.StatusInternalServerError
	case ErrorCodeServiceUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetCode extracts the error code from an error
func GetCode(err error) ErrorCode {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code
	}
	return ErrorCodeInternal
}

// GetHTTPStatus extracts the HTTP status from an error
func GetHTTPStatus(err error) int {
	if appErr, ok := err.(*AppError); ok {
		return appErr.HTTPStatus
	}
	return http.StatusInternalServerError
}
