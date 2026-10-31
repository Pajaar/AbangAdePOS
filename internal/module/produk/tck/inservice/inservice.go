// Package inservice berisi kontrak service modul produk.
package inservice

import (
	"context"

	"pos-jamu-backend/internal/module/produk/tck/domain"
)

type IKategoriInService interface {
	List(ctx context.Context, aktif *bool) ([]domain.Kategori, error)
	// GetByID mengembalikan apperror NotFound bila tidak ada.
	GetByID(ctx context.Context, id uint64) (domain.Kategori, error)
	// FindByNama: found=false bila tidak ada (bukan error).
	FindByNama(ctx context.Context, nama string) (k domain.Kategori, found bool, err error)
	Create(ctx context.Context, nama string) (uint64, error)
	Update(ctx context.Context, id uint64, nama string, aktif bool) error
}

type IProdukInService interface {
	List(ctx context.Context, f domain.ProdukFilter) (items []domain.Produk, total int64, err error)
	// GetByID termasuk daftar barcode. NotFound bila tidak ada.
	GetByID(ctx context.Context, id uint64) (domain.Produk, error)
	// Create membuat produk beserta barcode-nya dalam satu transaksi.
	Create(ctx context.Context, in domain.ProdukInput, barcodes []domain.BarcodeInput) (uint64, error)
	Update(ctx context.Context, id uint64, in domain.ProdukInput) error
}

type IBarcodeInService interface {
	Add(ctx context.Context, produkID uint64, in domain.BarcodeInput) (domain.Barcode, error)
	Delete(ctx context.Context, id uint64) error
	// FindByKode mencari barcode (sudah dinormalisasi). found=false bila tidak ada.
	FindByKode(ctx context.Context, kode string) (b domain.Barcode, found bool, err error)
	// FindForKasir: hanya database lokal. NotFound (404) bila tidak ada atau
	// produk nonaktif, Conflict (409) bila produk masih draft.
	FindForKasir(ctx context.Context, raw string) (domain.Produk, domain.Barcode, error)
	// Lookup: lokal -> cache -> API publik (lihat CLAUDE.md bagian Barcode).
	Lookup(ctx context.Context, raw string) (domain.LookupResult, error)
}

type IPenyesuaianStokInService interface {
	// Create mengubah stok (dan HPP untuk stok awal), mencatat penyesuaian,
	// dan membuat jurnal, semuanya dalam satu transaksi.
	Create(ctx context.Context, in domain.PenyesuaianInput) (domain.PenyesuaianStok, error)
	ListByProduk(ctx context.Context, produkID uint64, limit int32) ([]domain.PenyesuaianStok, error)
}

// IBarcodeSource = satu API publik barcode (Open Food Facts, UPCitemdb).
// Implementasinya di internal/integration.
type IBarcodeSource interface {
	Name() string
	// Lookup: Found=false dengan err=nil berarti API menjawab "tidak ada".
	// err != nil berarti error/timeout (jangan di-cache).
	Lookup(ctx context.Context, kode string) (domain.ExternalProduct, error)
}
