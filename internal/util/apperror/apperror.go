// Package apperror mendefinisikan error aplikasi yang membawa status HTTP,
// supaya handler/service tidak perlu memutuskan status code di tempat lain.
package apperror

import "net/http"

type AppError struct {
	Status  int
	Code    string
	Message string
	Details any // opsional, ikut dikirim ke klien (misalnya daftar error validasi)
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(status int, code, message string) *AppError {
	return &AppError{Status: status, Code: code, Message: message}
}

func BadRequest(message string) *AppError {
	return New(http.StatusBadRequest, "bad_request", message)
}

// Validation dipakai processor.Validate untuk pelanggaran aturan bisnis.
func Validation(details any) *AppError {
	return &AppError{
		Status:  http.StatusBadRequest,
		Code:    "validation_error",
		Message: "Data tidak valid",
		Details: details,
	}
}

func Unauthorized(message string) *AppError {
	return New(http.StatusUnauthorized, "unauthorized", message)
}

func Forbidden(message string) *AppError {
	return New(http.StatusForbidden, "forbidden", message)
}

func NotFound(message string) *AppError {
	return New(http.StatusNotFound, "not_found", message)
}

func Conflict(message string) *AppError {
	return New(http.StatusConflict, "conflict", message)
}

func TooManyRequests(message string) *AppError {
	return New(http.StatusTooManyRequests, "too_many_requests", message)
}

// Internal membungkus error asli yang detailnya hanya dicatat ke log,
// tidak dikirim ke klien.
func Internal(err error) *AppError {
	return &AppError{
		Status:  http.StatusInternalServerError,
		Code:    "internal_error",
		Message: "Terjadi kesalahan pada server",
		Err:     err,
	}
}
