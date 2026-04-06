package errors

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	err := New(ErrorCodeInvalidRequest, "bad request")

	if err.Code != ErrorCodeInvalidRequest {
		t.Errorf("Expected code %s, got %s", ErrorCodeInvalidRequest, err.Code)
	}

	if err.Message != "bad request" {
		t.Errorf("Expected message 'bad request', got '%s'", err.Message)
	}

	if err.HTTPStatus != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, err.HTTPStatus)
	}
}

func TestError(t *testing.T) {
	err := New(ErrorCodeNotFound, "resource not found")
	expected := "[NOT_FOUND] resource not found"

	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}
}

func TestError_WithUnderlyingError(t *testing.T) {
	underlying := errors.New("database error")
	err := Wrap(underlying, ErrorCodeInternal, "failed to save")

	expected := "[INTERNAL_ERROR] failed to save: database error"
	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}
}

func TestUnwrap(t *testing.T) {
	underlying := errors.New("underlying error")
	err := Wrap(underlying, ErrorCodeInternal, "wrapper message")

	unwrapped := err.Unwrap()
	if unwrapped != underlying {
		t.Errorf("Expected unwrapped error to be the underlying error")
	}
}

func TestWithDetails(t *testing.T) {
	err := New(ErrorCodeInvalidRequest, "invalid input").
		WithDetails("Field 'email' is required")

	if err.Details != "Field 'email' is required" {
		t.Errorf("Expected details 'Field 'email' is required', got '%s'", err.Details)
	}
}

func TestInvalidRequest(t *testing.T) {
	err := InvalidRequest("missing required field")

	if err.Code != ErrorCodeInvalidRequest {
		t.Errorf("Expected code %s, got %s", ErrorCodeInvalidRequest, err.Code)
	}

	if err.HTTPStatus != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, err.HTTPStatus)
	}
}

func TestUnauthorized(t *testing.T) {
	err := Unauthorized("access denied")

	if err.Code != ErrorCodeUnauthorized {
		t.Errorf("Expected code %s, got %s", ErrorCodeUnauthorized, err.Code)
	}

	if err.HTTPStatus != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, err.HTTPStatus)
	}
}

func TestNotFound(t *testing.T) {
	err := NotFound("user")

	if err.Message != "user not found" {
		t.Errorf("Expected message 'user not found', got '%s'", err.Message)
	}

	if err.HTTPStatus != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, err.HTTPStatus)
	}
}

func TestInsufficientBalance(t *testing.T) {
	err := InsufficientBalance(1.5, 0.8, "BTC")

	if err.Code != ErrorCodeInsufficientBalance {
		t.Errorf("Expected code %s, got %s", ErrorCodeInsufficientBalance, err.Code)
	}

	expectedDetails := "Required: 1.50 BTC, Available: 0.80 BTC"
	if err.Details != expectedDetails {
		t.Errorf("Expected details '%s', got '%s'", expectedDetails, err.Details)
	}

	if err.HTTPStatus != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, err.HTTPStatus)
	}
}

func TestInvalidAddress(t *testing.T) {
	address := "0xabc123"
	err := InvalidAddress(address)

	if err.Code != ErrorCodeInvalidAddress {
		t.Errorf("Expected code %s, got %s", ErrorCodeInvalidAddress, err.Code)
	}

	if err.Details != fmt.Sprintf("Address: %s", address) {
		t.Errorf("Expected details to contain address")
	}
}

func TestTransactionFailed(t *testing.T) {
	txHash := "0x123...abc"
	reason := "insufficient gas"
	err := TransactionFailed(txHash, reason)

	if err.Code != ErrorCodeTransactionFailed {
		t.Errorf("Expected code %s, got %s", ErrorCodeTransactionFailed, err.Code)
	}

	if err.HTTPStatus != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, err.HTTPStatus)
	}
}

func TestInternal(t *testing.T) {
	err := Internal("database connection failed")

	if err.Code != ErrorCodeInternal {
		t.Errorf("Expected code %s, got %s", ErrorCodeInternal, err.Code)
	}

	if err.HTTPStatus != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, err.HTTPStatus)
	}
}

func TestIsAppError(t *testing.T) {
	appErr := New(ErrorCodeInvalidRequest, "test")
	stdErr := errors.New("standard error")

	if !IsAppError(appErr) {
		t.Error("Expected true for AppError")
	}

	if IsAppError(stdErr) {
		t.Error("Expected false for standard error")
	}
}

func TestGetCode(t *testing.T) {
	appErr := New(ErrorCodeInvalidRequest, "test")
	stdErr := errors.New("standard error")

	if GetCode(appErr) != ErrorCodeInvalidRequest {
		t.Errorf("Expected code %s, got %s", ErrorCodeInvalidRequest, GetCode(appErr))
	}

	if GetCode(stdErr) != ErrorCodeInternal {
		t.Errorf("Expected default code %s, got %s", ErrorCodeInternal, GetCode(stdErr))
	}
}

func TestGetHTTPStatus(t *testing.T) {
	appErr := New(ErrorCodeNotFound, "test")
	stdErr := errors.New("standard error")

	if GetHTTPStatus(appErr) != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, GetHTTPStatus(appErr))
	}

	if GetHTTPStatus(stdErr) != http.StatusInternalServerError {
		t.Errorf("Expected default status %d, got %d", http.StatusInternalServerError, GetHTTPStatus(stdErr))
	}
}

func TestStatusCodeMapping(t *testing.T) {
	tests := []struct {
		code     ErrorCode
		expected int
	}{
		{ErrorCodeInvalidRequest, http.StatusBadRequest},
		{ErrorCodeUnauthorized, http.StatusUnauthorized},
		{ErrorCodeForbidden, http.StatusForbidden},
		{ErrorCodeNotFound, http.StatusNotFound},
		{ErrorCodeConflict, http.StatusConflict},
		{ErrorCodeInsufficientBalance, http.StatusBadRequest},
		{ErrorCodeTransactionFailed, http.StatusInternalServerError},
		{ErrorCodeServiceUnavailable, http.StatusServiceUnavailable},
	}

	for _, tt := range tests {
		t.Run(string(tt.code), func(t *testing.T) {
			err := New(tt.code, "test")
			if err.HTTPStatus != tt.expected {
				t.Errorf("Expected status %d, got %d", tt.expected, err.HTTPStatus)
			}
		})
	}
}

func TestErrorChaining(t *testing.T) {
	err1 := errors.New("level 1 error")
	err2 := Wrap(err1, ErrorCodeInvalidRequest, "level 2 error")
	err3 := Wrap(err2, ErrorCodeInternal, "level 3 error")

	// Unwrapping should return the chain
	if err3.Unwrap() != err2 {
		t.Error("Unwrap chain failed")
	}

	if err2.Unwrap() != err1 {
		t.Error("Unwrap chain failed at level 2")
	}
}
