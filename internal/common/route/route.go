// Package route menyediakan grup route yang dibagikan router ke setiap
// modul, supaya urutan middleware tetap diatur di satu tempat (router).
package route

import "github.com/gin-gonic/gin"

type Groups struct {
	Public        *gin.RouterGroup // tanpa Auth
	Authenticated *gin.RouterGroup // Auth saja (semua user login)
	Protected     *gin.RouterGroup // Auth + Authorize
}

// Registrar diimplementasikan config setiap modul (setara ControllerConfig).
type Registrar interface {
	RegisterRoutes(g Groups)
}
