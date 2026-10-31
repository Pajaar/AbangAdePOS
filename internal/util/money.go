package util

import "github.com/shopspring/decimal"

// Semua perhitungan uang memakai decimal.Decimal (aturan wajib nomor 1).

var (
	kelipatanHarga = decimal.NewFromInt(500)
	seratus        = decimal.NewFromInt(100)
	dua            = decimal.NewFromInt(2)
)

// Round2 membulatkan ke 2 desimal (DECIMAL(14,2)), setengah ke atas.
func Round2(d decimal.Decimal) decimal.Decimal {
	return d.Round(2)
}

// CeilToMultiple membulatkan ke atas ke kelipatan m (m > 0).
func CeilToMultiple(d, m decimal.Decimal) decimal.Decimal {
	return d.Div(m).Ceil().Mul(m)
}

// HPPRataRataBaru menghitung HPP moving average setelah barang masuk.
//
//	hpp_baru = (stok_lama * hpp_lama + qty * harga) / (stok_lama + qty)
//	jika stok_lama <= 0: hpp_baru = harga
func HPPRataRataBaru(stokLama int64, hppLama decimal.Decimal, qty int64, harga decimal.Decimal) decimal.Decimal {
	if stokLama <= 0 {
		return Round2(harga)
	}
	sl := decimal.NewFromInt(stokLama)
	q := decimal.NewFromInt(qty)
	total := sl.Mul(hppLama).Add(q.Mul(harga))
	return Round2(total.Div(sl.Add(q)))
}

// HargaNett = hpp * (1 + margin/100), dibulatkan ke atas ke kelipatan Rp500.
func HargaNett(hpp, marginPersen decimal.Decimal) decimal.Decimal {
	faktor := decimal.NewFromInt(1).Add(marginPersen.Div(seratus))
	return CeilToMultiple(hpp.Mul(faktor), kelipatanHarga)
}

// HargaTengah = titik tengah harga normal dan nett, dibulatkan ke atas ke kelipatan Rp500.
func HargaTengah(hargaNormal, nett decimal.Decimal) decimal.Decimal {
	return CeilToMultiple(hargaNormal.Add(nett).Div(dua), kelipatanHarga)
}
