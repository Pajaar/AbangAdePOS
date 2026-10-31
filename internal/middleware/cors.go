package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS mengizinkan origin dari daftar yang dikonfigurasi lewat env
// (CORS_ALLOWED_ORIGINS). AllowCredentials=true karena cookie refresh token
// dikirim cross-origin — origin TIDAK boleh wildcard "*" (browser menolaknya
// kombinasi dengan kredensial), jadi daftar origin harus eksplisit.
func CORS(allowedOrigins []string) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
