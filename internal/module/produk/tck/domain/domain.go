// Package domain berisi tipe domain modul produk: kategori, produk, barcode,
// penyesuaian stok, dan hasil lookup barcode.
package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Kategori struct {
	ID    uint64
	Nama  string
	Aktif bool
}

type Status string

const (
	StatusAktif Status = "aktif"
	// StatusDraft: produk belum siap jual (harga normal masih 0), tidak tampil di kasir.
	StatusDraft Status = "draft"
)

// StatusDariHarga: produk tanpa harga jual otomatis berstatus draft.
func StatusDariHarga(hargaNormal decimal.Decimal) Status {
	if hargaNormal.IsPositive() {
		return StatusAktif
	}
	return StatusDraft
}

type Produk struct {
	ID              uint64
	KategoriID      uint64
	KategoriNama    string
	Nama            string
	HargaNormal     decimal.Decimal
	BisaDitawar     bool
	MarginMinPersen *decimal.Decimal
	Stok            int64
	HPPRataRata     decimal.Decimal
	AturanPakai     *string
	Status          Status
	Aktif           bool
	CreatedAt       time.Time
	Barcodes        []Barcode // hanya terisi di detail
}

type Barcode struct {
	ID            uint64
	ProdukID      uint64
	Kode          string // sudah dinormalisasi
	Jenis         string // EAN13, EAN8, UPCA, INTERNAL
	IsiPerKemasan int32
	Sumber        string // manual, openfoodfacts, upcitemdb
	CreatedAt     time.Time
}

type ProdukFilter struct {
	Q          *string
	KategoriID *uint64
	Status     *Status
	Aktif      *bool
	Limit      int32
	Offset     int32
}

// ProdukInput dipakai untuk tambah dan ubah produk. Stok dan HPP tidak
// pernah diisi dari sini.
type ProdukInput struct {
	KategoriID      uint64
	Nama            string
	HargaNormal     decimal.Decimal
	BisaDitawar     bool
	MarginMinPersen *decimal.Decimal
	AturanPakai     *string
	Aktif           bool // hanya dipakai saat ubah
}

type BarcodeInput struct {
	Kode          string // mentah, dinormalisasi di service
	IsiPerKemasan int32
	Sumber        string // kosong = manual
}

type JenisPenyesuaian string

const (
	JenisStokAwal JenisPenyesuaian = "stok_awal"
	JenisOpname   JenisPenyesuaian = "opname"
	JenisRusak    JenisPenyesuaian = "rusak"
)

// PenyesuaianInput:
//   - stok_awal: Qty (> 0) dan HargaModal per unit
//   - opname:    StokFisik (hasil hitung fisik, >= 0)
//   - rusak:     Qty = jumlah unit rusak (> 0)
type PenyesuaianInput struct {
	ProdukID   uint64
	UserID     uint64
	Jenis      JenisPenyesuaian
	Qty        int64
	StokFisik  *int64
	HargaModal decimal.Decimal
	Keterangan string
}

type PenyesuaianStok struct {
	ID          uint64
	Tanggal     time.Time
	UserID      uint64
	UserNama    string
	ProdukID    uint64
	Jenis       JenisPenyesuaian
	Qty         int64 // perubahan stok, boleh minus
	StokSebelum int64
	StokSesudah int64
	HPP         decimal.Decimal // nilai per unit saat penyesuaian
	Keterangan  *string
}

// ExternalProduct = hasil satu API publik barcode.
type ExternalProduct struct {
	Found bool
	Nama  string
	Merek string
	Raw   []byte // JSON asli, disimpan di cache
}

// Sumber hasil lookup.
const (
	SumberLokal = "lokal"
	SumberSemua = "semua" // dicari di semua API, tidak ditemukan
)

type LookupResult struct {
	Kode      string
	Jenis     string
	Ditemukan bool
	Sumber    string // lokal, openfoodfacts, upcitemdb, semua
	DariCache bool
	Nama      *string // saran nama dari API, user wajib cek
	Merek     *string
	Produk    *Produk // terisi bila barcode sudah terdaftar di toko
	// SebagianGagal: ada API yang error/timeout, jadi "tidak ditemukan"
	// belum pasti dan hasilnya tidak di-cache.
	SebagianGagal bool
}
