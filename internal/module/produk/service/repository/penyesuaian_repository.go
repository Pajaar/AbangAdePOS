package repository

import (
	"context"

	"pos-jamu-backend/internal/module/produk/tck/domain"
	sqlcgen "pos-jamu-backend/internal/repository/sqlc"
)

type IPenyesuaianRepository interface {
	Create(ctx context.Context, p domain.PenyesuaianStok) (uint64, error)
	GetByID(ctx context.Context, id uint64) (domain.PenyesuaianStok, error)
	ListByProduk(ctx context.Context, produkID uint64, limit int32) ([]domain.PenyesuaianStok, error)
}

type penyesuaianRepository struct {
	q *sqlcgen.Queries
}

func NewPenyesuaianRepository(db sqlcgen.DBTX) IPenyesuaianRepository {
	return &penyesuaianRepository{q: sqlcgen.New(db)}
}

func (r *penyesuaianRepository) Create(ctx context.Context, p domain.PenyesuaianStok) (uint64, error) {
	id, err := r.q.CreatePenyesuaianStok(ctx, sqlcgen.CreatePenyesuaianStokParams{
		UserID:      p.UserID,
		ProdukID:    p.ProdukID,
		Jenis:       sqlcgen.PenyesuaianStokJenis(p.Jenis),
		Qty:         int32(p.Qty),
		StokSebelum: int32(p.StokSebelum),
		Hpp:         p.HPP,
		Keterangan:  nullString(p.Keterangan),
	})
	return uint64(id), err
}

func (r *penyesuaianRepository) GetByID(ctx context.Context, id uint64) (domain.PenyesuaianStok, error) {
	row, err := r.q.GetPenyesuaianStokByID(ctx, id)
	if err != nil {
		return domain.PenyesuaianStok{}, err
	}
	return penyesuaianFromRow(sqlcgen.ListPenyesuaianStokByProdukRow(row)), nil
}

func (r *penyesuaianRepository) ListByProduk(ctx context.Context, produkID uint64, limit int32) ([]domain.PenyesuaianStok, error) {
	rows, err := r.q.ListPenyesuaianStokByProduk(ctx, sqlcgen.ListPenyesuaianStokByProdukParams{
		ProdukID: produkID, Limit: limit,
	})
	if err != nil {
		return nil, err
	}
	out := make([]domain.PenyesuaianStok, 0, len(rows))
	for _, row := range rows {
		out = append(out, penyesuaianFromRow(row))
	}
	return out, nil
}

func penyesuaianFromRow(row sqlcgen.ListPenyesuaianStokByProdukRow) domain.PenyesuaianStok {
	qty := int64(row.Qty)
	sebelum := int64(row.StokSebelum)
	return domain.PenyesuaianStok{
		ID:          row.ID,
		Tanggal:     row.Tanggal,
		UserID:      row.UserID,
		UserNama:    row.UserNama,
		ProdukID:    row.ProdukID,
		Jenis:       domain.JenisPenyesuaian(row.Jenis),
		Qty:         qty,
		StokSebelum: sebelum,
		StokSesudah: sebelum + qty,
		HPP:         row.Hpp,
		Keterangan:  strPtr(row.Keterangan),
	}
}
