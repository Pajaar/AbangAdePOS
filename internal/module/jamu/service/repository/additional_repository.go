package repository

import (
	"context"
	"database/sql"

	"github.com/shopspring/decimal"

	"pos-jamu-backend/internal/module/jamu/tck/domain"
	sqlcgen "pos-jamu-backend/internal/repository/sqlc"
)

type IAdditionalRepository interface {
	List(ctx context.Context, aktif *bool) ([]domain.AdditionalItem, error)
	GetByID(ctx context.Context, id uint64) (domain.AdditionalItem, error)
	Create(ctx context.Context, in domain.AdditionalInput) (uint64, error)
	Update(ctx context.Context, id uint64, in domain.AdditionalInput) error
}

type additionalRepository struct {
	q *sqlcgen.Queries
}

func NewAdditionalRepository(db sqlcgen.DBTX) IAdditionalRepository {
	return &additionalRepository{q: sqlcgen.New(db)}
}

func (r *additionalRepository) List(ctx context.Context, aktif *bool) ([]domain.AdditionalItem, error) {
	rows, err := r.q.ListAdditionalItems(ctx, sqlcgen.ListAdditionalItemsParams{Aktif: nullBool(aktif)})
	if err != nil {
		return nil, err
	}
	out := make([]domain.AdditionalItem, 0, len(rows))
	for _, row := range rows {
		out = append(out, additionalFromRow(sqlcgen.GetAdditionalItemByIDRow(row)))
	}
	return out, nil
}

func (r *additionalRepository) GetByID(ctx context.Context, id uint64) (domain.AdditionalItem, error) {
	row, err := r.q.GetAdditionalItemByID(ctx, id)
	if err != nil {
		return domain.AdditionalItem{}, err
	}
	return additionalFromRow(row), nil
}

func (r *additionalRepository) Create(ctx context.Context, in domain.AdditionalInput) (uint64, error) {
	id, err := r.q.CreateAdditionalItem(ctx, sqlcgen.CreateAdditionalItemParams{
		Nama: in.Nama, HargaJual: in.HargaJual, ProdukID: nullID(in.ProdukID),
	})
	return uint64(id), err
}

func (r *additionalRepository) Update(ctx context.Context, id uint64, in domain.AdditionalInput) error {
	return r.q.UpdateAdditionalItem(ctx, sqlcgen.UpdateAdditionalItemParams{
		Nama: in.Nama, HargaJual: in.HargaJual, ProdukID: nullID(in.ProdukID), Aktif: in.Aktif, ID: id,
	})
}

func additionalFromRow(row sqlcgen.GetAdditionalItemByIDRow) domain.AdditionalItem {
	a := domain.AdditionalItem{
		ID:        row.ID,
		Nama:      row.Nama,
		HargaJual: row.HargaJual,
		Aktif:     row.Aktif,
		ProdukHPP: decimal.Zero,
	}
	if row.ProdukID.Valid {
		id := uint64(row.ProdukID.Int64)
		a.ProdukID = &id
	}
	if row.ProdukNama.Valid {
		nama := row.ProdukNama.String
		a.ProdukNama = &nama
	}
	if row.ProdukHpp.Valid {
		a.ProdukHPP = row.ProdukHpp.Decimal
	}
	return a
}

func nullID(p *uint64) sql.NullInt64 {
	if p == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(*p), Valid: true}
}
