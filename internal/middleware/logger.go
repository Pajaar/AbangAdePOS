package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger mencatat setiap request sebagai satu baris JSON terstruktur
// (method, path, status, durasi, request_id, client_ip).
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		slog.InfoContext(c.Request.Context(), "request",
			"method", c.Request.Method,
			"path", path,
			"status", c.Writer.Status(),
			"duration_ms", duration.Milliseconds(),
			"request_id", requestIDFromContext(c),
			"client_ip", c.ClientIP(),
		)
	}
}
