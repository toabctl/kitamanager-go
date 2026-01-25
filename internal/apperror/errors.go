package apperror

import (
	"errors"
	"net/http"
)

// Sentinel errors for domain operations
var (
	ErrNotFound       = errors.New("resource not found")
	ErrBadRequest     = errors.New("bad request")
	ErrConflict       = errors.New("resource conflict")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrInternalServer = errors.New("internal server error")
)

// AppError wraps errors with HTTP context
type AppError struct {
	Err     error
	Message string
	Code    int
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// NotFound creates a not found error
func NotFound(resource string) *AppError {
	return &AppError{Err: ErrNotFound, Message: resource + " not found", Code: http.StatusNotFound}
}

// BadRequest creates a bad request error
func BadRequest(msg string) *AppError {
	return &AppError{Err: ErrBadRequest, Message: msg, Code: http.StatusBadRequest}
}

// Conflict creates a conflict error
func Conflict(msg string) *AppError {
	return &AppError{Err: ErrConflict, Message: msg, Code: http.StatusConflict}
}

// Forbidden creates a forbidden error
func Forbidden(msg string) *AppError {
	return &AppError{Err: ErrForbidden, Message: msg, Code: http.StatusForbidden}
}

// Internal creates an internal server error
func Internal(msg string) *AppError {
	return &AppError{Err: ErrInternalServer, Message: msg, Code: http.StatusInternalServerError}
}

// HTTPStatus returns appropriate status code for error
func HTTPStatus(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	if errors.Is(err, ErrNotFound) {
		return http.StatusNotFound
	}
	if errors.Is(err, ErrBadRequest) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrConflict) {
		return http.StatusConflict
	}
	if errors.Is(err, ErrForbidden) {
		return http.StatusForbidden
	}
	if errors.Is(err, ErrUnauthorized) {
		return http.StatusUnauthorized
	}
	return http.StatusInternalServerError
}
