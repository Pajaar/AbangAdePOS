package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"pos-jamu-backend/internal/auth"
	"pos-jamu-backend/internal/util"
	"pos-jamu-backend/internal/util/apperror"
)

// Auth memverifikasi Bearer access token dan mengisi context dengan
// user_id dan role untuk dipakai handler / Authorize berikutnya.
func Auth(jwtSecret string) gin.HandlerFunc {
	const prefix = "Bearer "

	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, prefix) || len(header) <= len(prefix) {
			util.Error(c, apperror.Unauthorized("Belum login"))
			c.Abort()
			return
		}

		claims, err := auth.ParseAccessToken(header[len(prefix):], jwtSecret)
		if err != nil {
			util.Error(c, apperror.Unauthorized("Token tidak valid atau sudah habis"))
			c.Abort()
			return
		}

		userID, err := claims.UserID()
		if err != nil {
			util.Error(c, apperror.Unauthorized("Token tidak valid"))
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
