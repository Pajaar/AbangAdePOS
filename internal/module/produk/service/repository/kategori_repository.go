package repository

import (
	"context"
	"database/sql"

	"pos-jamu-backend/internal/module/produk/tck/domain"
	sqlcgen "pos-jamu-backend/internal/repository/sqlc"
)

type IKategoriRepository interface {
	List(ctx context.Context, aktif *bool) ([]domain.Kategori, error)
	GetByID(ctx context.Context, id uint64) (domain.Kategori, error)
	GetByNama(ctx context.Context, nama string) (domain.Kategori, error)
	Create(ctx context.Context, nama string) (uint64, error)
	Update(ctx context.Context, id uint64, nama string, aktif bool) error
}

type kategoriRepository struct {
	q *sqlcgen.Queries
}

func NewKategoriRepository(db sqlcgen.DBTX) IKategoriRepository {
	return &kategoriRepository{q: sqlcgen.New(db)}
}

func (r *kategoriRepository) List(ctx context.Context, aktif *bool) ([]domain.Kategori, error) {
	param := sql.NullBool{}
	if aktif != nil {
		param = sql.NullBool{Bool: *aktif, Valid: true}
	}
	rows, err := r.q.ListKategori(ctx, sqlcgen.ListKategoriParams{Aktif: param})
	if err != nil {
		return nil, err
	}
	out := make([]domain.Kategori, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.Kategori{ID: row.ID, Nama: row.Nama, Aktif: row.Aktif})
	}
	return out, nil
}

func (r *kategoriRepository) GetByID(ctx context.Context, id uint64) (domain.Kategori, error) {
	row, err := r.q.GetKategoriByID(ctx, id)
	return domain.Kategori{ID: row.ID, Nama: row.Nama, Aktif: row.Aktif}, err
}

func (r *kategoriRepository) GetByNama(ctx context.Context, nama string) (domain.Kategori, error) {
	row, err := r.q.GetKategoriByNama(ctx, nama)
	return domain.Kategori{ID: row.ID, Nama: row.Nama, Aktif: row.Aktif}, err
}

func (r *kategoriRepository) Create(ctx context.Context, nama string) (uint64, error) {
	id, err := r.q.CreateKategori(ctx, nama)
	return uint64(id), err
}

func (r *kategoriRepository) Update(ctx context.Context, id uint64, nama string, aktif bool) error {
	return r.q.UpdateKategori(ctx, sqlcgen.UpdateKategoriParams{Nama: nama, Aktif: aktif, ID: id})
}
