// Package domain berisi tipe domain modul auth (setara tck/domain di
// primecash). Lapisan api dan service bertukar data lewat tipe ini, bukan
// lewat struct hasil sqlc.
package domain

import "time"

type User struct {
	ID           uint64
	Nama         string
	Username     string
	PasswordHash string
	RoleID       uint64
	RoleKode     string
	Aktif        bool
}

type NewUser struct {
	Nama     string
	Username string
	Password string // plain, di-hash di service
	RoleKode string
}

type MenuTipe string

const (
	MenuTipeMenu MenuTipe = "menu"
	MenuTipeAksi MenuTipe = "aksi"
)

type Menu struct {
	ID       uint64
	ParentID *uint64
	Kode     string
	Nama     string
	Tipe     MenuTipe
	Path     *string
	Icon     *string
	Urutan   int32
}

// Permission = satu endpoint yang boleh dipanggil sebuah role.
type Permission struct {
	RoleKode string
	Method   string
	Path     string // pola route Gin, misalnya /api/v1/produk/:id
}

type RefreshToken struct {
	ID        uint64
	UserID    uint64
	ExpiresAt time.Time
	Revoked   bool
}

// ClientInfo dicatat bersama refresh token untuk jejak audit.
type ClientInfo struct {
	UserAgent string
	IPAddress string
}

// Session = refresh token baru hasil login atau rotasi.
// RefreshToken berisi token asli (bukan hash), hanya untuk dikirim di cookie.
type Session struct {
	UserID           uint64
	RoleKode         string
	RefreshToken     string
	RefreshExpiresAt time.Time
}
