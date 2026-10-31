// Package mapper memetakan domain produk ke DTO. Semua penyaringan field HPP
// (izin lihat_hpp) terjadi di sini lewat parameter canSeeHPP.
package mapper

import (
	"github.com/shopspring/decimal"

	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/module/produk/tck/domain"
	"pos-jamu-backend/internal/util"
)

func ToKategoriResponse(k domain.Kategori) dto.KategoriResponse {
	return dto.KategoriResponse{ID: k.ID, Nama: k.Nama, Aktif: k.Aktif}
}

func ToBarcodeResponse(b domain.Barcode) dto.BarcodeResponse {
	return dto.BarcodeResponse{
		ID: b.ID, ProdukID: b.ProdukID, Kode: b.Kode, Jenis: b.Jenis,
		IsiPerKemasan: b.IsiPerKemasan, Sumber: b.Sumber, CreatedAt: b.CreatedAt,
	}
}

func ToProdukResponse(p domain.Produk, canSeeHPP bool) dto.ProdukResponse {
	res := dto.ProdukResponse{
		ID:           p.ID,
		KategoriID:   p.KategoriID,
		KategoriNama: p.KategoriNama,
		Nama:         p.Nama,
		HargaNormal:  p.HargaNormal,
		BisaDitawar:  p.BisaDitawar,
		Stok:         p.Stok,
		AturanPakai:  p.AturanPakai,
		Status:       string(p.Status),
		Aktif:        p.Aktif,
		CreatedAt:    p.CreatedAt,
	}
	for _, b := range p.Barcodes {
		res.Barcodes = append(res.Barcodes, ToBarcodeResponse(b))
	}
	if !canSeeHPP {
		return res
	}

	hpp := p.HPPRataRata
	res.HPPRataRata = &hpp
	res.MarginMinPersen = p.MarginMinPersen

	// Nett dan tengah tidak disimpan, dihitung saat ditampilkan. Hanya untuk
	// produk yang bisa ditawar dan HPP-nya sudah diketahui.
	if p.BisaDitawar && hpp.IsPositive() {
		margin := decimal.Zero
		if p.MarginMinPersen != nil {
			margin = *p.MarginMinPersen
		}
		nett := util.HargaNett(hpp, margin)
		tengah := util.HargaTengah(p.HargaNormal, nett)
		perlu := nett.GreaterThanOrEqual(p.HargaNormal)
		res.HargaNett, res.HargaTengah, res.HargaPerluDinaikkan = &nett, &tengah, &perlu
	}
	return res
}

func ToLookupResponse(r domain.LookupResult, canSeeHPP bool) dto.BarcodeLookupResponse {
	res := dto.BarcodeLookupResponse{
		Kode: r.Kode, Jenis: r.Jenis, Ditemukan: r.Ditemukan, Sumber: r.Sumber,
		DariCache: r.DariCache, Nama: r.Nama, Merek: r.Merek, SebagianGagal: r.SebagianGagal,
	}
	if r.Produk != nil {
		p := ToProdukResponse(*r.Produk, canSeeHPP)
		res.Produk = &p
	}
	return res
}

func ToPenyesuaianResponse(p domain.PenyesuaianStok, canSeeHPP bool) dto.PenyesuaianResponse {
	res := dto.PenyesuaianResponse{
		ID: p.ID, Tanggal: p.Tanggal, UserNama: p.UserNama, ProdukID: p.ProdukID,
		Jenis: string(p.Jenis), Qty: p.Qty, StokSebelum: p.StokSebelum, StokSesudah: p.StokSesudah,
		Keterangan: p.Keterangan,
	}
	if canSeeHPP {
		hpp := p.HPP
		res.HPP = &hpp
	}
	return res
}
