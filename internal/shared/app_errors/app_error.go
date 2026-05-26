package app_errors

import (
	"fmt"
	"net/http"

	"github.com/boilerplate/internal/shared/enums"
)

var errorCodeMap = map[string]int{
	ErrCodeValidation: http.StatusBadRequest,
	ErrCodeAuth:       http.StatusUnauthorized,
	ErrCodeForbidden:  http.StatusForbidden,
	ErrCodeNotFound:   http.StatusNotFound,
	ErrCodeConflict:   http.StatusConflict,
	ErrCodeRateLimit:  http.StatusTooManyRequests,
	ErrCodeDatabase:   http.StatusInternalServerError,
	ErrCodeInternal:   http.StatusInternalServerError,
}

const (
	ErrCodeValidation = "VALIDATION_ERROR"
	ErrCodeAuth       = "AUTH_ERROR"
	ErrCodeForbidden  = "FORBIDDEN_ERROR"
	ErrCodeNotFound   = "NOT_FOUND_ERROR"
	ErrCodeConflict   = "CONFLICT_ERROR"
	ErrCodeRateLimit  = "RATE_LIMIT_ERROR"
	ErrCodeDatabase   = "DATABASE_ERROR"
	ErrCodeInternal   = "INTERNAL_ERROR"
)

type AppError struct {
	errorType enums.ErrorType
	Code      string
	Message   string
	Status    int
	Cause     error
}

func New(code string, message string) *AppError {
	status, ok := errorCodeMap[code]
	if !ok {
		status = http.StatusInternalServerError
	}
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

func (e *AppError) WithCause(err error) *AppError {
	e.Cause = err
	return e
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Named constructors
func NotFound(resource string) *AppError {
	return New(ErrCodeNotFound, resource+" not found")
}

func InvalidInput() *AppError {
	return New(ErrCodeValidation, "invalid input")
}

func DatabaseFailure(err error) *AppError {
	return New(ErrCodeDatabase, "an internal error occurred").WithCause(err)
}

func Forbidden(message string) *AppError {
	if message == "" {
		message = "access denied"
	}
	return New(ErrCodeForbidden, message)
}

func Unauthorized(message string) *AppError {
	if message == "" {
		message = "unauthorized"
	}
	return New(ErrCodeAuth, message)
}

func Conflict(message string) *AppError {
	return New(ErrCodeConflict, message)
}

func ValidationError(message string) *AppError {
	return New(ErrCodeValidation, message)
}

func InternalError(message string) *AppError {
	return New(ErrCodeInternal, message)
}
