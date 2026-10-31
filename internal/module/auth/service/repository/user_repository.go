// Package repository membungkus query sqlc modul auth (setara
// service/repository di primecash) dan memetakan hasilnya ke tipe domain.
// Konstruktor menerima sqlcgen.DBTX supaya bisa dipakai dengan *sql.DB
// maupun *sql.Tx dari database.WithTx.
package repository

import (
	"context"

	"pos-jamu-backend/internal/module/auth/tck/domain"
	sqlcgen "pos-jamu-backend/internal/repository/sqlc"
)

type IUserRepository interface {
	GetByUsername(ctx context.Context, username string) (domain.User, error)
	GetByID(ctx context.Context, id uint64) (domain.User, error)
	UpdateLastLogin(ctx context.Context, id uint64) error
	Create(ctx context.Context, nama, username, passwordHash string, roleID uint64) (uint64, error)
}

type userRepository struct {
	q *sqlcgen.Queries
}

func NewUserRepository(db sqlcgen.DBTX) IUserRepository {
	return &userRepository{q: sqlcgen.New(db)}
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	row, err := r.q.GetUserByUsername(ctx, username)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		ID: row.ID, Nama: row.Nama, Username: row.Username, PasswordHash: row.PasswordHash,
		RoleID: row.RoleID, RoleKode: row.RoleKode, Aktif: row.Aktif,
	}, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uint64) (domain.User, error) {
	row, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		ID: row.ID, Nama: row.Nama, Username: row.Username, PasswordHash: row.PasswordHash,
		RoleID: row.RoleID, RoleKode: row.RoleKode, Aktif: row.Aktif,
	}, nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, id uint64) error {
	return r.q.UpdateLastLogin(ctx, id)
}

func (r *userRepository) Create(ctx context.Context, nama, username, passwordHash string, roleID uint64) (uint64, error) {
	id, err := r.q.CreateUser(ctx, sqlcgen.CreateUserParams{
		Nama: nama, Username: username, PasswordHash: passwordHash, RoleID: roleID,
	})
	return uint64(id), err
}
