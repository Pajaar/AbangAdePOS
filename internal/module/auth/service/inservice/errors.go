// Package inservice berisi implementasi kontrak tck/inservice (setara
// service/internal di primecash). Transaksi database dibuka di sini lewat
// database.WithTx, sesuai aturan wajib nomor 2 di CLAUDE.md.
package inservice

import (
	"errors"

	"pos-jamu-backend/internal/util/apperror"
)

var (
	errLoginFailed    = apperror.Unauthorized("Username atau password salah")
	errSessionInvalid = apperror.Unauthorized("Sesi tidak valid, silakan login ulang")
)

// asAppError meneruskan *apperror.AppError apa adanya dan membungkus error
// lain (database, dll) sebagai Internal.
func asAppError(err error) error {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return apperror.Internal(err)
}
