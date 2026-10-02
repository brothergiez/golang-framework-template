package errors

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrConflict      = errors.New("conflict")
	ErrValidation    = errors.New("validation error")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrInternal      = errors.New("internal error")
	ErrBadRequest    = errors.New("bad request")
	ErrUnavailable   = errors.New("service unavailable")
)

type AppError struct {
	Code       string
	Message    string
	HTTPStatus int
	Err        error
	Details    any
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code, message string, httpStatus int, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Err:        err,
	}
}

func NotFound(message string) *AppError {
	return New("NOT_FOUND", message, http.StatusNotFound, ErrNotFound)
}

func Conflict(message string) *AppError {
	return New("CONFLICT", message, http.StatusConflict, ErrConflict)
}

func Validation(message string, details any) *AppError {
	return &AppError{
		Code:       "VALIDATION_ERROR",
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
		Err:        ErrValidation,
		Details:    details,
	}
}

func BadRequest(message string) *AppError {
	return New("BAD_REQUEST", message, http.StatusBadRequest, ErrBadRequest)
}

func Internal(message string, err error) *AppError {
	return New("INTERNAL_ERROR", message, http.StatusInternalServerError, err)
}

func Unauthorized(message string) *AppError {
	return New("UNAUTHORIZED", message, http.StatusUnauthorized, ErrUnauthorized)
}

func Forbidden(message string) *AppError {
	return New("FORBIDDEN", message, http.StatusForbidden, ErrForbidden)
}

func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
