package domain

import "fmt"

// AppError dipetakan ke response JSON oleh httpapi.
type AppError struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Status  int            `json:"-"`
	Details map[string]any `json:"details,omitempty"`
}

func (e *AppError) Error() string { return fmt.Sprintf("%s: %s", e.Code, e.Message) }

func NewError(status int, code, message string) *AppError {
	return &AppError{Code: code, Message: message, Status: status}
}

func Validation(message string, details map[string]any) *AppError {
	return &AppError{Code: "VALIDATION_ERROR", Message: message, Status: 400, Details: details}
}

func NotFound(entity string) *AppError {
	return NewError(404, "NOT_FOUND", entity+" tidak ditemukan")
}

func Conflict(message string) *AppError {
	return NewError(409, "CONFLICT", message)
}

func Forbidden(message string) *AppError {
	return NewError(403, "FORBIDDEN", message)
}

func Unauthorized(message string) *AppError {
	return NewError(401, "UNAUTHENTICATED", message)
}

func Internal(err error) *AppError {
	return &AppError{Code: "INTERNAL", Message: "terjadi kesalahan internal", Status: 500, Details: map[string]any{"cause": err.Error()}}
}
