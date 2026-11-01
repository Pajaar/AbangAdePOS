package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

// ---------- Jamu seduh ----------

type JamuListRequest struct {
	Aktif    *bool  `form:"aktif"`
	RoleKode string `form:"-"`
}

// JamuRequest: satu gelas = tepat satu sachet dari produk_id.
type JamuRequest struct {
	ID       uint64 `json:"-"`
	RoleKode string `json:"-"`

	Nama      string          `json:"nama" binding:"required,max=100"`
	ProdukID  uint64          `json:"produk_id" binding:"required"` // produk sachet yang diseduh
	HargaJual decimal.Decimal `json:"harga_jual"`                   // harga seduh per gelas
	Aktif     *bool           `json:"aktif"`                        // hanya untuk ubah; kosong = tidak berubah
}

// JamuResponse: hpp (per gelas = HPP sachet) dan margin (harga seduh - hpp)
// hanya dikirim untuk role dengan izin lihat_hpp.
type JamuResponse struct {
	ID         uint64          `json:"id"`
	Nama       string          `json:"nama"`
	HargaJual  decimal.Decimal `json:"harga_jual"`
	Aktif      bool            `json:"aktif"`
	CreatedAt  time.Time       `json:"created_at"`
	ProdukID   uint64          `json:"produk_id"`
	ProdukNama string          `json:"produk_nama"`
	StokSachet int64           `json:"stok_sachet"`

	HPP    *decimal.Decimal `json:"hpp,omitempty"`
	Margin *decimal.Decimal `json:"margin,omitempty"`
}

// ---------- Item tambahan ----------

type AdditionalListRequest struct {
	Aktif    *bool  `form:"aktif"`
	RoleKode string `form:"-"`
}

// AdditionalRequest: tambahan harga tetap (telur, madu). produk_id diisi
// hanya bila tambahan itu memotong stok produk tertentu; kosong = tidak
// dilacak (HPP 0). Sachet tambahan dari produk lain tidak perlu didaftarkan
// di sini; kasir memilih produknya langsung.
type AdditionalRequest struct {
	ID       uint64 `json:"-"`
	RoleKode string `json:"-"`

	Nama      string          `json:"nama" binding:"required,max=100"`
	HargaJual decimal.Decimal `json:"harga_jual"`
	ProdukID  *uint64         `json:"produk_id"`
	Aktif     *bool           `json:"aktif"` // hanya untuk ubah; kosong = tidak berubah
}

// AdditionalResponse: hpp hanya untuk role lihat_hpp.
type AdditionalResponse struct {
	ID          uint64          `json:"id"`
	Nama        string          `json:"nama"`
	HargaJual   decimal.Decimal `json:"harga_jual"`
	ProdukID    *uint64         `json:"produk_id"`
	ProdukNama  *string         `json:"produk_nama"`
	DilacakStok bool            `json:"dilacak_stok"`
	Aktif       bool            `json:"aktif"`

	HPP *decimal.Decimal `json:"hpp,omitempty"`
}
