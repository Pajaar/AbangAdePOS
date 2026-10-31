package repository

import (
	"context"
	"database/sql"

	"github.com/shopspring/decimal"

	"pos-jamu-backend/internal/module/produk/tck/domain"
	sqlcgen "pos-jamu-backend/internal/repository/sqlc"
)

// StokHPP = stok dan HPP produk yang sedang dikunci (FOR UPDATE).
type StokHPP struct {
	Stok        int64
	HPPRataRata decimal.Decimal
}

type IProdukRepository interface {
	List(ctx context.Context, f domain.ProdukFilter) ([]domain.Produk, error)
	Count(ctx context.Context, f domain.ProdukFilter) (int64, error)
	GetByID(ctx context.Context, id uint64) (domain.Produk, error)
	// GetStokForUpdate wajib dipanggil dengan *sql.Tx.
	GetStokForUpdate(ctx context.Context, id uint64) (StokHPP, error)
	Create(ctx context.Context, in domain.ProdukInput, status domain.Status) (uint64, error)
	Update(ctx context.Context, id uint64, in domain.ProdukInput, status domain.Status) error
	UpdateStokHPP(ctx context.Context, id uint64, stok int64, hpp decimal.Decimal) error
}

type produkRepository struct {
	q *sqlcgen.Queries
}

func NewProdukRepository(db sqlcgen.DBTX) IProdukRepository {
	return &produkRepository{q: sqlcgen.New(db)}
}

type filterParams struct {
	q          sql.NullString
	kategoriID sql.NullInt64
	status     sqlcgen.NullProdukStatus
	aktif      sql.NullBool
}

func toFilterParams(f domain.ProdukFilter) filterParams {
	var p filterParams
	if f.Q != nil && *f.Q != "" {
		p.q = sql.NullString{String: *f.Q, Valid: true}
	}
	if f.KategoriID != nil {
		p.kategoriID = sql.NullInt64{Int64: int64(*f.KategoriID), Valid: true}
	}
	if f.Status != nil {
		p.status = sqlcgen.NullProdukStatus{ProdukStatus: sqlcgen.ProdukStatus(*f.Status), Valid: true}
	}
	if f.Aktif != nil {
		p.aktif = sql.NullBool{Bool: *f.Aktif, Valid: true}
	}
	return p
}

func (r *produkRepository) List(ctx context.Context, f domain.ProdukFilter) ([]domain.Produk, error) {
	p := toFilterParams(f)
	rows, err := r.q.ListProduk(ctx, sqlcgen.ListProdukParams{
		Q: p.q, KategoriID: p.kategoriID, Status: p.status, Aktif: p.aktif,
		Limit: f.Limit, Offset: f.Offset,
	})
	if err != nil {
		return nil, err
	}
	out := make([]domain.Produk, 0, len(rows))
	for _, row := range rows {
		out = append(out, produkFromRow(sqlcgen.GetProdukByIDRow(row)))
	}
	return out, nil
}

func (r *produkRepository) Count(ctx context.Context, f domain.ProdukFilter) (int64, error) {
	p := toFilterParams(f)
	return r.q.CountProduk(ctx, sqlcgen.CountProdukParams{
		Q: p.q, KategoriID: p.kategoriID, Status: p.status, Aktif: p.aktif,
	})
}

func (r *produkRepository) GetByID(ctx context.Context, id uint64) (domain.Produk, error) {
	row, err := r.q.GetProdukByID(ctx, id)
	if err != nil {
		return domain.Produk{}, err
	}
	return produkFromRow(row), nil
}

func (r *produkRepository) GetStokForUpdate(ctx context.Context, id uint64) (StokHPP, error) {
	row, err := r.q.GetProdukStokForUpdate(ctx, id)
	return StokHPP{Stok: int64(row.Stok), HPPRataRata: row.HppRataRata}, err
}

func (r *produkRepository) Create(ctx context.Context, in domain.ProdukInput, status domain.Status) (uint64, error) {
	id, err := r.q.CreateProduk(ctx, sqlcgen.CreateProdukParams{
		KategoriID:      in.KategoriID,
		Nama:            in.Nama,
		HargaNormal:     in.HargaNormal,
		BisaDitawar:     in.BisaDitawar,
		MarginMinPersen: nullDecimal(in.MarginMinPersen),
		AturanPakai:     nullString(in.AturanPakai),
		Status:          sqlcgen.ProdukStatus(status),
	})
	return uint64(id), err
}

func (r *produkRepository) Update(ctx context.Context, id uint64, in domain.ProdukInput, status domain.Status) error {
	return r.q.UpdateProduk(ctx, sqlcgen.UpdateProdukParams{
		KategoriID:      in.KategoriID,
		Nama:            in.Nama,
		HargaNormal:     in.HargaNormal,
		BisaDitawar:     in.BisaDitawar,
		MarginMinPersen: nullDecimal(in.MarginMinPersen),
		AturanPakai:     nullString(in.AturanPakai),
		Status:          sqlcgen.ProdukStatus(status),
		Aktif:           in.Aktif,
		ID:              id,
	})
}

func (r *produkRepository) UpdateStokHPP(ctx context.Context, id uint64, stok int64, hpp decimal.Decimal) error {
	return r.q.UpdateProdukStokHPP(ctx, sqlcgen.UpdateProdukStokHPPParams{
		Stok: int32(stok), HppRataRata: hpp, ID: id,
	})
}

func produkFromRow(row sqlcgen.GetProdukByIDRow) domain.Produk {
	return domain.Produk{
		ID:              row.ID,
		KategoriID:      row.KategoriID,
		KategoriNama:    row.KategoriNama,
		Nama:            row.Nama,
		HargaNormal:     row.HargaNormal,
		BisaDitawar:     row.BisaDitawar,
		MarginMinPersen: decPtr(row.MarginMinPersen),
		Stok:            int64(row.Stok),
		HPPRataRata:     row.HppRataRata,
		AturanPakai:     strPtr(row.AturanPakai),
		Status:          domain.Status(row.Status),
		Aktif:           row.Aktif,
		CreatedAt:       row.CreatedAt,
	}
}
