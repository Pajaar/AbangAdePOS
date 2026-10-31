// Package mapper memetakan tipe domain ke DTO respons (setara api/mapper
// di primecash). Executor tidak menyusun DTO sendiri.
package mapper

import (
	"pos-jamu-backend/internal/dto"
	"pos-jamu-backend/internal/module/auth/tck/domain"
)

func ToUserInfo(u domain.User) dto.UserInfo {
	return dto.UserInfo{
		ID:       u.ID,
		Nama:     u.Nama,
		Username: u.Username,
		Role:     u.RoleKode,
	}
}

// ToMenusResponse memisahkan menu: tipe 'menu' disusun jadi tree lewat
// ParentID, tipe 'aksi' dikumpulkan jadi daftar kode (misalnya lihat_hpp).
// Urutan input (urutan menu) dipertahankan.
func ToMenusResponse(menus []domain.Menu) dto.MenusResponse {
	nodes := make(map[uint64]*dto.MenuNode)
	actions := []string{}

	for _, m := range menus {
		if m.Tipe == domain.MenuTipeAksi {
			actions = append(actions, m.Kode)
			continue
		}
		nodes[m.ID] = &dto.MenuNode{
			ID:     m.ID,
			Kode:   m.Kode,
			Nama:   m.Nama,
			Path:   m.Path,
			Icon:   m.Icon,
			Urutan: m.Urutan,
		}
	}

	var roots []*dto.MenuNode
	for _, m := range menus {
		if m.Tipe != domain.MenuTipeMenu {
			continue
		}
		node := nodes[m.ID]
		if m.ParentID != nil {
			if parent, ok := nodes[*m.ParentID]; ok {
				parent.Children = append(parent.Children, node)
				continue
			}
		}
		roots = append(roots, node)
	}

	return dto.MenusResponse{Menus: roots, Actions: actions}
}
