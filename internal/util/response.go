// Package util menyediakan helper lintas layer: respons JSON seragam dan
// pemetaan error aplikasi ke status HTTP.
package util

import (
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"

	"pos-jamu-backend/internal/util/apperror"
)

type successEnvelope struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

type errorEnvelope struct {
	Success bool      `json:"success"`
	Error   errorBody `json:"error"`
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func Success(c *gin.Context, status int, data any) {
	c.JSON(status, successEnvelope{Success: true, Data: data})
}

// Error memetakan error ke respons JSON seragam. Error yang bukan *apperror.AppError
// dianggap tak terduga: dicatat sebagai Internal, detail aslinya TIDAK dikirim ke klien.
func Error(c *gin.Context, err error) {
	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		appErr = apperror.Internal(err)
	}

	if appErr.Status >= 500 {
		slog.ErrorContext(c.Request.Context(), "request_error",
			"request_id", requestID(c),
			"error", appErr.Error(),
		)
	}

	c.JSON(appErr.Status, errorEnvelope{
		Success: false,
		Error: errorBody{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
		},
	})
}

func requestID(c *gin.Context) string {
	if v, ok := c.Get("request_id"); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
