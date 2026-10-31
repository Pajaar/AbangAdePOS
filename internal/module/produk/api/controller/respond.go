// Package controller berisi handler HTTP modul produk. Controller hanya
// bind input, mengisi data dari path/context, memanggil executor, lalu
// mengirim respons seragam.
package controller

import (
	"github.com/gin-gonic/gin"

	"pos-jamu-backend/internal/util"
)

func respond[T any](c *gin.Context, status int, res T, err error) {
	if err != nil {
		util.Error(c, err)
		return
	}
	util.Success(c, status, res)
}
