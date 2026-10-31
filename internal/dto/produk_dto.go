package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

// Uang dikirim dan diterima sebagai string desimal, misalnya "15000" atau
// "12500.50" (decimal.Decimal), supaya tidak ada pembulatan float.
// Field bertanda json:"-" diisi controller dari path/context, bukan dari body.

// ---------- Kategori ----------

type KategoriListRequest struct {
	Aktif *bool `form:"aktif"`
}

type KategoriRequest struct {
	ID    uint64 `json:"-"`
	Nama  string `json:"nama" binding:"required,max=100"`
	Aktif *bool  `json:"aktif"` // hanya untuk ubah; kosong = true
}

type KategoriResponse struct {
	ID    uint64 `json:"id"`
	Nama  string `json:"nama"`
	Aktif bool   `json:"aktif"`
}

// ---------- Produk ----------

type ProdukListRequest struct {
	Q          string  `form:"q"` // sebagian nama atau barcode persis
	KategoriID *uint64 `form:"kategori_id"`
	Status     string  `form:"status" binding:"omitempty,oneof=aktif draft"`
	Aktif      *bool   `form:"aktif"`
	Limit      int32   `form:"limit" binding:"omitempty,min=1,max=200"`
	Offset     int32   `form:"offset" binding:"omitempty,min=0"`
	RoleKode   string  `form:"-"`
}

type BarcodeRequest struct {
	Kode          string `json:"kode" binding:"required"`
	IsiPerKemasan int32  `json:"isi_per_kemasan" binding:"omitempty,min=1"` // kosong = 1
	Sumber        string `json:"sumber" binding:"omitempty,oneof=manual openfoodfacts upcitemdb"`
}

type ProdukRequest struct {
	ID       uint64 `json:"-"`
	RoleKode string `json:"-"`

	KategoriID      uint64           `json:"kategori_id" binding:"required"`
	Nama            string           `json:"nama" binding:"required,max=150"`
	HargaNormal     decimal.Decimal  `json:"harga_normal"` // 0 = produk disimpan sebagai draft
	BisaDitawar     bool             `json:"bisa_ditawar"`
	MarginMinPersen *decimal.Decimal `json:"margin_min_persen"` // wajib bila bisa_ditawar
	AturanPakai     *string          `json:"aturan_pakai" binding:"omitempty,max=255"`
	Aktif           *bool            `json:"aktif"`                             // hanya untuk ubah; kosong = true
	Barcodes        []BarcodeRequest `json:"barcodes" binding:"omitempty,dive"` // hanya untuk tambah
}

type BarcodeResponse struct {
	ID            uint64    `json:"id"`
	ProdukID      uint64    `json:"produk_id"`
	Kode          string    `json:"kode"`
	Jenis         string    `json:"jenis"`
	IsiPerKemasan int32     `json:"isi_per_kemasan"`
	Sumber        string    `json:"sumber"`
	CreatedAt     time.Time `json:"created_at"`
}

// ProdukResponse: field HPP (hpp_rata_rata, margin, harga nett/tengah,
// harga_perlu_dinaikkan) hanya dikirim untuk role dengan izin lihat_hpp.
// Untuk role lain field itu tidak ada sama sekali di JSON.
type ProdukResponse struct {
	ID           uint64          `json:"id"`
	KategoriID   uint64          `json:"kategori_id"`
	KategoriNama string          `json:"kategori_nama"`
	Nama         string          `json:"nama"`
	HargaNormal  decimal.Decimal `json:"harga_normal"`
	BisaDitawar  bool            `json:"bisa_ditawar"`
	Stok         int64           `json:"stok"`
	AturanPakai  *string         `json:"aturan_pakai"`
	Status       string          `json:"status"`
	Aktif        bool            `json:"aktif"`
	CreatedAt    time.Time       `json:"created_at"`

	HPPRataRata         *decimal.Decimal `json:"hpp_rata_rata,omitempty"`
	MarginMinPersen     *decimal.Decimal `json:"margin_min_persen,omitempty"`
	HargaNett           *decimal.Decimal `json:"harga_nett,omitempty"`
	HargaTengah         *decimal.Decimal `json:"harga_tengah,omitempty"`
	HargaPerluDinaikkan *bool            `json:"harga_perlu_dinaikkan,omitempty"`

	Barcodes []BarcodeResponse `json:"barcodes,omitempty"`
}

// ---------- Barcode ----------

type BarcodeAddRequest struct {
	ProdukID uint64 `json:"-"`
	BarcodeRequest
}

type BarcodeKasirRequest struct {
	Kode     string
	RoleKode string
}

type BarcodeKasirResponse struct {
	Produk  ProdukResponse  `json:"produk"`
	Barcode BarcodeResponse `json:"barcode"`
}

type BarcodeLookupRequest struct {
	Kode     string `json:"kode" binding:"required"`
	RoleKode string `json:"-"`
}

type BarcodeLookupResponse struct {
	Kode          string          `json:"kode"`
	Jenis         string          `json:"jenis"`
	Ditemukan     bool            `json:"ditemukan"`
	Sumber        string          `json:"sumber,omitempty"`
	DariCache     bool            `json:"dari_cache"`
	Nama          *string         `json:"nama"`  // saran, user wajib cek
	Merek         *string         `json:"merek"` // saran
	Produk        *ProdukResponse `json:"produk,omitempty"`
	SebagianGagal bool            `json:"sebagian_gagal"`
}

// ---------- Penyesuaian stok ----------

// PenyesuaianRequest:
//   - stok_awal: qty (> 0) dan harga_modal per unit
//   - opname:    stok_fisik (hasil hitung fisik)
//   - rusak:     qty = jumlah unit rusak (> 0)
type PenyesuaianRequest struct {
	ProdukID uint64 `json:"-"`
	UserID   uint64 `json:"-"`
	RoleKode string `json:"-"`

	Jenis      string          `json:"jenis" binding:"required,oneof=stok_awal opname rusak"`
	Qty        int64           `json:"qty"`
	StokFisik  *int64          `json:"stok_fisik"`
	HargaModal decimal.Decimal `json:"harga_modal"`
	Keterangan string          `json:"keterangan" binding:"max=255"`
}

type PenyesuaianListRequest struct {
	ProdukID uint64 `form:"-"`
	RoleKode string `form:"-"`
	Limit    int32  `form:"limit" binding:"omitempty,min=1,max=200"`
}

// PenyesuaianResponse: hpp (nilai per unit) hanya untuk role lihat_hpp.
type PenyesuaianResponse struct {
	ID          uint64           `json:"id"`
	Tanggal     time.Time        `json:"tanggal"`
	UserNama    string           `json:"user_nama"`
	ProdukID    uint64           `json:"produk_id"`
	Jenis       string           `json:"jenis"`
	Qty         int64            `json:"qty"`
	StokSebelum int64            `json:"stok_sebelum"`
	StokSesudah int64            `json:"stok_sesudah"`
	HPP         *decimal.Decimal `json:"hpp,omitempty"`
	Keterangan  *string          `json:"keterangan"`
}
