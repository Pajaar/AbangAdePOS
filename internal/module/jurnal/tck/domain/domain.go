// Package domain berisi tipe jurnal umum dan kode akun yang dipakai kode Go.
package domain

import "github.com/shopspring/decimal"

// Kode akun dari seed migrasi 000004. Nama akun boleh diubah lewat data,
// tetapi kode ini dipakai oleh posting jurnal otomatis.
const (
	AkunKas               = "1-101"
	AkunPiutangUsaha      = "1-102"
	AkunPersediaan        = "1-103"
	AkunUtangUsaha        = "2-101"
	AkunModalPemilik      = "3-101"
	AkunPenjualan         = "4-101"
	AkunPotonganPenjualan = "4-102"
	AkunHPP               = "5-101"
	AkunBebanOperasional  = "6-101"
)

// Sumber jurnal (jurnal.sumber_tipe).
const (
	SumberPenyesuaianStok = "penyesuaian_stok"
)

// Baris berisi debit ATAU kredit (yang lain nol).
type Baris struct {
	AkunKode string
	Debit    decimal.Decimal
	Kredit   decimal.Decimal
}

func Debit(akun string, jumlah decimal.Decimal) Baris {
	return Baris{AkunKode: akun, Debit: jumlah, Kredit: decimal.Zero}
}

func Kredit(akun string, jumlah decimal.Decimal) Baris {
	return Baris{AkunKode: akun, Debit: decimal.Zero, Kredit: jumlah}
}

type Entri struct {
	UserID     uint64
	Keterangan string
	SumberTipe string
	SumberID   uint64
	Baris      []Baris
}
