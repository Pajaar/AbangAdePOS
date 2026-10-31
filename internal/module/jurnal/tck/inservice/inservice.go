// Package inservice berisi kontrak service jurnal.
package inservice

import (
	"context"
	"database/sql"

	"pos-jamu-backend/internal/module/jurnal/tck/domain"
)

type IJurnalInService interface {
	// Post menulis jurnal + detail di dalam transaksi milik pemanggil, supaya
	// jurnal ikut batal bila transaksi bisnisnya batal. Baris bernilai nol
	// diabaikan; bila semua baris nol, tidak ada jurnal yang dibuat (id = 0).
	// Gagal bila total debit tidak sama dengan total kredit.
	Post(ctx context.Context, tx *sql.Tx, e domain.Entri) (uint64, error)
}
