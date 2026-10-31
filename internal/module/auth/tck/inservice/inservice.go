// Package inservice berisi kontrak service modul auth (setara tck/internal
// di primecash). Executor dan komponen lain hanya bergantung pada interface
// ini, implementasinya ada di service/inservice.
//
// Nama folder bukan "internal" karena di Go folder itu membatasi siapa yang
// boleh meng-import.
package inservice

import (
	"context"

	"pos-jamu-backend/internal/module/auth/tck/domain"
)

type IUserInService interface {
	// Authenticate gagal dengan pesan umum yang sama untuk user tidak ada,
	// nonaktif, atau password salah.
	Authenticate(ctx context.Context, username, password string) (domain.User, error)
	GetByID(ctx context.Context, id uint64) (domain.User, error)
	Create(ctx context.Context, u domain.NewUser) (uint64, error)
}

type ISessionInService interface {
	// Start membuat refresh token baru dan mencatat waktu login.
	Start(ctx context.Context, user domain.User, client domain.ClientInfo) (domain.Session, error)
	// Rotate mengganti refresh token lama dengan yang baru. Token yang sudah
	// di-revoke dipakai lagi = indikasi dicuri: semua token user di-revoke.
	Rotate(ctx context.Context, refreshToken string, client domain.ClientInfo) (domain.Session, error)
	// End me-revoke refresh token (logout). Idempoten.
	End(ctx context.Context, refreshToken string) error
}

type IMenuInService interface {
	ListByRole(ctx context.Context, roleKode string) ([]domain.Menu, error)
}

type IPermissionInService interface {
	ListAll(ctx context.Context) ([]domain.Permission, error)
	ListAllActions(ctx context.Context) ([]domain.RoleAction, error)
}
