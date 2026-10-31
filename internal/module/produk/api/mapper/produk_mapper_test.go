package mapper

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/shopspring/decimal"

	"pos-jamu-backend/internal/module/produk/tck/domain"
)

func TestToProdukResponse_FilterHPP(t *testing.T) {
	margin := decimal.NewFromInt(20)
	p := domain.Produk{
		ID: 1, Nama: "Obat Batuk", HargaNormal: decimal.NewFromInt(15000), BisaDitawar: true,
		HPPRataRata: decimal.NewFromInt(10000), MarginMinPersen: &margin, Status: domain.StatusAktif, Aktif: true,
	}

	// Tanpa izin: field HPP tidak boleh ada di JSON sama sekali.
	raw, _ := json.Marshal(ToProdukResponse(p, false))
	for _, field := range []string{"hpp_rata_rata", "margin_min_persen", "harga_nett", "harga_tengah", "harga_perlu_dinaikkan"} {
		if strings.Contains(string(raw), field) {
			t.Errorf("field %s bocor ke role tanpa lihat_hpp: %s", field, raw)
		}
	}

	// Dengan izin: nett = 10000 * 1.2 = 12000, tengah = (15000+12000)/2 = 13500.
	res := ToProdukResponse(p, true)
	if res.HPPRataRata == nil || !res.HPPRataRata.Equal(decimal.NewFromInt(10000)) {
		t.Errorf("hpp = %v", res.HPPRataRata)
	}
	if res.HargaNett == nil || !res.HargaNett.Equal(decimal.NewFromInt(12000)) {
		t.Errorf("nett = %v", res.HargaNett)
	}
	if res.HargaTengah == nil || !res.HargaTengah.Equal(decimal.NewFromInt(13500)) {
		t.Errorf("tengah = %v", res.HargaTengah)
	}
	if res.HargaPerluDinaikkan == nil || *res.HargaPerluDinaikkan {
		t.Errorf("perlu dinaikkan = %v", res.HargaPerluDinaikkan)
	}

	// Nett >= harga normal: tandai harga perlu dinaikkan.
	p.HargaNormal = decimal.NewFromInt(11000)
	if res := ToProdukResponse(p, true); res.HargaPerluDinaikkan == nil || !*res.HargaPerluDinaikkan {
		t.Error("harga_perlu_dinaikkan harus true")
	}

	// HPP belum diketahui: tidak ada nett/tengah.
	p.HPPRataRata = decimal.Zero
	if res := ToProdukResponse(p, true); res.HargaNett != nil {
		t.Error("HPP 0 tidak boleh menghasilkan harga nett")
	}
}

func TestToPenyesuaianResponse_FilterHPP(t *testing.T) {
	ps := domain.PenyesuaianStok{ID: 1, HPP: decimal.NewFromInt(2000)}
	raw, _ := json.Marshal(ToPenyesuaianResponse(ps, false))
	if strings.Contains(string(raw), `"hpp"`) {
		t.Errorf("hpp bocor: %s", raw)
	}
	if res := ToPenyesuaianResponse(ps, true); res.HPP == nil {
		t.Error("hpp harus ada untuk lihat_hpp")
	}
}
