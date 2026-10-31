package inservice

import (
	"testing"

	"github.com/shopspring/decimal"

	"pos-jamu-backend/internal/module/jurnal/tck/domain"
)

func TestValidasiBaris(t *testing.T) {
	n := decimal.NewFromInt

	ok, err := ValidasiBaris([]domain.Baris{
		domain.Debit(domain.AkunPersediaan, n(20000)),
		domain.Kredit(domain.AkunModalPemilik, n(20000)),
		domain.Debit(domain.AkunHPP, decimal.Zero), // dibuang
	})
	if err != nil || len(ok) != 2 {
		t.Fatalf("seimbang harus lolos: %v %v", ok, err)
	}

	if _, err := ValidasiBaris([]domain.Baris{
		domain.Debit(domain.AkunPersediaan, n(20000)),
		domain.Kredit(domain.AkunModalPemilik, n(19000)),
	}); err == nil {
		t.Error("tidak seimbang harus gagal")
	}

	if _, err := ValidasiBaris([]domain.Baris{
		{AkunKode: domain.AkunKas, Debit: n(1), Kredit: n(1)},
	}); err == nil {
		t.Error("baris dua sisi harus gagal")
	}

	kosong, err := ValidasiBaris([]domain.Baris{domain.Debit(domain.AkunHPP, decimal.Zero), domain.Kredit(domain.AkunPersediaan, decimal.Zero)})
	if err != nil || len(kosong) != 0 {
		t.Errorf("semua nol = tidak ada jurnal: %v %v", kosong, err)
	}
}
