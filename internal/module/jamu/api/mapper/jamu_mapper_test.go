package mapper

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/shopspring/decimal"

	"pos-jamu-backend/internal/module/jamu/tck/domain"
)

func TestToJamuResponse(t *testing.T) {
	j := domain.Jamu{ID: 1, Nama: "Buyung Upik Seduh", HargaJual: decimal.NewFromInt(5000), ProdukHPP: decimal.NewFromInt(1800)}

	raw, _ := json.Marshal(ToJamuResponse(j, false))
	if strings.Contains(string(raw), `"hpp"`) || strings.Contains(string(raw), `"margin"`) {
		t.Errorf("hpp/margin bocor ke role tanpa lihat_hpp: %s", raw)
	}

	res := ToJamuResponse(j, true)
	if res.HPP == nil || !res.HPP.Equal(decimal.NewFromInt(1800)) {
		t.Errorf("hpp = %v", res.HPP)
	}
	if res.Margin == nil || !res.Margin.Equal(decimal.NewFromInt(3200)) {
		t.Errorf("margin = %v", res.Margin)
	}
}

func TestToAdditionalResponse(t *testing.T) {
	telur := domain.AdditionalItem{ID: 1, Nama: "Telur", HargaJual: decimal.NewFromInt(3000), ProdukHPP: decimal.Zero}
	res := ToAdditionalResponse(telur, false)
	if res.DilacakStok || res.HPP != nil {
		t.Errorf("telur tidak dilacak dan tanpa hpp: %+v", res)
	}

	pid := uint64(9)
	madu := domain.AdditionalItem{ID: 2, Nama: "Madu", ProdukID: &pid, ProdukHPP: decimal.NewFromInt(1000)}
	if res := ToAdditionalResponse(madu, true); !res.DilacakStok || res.HPP == nil || !res.HPP.Equal(decimal.NewFromInt(1000)) {
		t.Errorf("madu dilacak dengan hpp: %+v", res)
	}
}
