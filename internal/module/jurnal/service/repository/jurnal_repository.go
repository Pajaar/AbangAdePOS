package repository

import (
	"context"
	"database/sql"

	"github.com/shopspring/decimal"

	sqlcgen "pos-jamu-backend/internal/repository/sqlc"
)

type IJurnalRepository interface {
	// GetAkunIDByKode mengembalikan sql.ErrNoRows bila akun tidak ada.
	GetAkunIDByKode(ctx context.Context, kode string) (uint64, error)
	CreateJurnal(ctx context.Context, userID uint64, keterangan, sumberTipe string, sumberID uint64) (uint64, error)
	CreateDetail(ctx context.Context, jurnalID, akunID uint64, debit, kredit decimal.Decimal) error
}

type jurnalRepository struct {
	q *sqlcgen.Queries
}

func NewJurnalRepository(db sqlcgen.DBTX) IJurnalRepository {
	return &jurnalRepository{q: sqlcgen.New(db)}
}

func (r *jurnalRepository) GetAkunIDByKode(ctx context.Context, kode string) (uint64, error) {
	return r.q.GetAkunIDByKode(ctx, kode)
}

func (r *jurnalRepository) CreateJurnal(ctx context.Context, userID uint64, keterangan, sumberTipe string, sumberID uint64) (uint64, error) {
	id, err := r.q.CreateJurnal(ctx, sqlcgen.CreateJurnalParams{
		UserID:     userID,
		Keterangan: sql.NullString{String: keterangan, Valid: keterangan != ""},
		SumberTipe: sql.NullString{String: sumberTipe, Valid: sumberTipe != ""},
		SumberID:   sql.NullInt64{Int64: int64(sumberID), Valid: sumberID != 0},
	})
	return uint64(id), err
}

func (r *jurnalRepository) CreateDetail(ctx context.Context, jurnalID, akunID uint64, debit, kredit decimal.Decimal) error {
	return r.q.CreateJurnalDetail(ctx, sqlcgen.CreateJurnalDetailParams{
		JurnalID: jurnalID, AkunID: akunID, Debit: debit, Kredit: kredit,
	})
}
