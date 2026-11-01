// Package repository membungkus query sqlc modul menu jamu dan memetakan
// hasilnya ke tipe domain.
package repository

import (
	"context"
	"database/sql"

	"pos-jamu-backend/internal/module/jamu/tck/domain"
	sqlcgen "pos-jamu-backend/internal/repository/sqlc"
)

type IJamuRepository interface {
	List(ctx context.Context, aktif *bool) ([]domain.Jamu, error)
	GetByID(ctx context.Context, id uint64) (domain.Jamu, error)
	Create(ctx context.Context, in domain.JamuInput) (uint64, error)
	Update(ctx context.Context, id uint64, in domain.JamuInput) error
}

type jamuRepository struct {
	q *sqlcgen.Queries
}

func NewJamuRepository(db sqlcgen.DBTX) IJamuRepository {
	return &jamuRepository{q: sqlcgen.New(db)}
}

func (r *jamuRepository) List(ctx context.Context, aktif *bool) ([]domain.Jamu, error) {
	rows, err := r.q.ListJamu(ctx, sqlcgen.ListJamuParams{Aktif: nullBool(aktif)})
	if err != nil {
		return nil, err
	}
	out := make([]domain.Jamu, 0, len(rows))
	for _, row := range rows {
		out = append(out, jamuFromRow(sqlcgen.GetJamuByIDRow(row)))
	}
	return out, nil
}

func (r *jamuRepository) GetByID(ctx context.Context, id uint64) (domain.Jamu, error) {
	row, err := r.q.GetJamuByID(ctx, id)
	if err != nil {
		return domain.Jamu{}, err
	}
	return jamuFromRow(row), nil
}

func (r *jamuRepository) Create(ctx context.Context, in domain.JamuInput) (uint64, error) {
	id, err := r.q.CreateJamu(ctx, sqlcgen.CreateJamuParams{
		Nama: in.Nama, HargaJual: in.HargaJual, ProdukID: in.ProdukID,
	})
	return uint64(id), err
}

func (r *jamuRepository) Update(ctx context.Context, id uint64, in domain.JamuInput) error {
	return r.q.UpdateJamu(ctx, sqlcgen.UpdateJamuParams{
		Nama: in.Nama, HargaJual: in.HargaJual, ProdukID: in.ProdukID, Aktif: in.Aktif, ID: id,
	})
}

func jamuFromRow(row sqlcgen.GetJamuByIDRow) domain.Jamu {
	return domain.Jamu{
		ID:         row.ID,
		Nama:       row.Nama,
		HargaJual:  row.HargaJual,
		Aktif:      row.Aktif,
		CreatedAt:  row.CreatedAt,
		ProdukID:   row.ProdukID,
		ProdukNama: row.ProdukNama,
		ProdukStok: int64(row.ProdukStok),
		ProdukHPP:  row.ProdukHpp,
	}
}

func nullBool(p *bool) sql.NullBool {
	if p == nil {
		return sql.NullBool{}
	}
	return sql.NullBool{Bool: *p, Valid: true}
}
