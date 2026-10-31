package database

import (
	"context"
	"database/sql"
	"fmt"
)

// WithTx membuka transaksi, menjalankan fn, lalu commit kalau sukses atau
// rollback kalau fn error / panic. Semua langkah satu operasi bisnis
// (nota, detail, potong stok, piutang, jurnal, rotasi refresh token, dll)
// harus dibungkus di sini, dipanggil dari service — bukan dari middleware.
func WithTx(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("%w (rollback error: %v)", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
