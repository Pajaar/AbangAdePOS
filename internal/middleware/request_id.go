package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

const requestIDKey = "request_id"

// RequestID generate ID acak per request, disimpan di context dan
// dikirim balik lewat header X-Request-ID supaya bisa dikorelasikan
// dengan log Logger middleware.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := generateRequestID()
		c.Set(requestIDKey, id)
		c.Writer.Header().Set("X-Request-ID", id)
		c.Next()
	}
}

func generateRequestID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(b)
}

func requestIDFromContext(c *gin.Context) string {
	if v, ok := c.Get(requestIDKey); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
