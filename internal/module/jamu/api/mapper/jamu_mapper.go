// Package mapper memetakan domain menu jamu ke DTO, termasuk penyaringan
// field HPP lewat canSeeHPP.
package mapper

import (
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/module/jamu/tck/domain"
)

func ToJamuResponse(j domain.Jamu, canSeeHPP bool) dto.JamuResponse {
	res := dto.JamuResponse{
		ID:         j.ID,
		Nama:       j.Nama,
		HargaJual:  j.HargaJual,
		Aktif:      j.Aktif,
		CreatedAt:  j.CreatedAt,
		ProdukID:   j.ProdukID,
		ProdukNama: j.ProdukNama,
		StokSachet: j.ProdukStok,
	}
	if canSeeHPP {
		hpp := j.ProdukHPP
		margin := j.HargaJual.Sub(hpp)
		res.HPP, res.Margin = &hpp, &margin
	}
	return res
}

func ToAdditionalResponse(a domain.AdditionalItem, canSeeHPP bool) dto.AdditionalResponse {
	res := dto.AdditionalResponse{
		ID:          a.ID,
		Nama:        a.Nama,
		HargaJual:   a.HargaJual,
		ProdukID:    a.ProdukID,
		ProdukNama:  a.ProdukNama,
		DilacakStok: a.ProdukID != nil,
		Aktif:       a.Aktif,
	}
	if canSeeHPP {
		hpp := a.ProdukHPP
		res.HPP = &hpp
	}
	return res
}
