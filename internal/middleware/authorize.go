package middleware

import (
	"github.com/gin-gonic/gin"

	"pos-jamu-backend/internal/authorize"
	"pos-jamu-backend/internal/util"
	"pos-jamu-backend/internal/util/apperror"
)

// Authorize membandingkan method + c.FullPath() (pola route, bukan URL
// asli) terhadap izin role. Headadmin TIDAK di-bypass di sini — izinnya
// datang murni dari data seed role_menus (lihat migrations/000002).
func Authorize(store *authorize.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, ok := c.Get("role")
		roleKode, _ := roleVal.(string)
		if !ok || roleKode == "" {
			util.Error(c, apperror.Unauthorized("Belum login"))
			c.Abort()
			return
		}

		if !store.IsAllowed(roleKode, c.Request.Method, c.FullPath()) {
			util.Error(c, apperror.Forbidden("Anda tidak punya akses ke fitur ini"))
			c.Abort()
			return
		}

		c.Next()
	}
}
