package dto

import "time"

// LoginRequest: UserAgent dan IPAddress diisi controller dari HTTP request,
// tidak dibaca dari body JSON.
type LoginRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	UserAgent string `json:"-"`
	IPAddress string `json:"-"`
}

// RefreshRequest dan LogoutRequest diisi controller dari cookie.
type RefreshRequest struct {
	RefreshToken string
	UserAgent    string
	IPAddress    string
}

type LogoutRequest struct {
	RefreshToken string
}

// MeRequest dan MenusRequest diisi controller dari context hasil middleware Auth.
type MeRequest struct {
	UserID uint64
}

type MenusRequest struct {
	RoleKode string
}

// LoginResult dan RefreshResult adalah hasil executor. Response dikirim
// sebagai body JSON; RefreshToken dipasang controller di cookie httpOnly.
type LoginResult struct {
	Response         LoginResponse
	RefreshToken     string
	RefreshExpiresAt time.Time
}

type RefreshResult struct {
	Response         RefreshResponse
	RefreshToken     string
	RefreshExpiresAt time.Time
}

type UserInfo struct {
	ID       uint64 `json:"id"`
	Nama     string `json:"nama"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type LoginResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
	User        UserInfo  `json:"user"`
}

type RefreshResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type MeResponse struct {
	UserInfo
}

// MenuNode adalah satu baris menu sidebar (tipe='menu'), disusun jadi tree
// lewat parent_id. Menu tipe='aksi' (misal lihat_hpp) TIDAK muncul di sini,
// lihat MenusResponse.Actions.
type MenuNode struct {
	ID       uint64      `json:"id"`
	Kode     string      `json:"kode"`
	Nama     string      `json:"nama"`
	Path     *string     `json:"path,omitempty"`
	Icon     *string     `json:"icon,omitempty"`
	Urutan   int32       `json:"urutan"`
	Children []*MenuNode `json:"children,omitempty"`
}

// MenusResponse: Menus untuk sidebar, Actions untuk izin tanpa halaman
// (misal "lihat_hpp") supaya frontend bisa toggle tampilan tanpa endpoint lain.
type MenusResponse struct {
	Menus   []*MenuNode `json:"menus"`
	Actions []string    `json:"actions"`
}
