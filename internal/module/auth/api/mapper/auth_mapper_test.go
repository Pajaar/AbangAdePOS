package mapper

import (
	"testing"

	"pos-jamu-backend/internal/module/auth/tck/domain"
)

func TestToMenusResponse(t *testing.T) {
	parent := uint64(1)
	menus := []domain.Menu{
		{ID: 1, Kode: "laporan", Tipe: domain.MenuTipeMenu, Urutan: 1},
		{ID: 2, Kode: "laba_rugi", Tipe: domain.MenuTipeMenu, ParentID: &parent, Urutan: 2},
		{ID: 3, Kode: "kasir", Tipe: domain.MenuTipeMenu, Urutan: 3},
		{ID: 9, Kode: "lihat_hpp", Tipe: domain.MenuTipeAksi, Urutan: 90},
	}

	res := ToMenusResponse(menus)

	if len(res.Menus) != 2 || res.Menus[0].Kode != "laporan" || res.Menus[1].Kode != "kasir" {
		t.Fatalf("root menu salah: %+v", res.Menus)
	}
	if len(res.Menus[0].Children) != 1 || res.Menus[0].Children[0].Kode != "laba_rugi" {
		t.Errorf("child menu salah: %+v", res.Menus[0].Children)
	}
	if len(res.Actions) != 1 || res.Actions[0] != "lihat_hpp" {
		t.Errorf("actions salah: %v", res.Actions)
	}
}
