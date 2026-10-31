package repository

import (
	"context"

	"pos-jamu-backend/internal/module/produk/tck/domain"
	sqlcgen "pos-jamu-backend/internal/repository/sqlc"
)

// CachedLookup = isi barcode_lookup_cache.
type CachedLookup struct {
	Sumber    string
	Ditemukan bool
	Nama      *string
	Merek     *string
}

type IBarcodeRepository interface {
	ListByProduk(ctx context.Context, produkID uint64) ([]domain.Barcode, error)
	GetByID(ctx context.Context, id uint64) (domain.Barcode, error)
	GetByKode(ctx context.Context, kode string) (domain.Barcode, error)
	Create(ctx context.Context, b domain.Barcode) (uint64, error)
	Delete(ctx context.Context, id uint64) error

	GetCache(ctx context.Context, kode string) (CachedLookup, error)
	UpsertCache(ctx context.Context, kode string, c CachedLookup, raw []byte) error
}

type barcodeRepository struct {
	q *sqlcgen.Queries
}

func NewBarcodeRepository(db sqlcgen.DBTX) IBarcodeRepository {
	return &barcodeRepository{q: sqlcgen.New(db)}
}

func (r *barcodeRepository) ListByProduk(ctx context.Context, produkID uint64) ([]domain.Barcode, error) {
	rows, err := r.q.ListBarcodeByProduk(ctx, produkID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.Barcode, 0, len(rows))
	for _, row := range rows {
		out = append(out, barcodeFromRow(row))
	}
	return out, nil
}

func (r *barcodeRepository) GetByID(ctx context.Context, id uint64) (domain.Barcode, error) {
	row, err := r.q.GetBarcodeByID(ctx, id)
	return barcodeFromRow(row), err
}

func (r *barcodeRepository) GetByKode(ctx context.Context, kode string) (domain.Barcode, error) {
	row, err := r.q.GetBarcodeByKode(ctx, kode)
	return barcodeFromRow(row), err
}

func (r *barcodeRepository) Create(ctx context.Context, b domain.Barcode) (uint64, error) {
	id, err := r.q.CreateBarcode(ctx, sqlcgen.CreateBarcodeParams{
		ProdukID:      b.ProdukID,
		Barcode:       b.Kode,
		Jenis:         sqlcgen.ProdukBarcodeJenis(b.Jenis),
		IsiPerKemasan: b.IsiPerKemasan,
		Sumber:        sqlcgen.ProdukBarcodeSumber(b.Sumber),
	})
	return uint64(id), err
}

func (r *barcodeRepository) Delete(ctx context.Context, id uint64) error {
	return r.q.DeleteBarcode(ctx, id)
}

func (r *barcodeRepository) GetCache(ctx context.Context, kode string) (CachedLookup, error) {
	row, err := r.q.GetBarcodeLookupCache(ctx, kode)
	if err != nil {
		return CachedLookup{}, err
	}
	return CachedLookup{
		Sumber: row.Sumber, Ditemukan: row.Ditemukan,
		Nama: strPtr(row.Nama), Merek: strPtr(row.Merek),
	}, nil
}

func (r *barcodeRepository) UpsertCache(ctx context.Context, kode string, c CachedLookup, raw []byte) error {
	return r.q.UpsertBarcodeLookupCache(ctx, sqlcgen.UpsertBarcodeLookupCacheParams{
		Barcode:   kode,
		Sumber:    c.Sumber,
		Ditemukan: c.Ditemukan,
		Nama:      nullString(c.Nama),
		Merek:     nullString(c.Merek),
		RawJson:   raw,
	})
}

func barcodeFromRow(row sqlcgen.ProdukBarcode) domain.Barcode {
	return domain.Barcode{
		ID:            row.ID,
		ProdukID:      row.ProdukID,
		Kode:          row.Barcode,
		Jenis:         string(row.Jenis),
		IsiPerKemasan: row.IsiPerKemasan,
		Sumber:        string(row.Sumber),
		CreatedAt:     row.CreatedAt,
	}
}
