package inservice

import (
	"testing"

	"github.com/shopspring/decimal"

	jurnal "pos-jamu-backend/internal/module/jurnal/service/inservice"
	jd "pos-jamu-backend/internal/module/jurnal/tck/domain"
	"pos-jamu-backend/internal/module/produk/tck/domain"
)

func dec(s string) decimal.Decimal { return decimal.RequireFromString(s) }

func TestHitungPenyesuaian(t *testing.T) {
	fisik := func(n int64) *int64 { return &n }

	cases := []struct {
		name           string
		in             domain.PenyesuaianInput
		stok           int64
		hpp            string
		qty, stokBaru  int64
		hppBaru, nilai string
		debit, kredit  string // akun
	}{
		{"stok awal dari nol", domain.PenyesuaianInput{Jenis: domain.JenisStokAwal, Qty: 10, HargaModal: dec("2000")},
			0, "0", 10, 10, "2000", "20000", jd.AkunPersediaan, jd.AkunModalPemilik},
		{"stok awal rata-rata", domain.PenyesuaianInput{Jenis: domain.JenisStokAwal, Qty: 10, HargaModal: dec("3000")},
			10, "2000", 10, 20, "2500", "30000", jd.AkunPersediaan, jd.AkunModalPemilik},
		{"opname berkurang", domain.PenyesuaianInput{Jenis: domain.JenisOpname, StokFisik: fisik(7)},
			10, "2500", -3, 7, "2500", "7500", jd.AkunHPP, jd.AkunPersediaan},
		{"opname bertambah", domain.PenyesuaianInput{Jenis: domain.JenisOpname, StokFisik: fisik(12)},
			10, "2500", 2, 12, "2500", "5000", jd.AkunPersediaan, jd.AkunHPP},
		{"rusak", domain.PenyesuaianInput{Jenis: domain.JenisRusak, Qty: 2},
			5, "1000", -2, 3, "1000", "2000", jd.AkunHPP, jd.AkunPersediaan},
		{"rusak boleh bikin stok minus", domain.PenyesuaianInput{Jenis: domain.JenisRusak, Qty: 3},
			1, "1000", -3, -2, "1000", "3000", jd.AkunHPP, jd.AkunPersediaan},
	}

	for _, c := range cases {
		r, err := HitungPenyesuaian(c.in, c.stok, dec(c.hpp))
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if r.qty != c.qty || r.stokBaru != c.stokBaru || !r.hppBaru.Equal(dec(c.hppBaru)) {
			t.Errorf("%s: qty=%d stok=%d hpp=%s", c.name, r.qty, r.stokBaru, r.hppBaru)
		}
		baris, err := jurnal.ValidasiBaris(r.baris)
		if err != nil || len(baris) != 2 {
			t.Fatalf("%s: jurnal tidak valid: %v", c.name, err)
		}
		if baris[0].AkunKode != c.debit || !baris[0].Debit.Equal(dec(c.nilai)) ||
			baris[1].AkunKode != c.kredit || !baris[1].Kredit.Equal(dec(c.nilai)) {
			t.Errorf("%s: jurnal = %+v", c.name, baris)
		}
	}
}

func TestHitungPenyesuaian_Ditolak(t *testing.T) {
	sama := int64(5)
	minus := int64(-1)
	bad := []domain.PenyesuaianInput{
		{Jenis: domain.JenisStokAwal, Qty: 0, HargaModal: dec("1000")},
		{Jenis: domain.JenisStokAwal, Qty: 1, HargaModal: dec("-1")},
		{Jenis: domain.JenisOpname},
		{Jenis: domain.JenisOpname, StokFisik: &sama},
		{Jenis: domain.JenisOpname, StokFisik: &minus},
		{Jenis: domain.JenisRusak, Qty: 0},
		{Jenis: "hilang", Qty: 1},
	}
	for i, in := range bad {
		if _, err := HitungPenyesuaian(in, 5, dec("1000")); err == nil {
			t.Errorf("kasus %d (%s) seharusnya ditolak", i, in.Jenis)
		}
	}
}

func TestHitungPenyesuaian_HPPNolTanpaJurnal(t *testing.T) {
	stok := int64(3)
	r, err := HitungPenyesuaian(domain.PenyesuaianInput{Jenis: domain.JenisOpname, StokFisik: &stok}, 5, decimal.Zero)
	if err != nil {
		t.Fatal(err)
	}
	baris, err := jurnal.ValidasiBaris(r.baris)
	if err != nil || len(baris) != 0 {
		t.Errorf("HPP 0 = jurnal kosong, got %v %v", baris, err)
	}
}
