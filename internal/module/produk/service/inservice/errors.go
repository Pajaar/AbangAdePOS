// Package inservice berisi implementasi kontrak tck/inservice modul produk.
// Transaksi dibuka di sini lewat database.WithTx.
package inservice

import (
	"errors"

	"github.com/go-sql-driver/mysql"

	"pos-jamu-backend/internal/util/apperror"
)

const (
	mysqlErrDuplicate  = 1062
	mysqlErrForeignKey = 1452
)

func mysqlErrNo(err error) uint16 {
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		return me.Number
	}
	return 0
}

// asAppError meneruskan *apperror.AppError apa adanya dan membungkus error
// lain sebagai Internal.
func asAppError(err error) error {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return apperror.Internal(err)
}
