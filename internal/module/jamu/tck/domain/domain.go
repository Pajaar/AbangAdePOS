// Package domain berisi tipe domain modul menu jamu: jamu seduh dan item
// tambahan.
package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// Jamu = menu jamu seduh. Satu gelas = tepat satu sachet dari ProdukID.
// Harga seduh (HargaJual) berbeda dari harga sachet utuh di produk.
type Jamu struct {
	ID         uint64
	Nama       string
	HargaJual  decimal.Decimal
	Aktif      bool
	CreatedAt  time.Time
	ProdukID   uint64
	ProdukNama string
	ProdukStok int64           // stok sachet, untuk peringatan
	ProdukHPP  decimal.Decimal // HPP per gelas = HPP sachet
}

type JamuInput struct {
	Nama      string
	HargaJual decimal.Decimal
	ProdukID  uint64
	Aktif     bool // hanya dipakai saat ubah
}

// AdditionalItem = tambahan harga tetap untuk jamu seduh (telur, madu, ...).
// ProdukID kosong = tidak dilacak stoknya, HPP 0. Sachet tambahan dari produk
// lain TIDAK lewat sini; di nota langsung menunjuk produk (harga = harga sachet).
type AdditionalItem struct {
	ID         uint64
	Nama       string
	HargaJual  decimal.Decimal
	ProdukID   *uint64
	ProdukNama *string
	ProdukHPP  decimal.Decimal // 0 bila tidak terhubung produk
	Aktif      bool
}

type AdditionalInput struct {
	Nama      string
	HargaJual decimal.Decimal
	ProdukID  *uint64
	Aktif     bool // hanya dipakai saat ubah
}
