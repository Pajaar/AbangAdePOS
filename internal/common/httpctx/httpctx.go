// Package httpctx membaca data user yang diisi middleware Auth ke context Gin.
package httpctx

import "github.com/gin-gonic/gin"

// Kunci context yang diisi middleware.Auth.
const (
	KeyUserID = "user_id"
	KeyRole   = "role"
)

func UserID(c *gin.Context) (uint64, bool) {
	v, ok := c.Get(KeyUserID)
	if !ok {
		return 0, false
	}
	id, ok := v.(uint64)
	return id, ok && id != 0
}

func RoleKode(c *gin.Context) (string, bool) {
	v, ok := c.Get(KeyRole)
	if !ok {
		return "", false
	}
	role, ok := v.(string)
	return role, ok && role != ""
}
