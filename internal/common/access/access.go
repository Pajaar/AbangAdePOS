// Package access berisi kontrak pengecekan izin aksi (menus.tipe = 'aksi')
// yang dipakai executor, tanpa bergantung langsung ke authorize.Store.
package access

// ActionLihatHPP: role boleh menerima HPP, HPP rata-rata, harga nett/tengah,
// dan margin di respons API. Tanpa izin ini field tersebut tidak dikirim.
const ActionLihatHPP = "lihat_hpp"

type IActionChecker interface {
	HasAction(roleKode, aksiKode string) bool
}
