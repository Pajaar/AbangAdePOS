package middleware

import (
	"fmt"
	"log/slog"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"pos-jamu-backend/internal/util"
	"pos-jamu-backend/internal/util/apperror"
)

// Recovery menangkap panic di handler/middleware berikutnya supaya proses
// tidak crash, membalas 500 JSON seragam, dan mencatat stack trace.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.ErrorContext(c.Request.Context(), "panic_recovered",
					"request_id", requestIDFromContext(c),
					"panic", fmt.Sprintf("%v", rec),
					"stack", string(debug.Stack()),
				)
				util.Error(c, apperror.Internal(fmt.Errorf("panic: %v", rec)))
				c.Abort()
			}
		}()
		c.Next()
	}
}
