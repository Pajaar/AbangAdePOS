package inservice

import (
	"context"
	"database/sql"
	"errors"

	"pos-jamu-backend/internal/auth"
	"pos-jamu-backend/internal/module/auth/service/repository"
	"pos-jamu-backend/internal/module/auth/tck/domain"
	tck "pos-jamu-backend/internal/module/auth/tck/inservice"
	"pos-jamu-backend/internal/util/apperror"
)

type UserInService struct {
	db *sql.DB
	// dummyPasswordHash dipakai sebagai target bcrypt.Compare kalau user
	// tidak ketemu / nonaktif, supaya waktu eksekusi tidak membedakan
	// kasus itu dari "password salah" (mitigasi timing attack).
	dummyPasswordHash string
}

var _ tck.IUserInService = (*UserInService)(nil)

func NewUserInService(db *sql.DB) (*UserInService, error) {
	dummyHash, err := auth.HashPassword("dummy-password-untuk-constant-time-compare")
	if err != nil {
		return nil, err
	}
	return &UserInService{db: db, dummyPasswordHash: dummyHash}, nil
}

func (s *UserInService) Authenticate(ctx context.Context, username, password string) (domain.User, error) {
	user, err := repository.NewUserRepository(s.db).GetByUsername(ctx, username)
	notFound := errors.Is(err, sql.ErrNoRows)
	if err != nil && !notFound {
		return domain.User{}, apperror.Internal(err)
	}

	validUser := !notFound && user.Aktif
	hashToCompare := s.dummyPasswordHash
	if validUser {
		hashToCompare = user.PasswordHash
	}

	cmpErr := auth.ComparePassword(hashToCompare, password)
	if !validUser || cmpErr != nil {
		return domain.User{}, errLoginFailed
	}
	return user, nil
}

func (s *UserInService) GetByID(ctx context.Context, id uint64) (domain.User, error) {
	user, err := repository.NewUserRepository(s.db).GetByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, apperror.NotFound("User tidak ditemukan")
	}
	if err != nil {
		return domain.User{}, apperror.Internal(err)
	}
	return user, nil
}

func (s *UserInService) Create(ctx context.Context, u domain.NewUser) (uint64, error) {
	roleID, err := repository.NewRoleRepository(s.db).GetIDByKode(ctx, u.RoleKode)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, apperror.BadRequest("Role " + u.RoleKode + " tidak ditemukan")
	}
	if err != nil {
		return 0, apperror.Internal(err)
	}

	users := repository.NewUserRepository(s.db)
	if _, err := users.GetByUsername(ctx, u.Username); err == nil {
		return 0, apperror.BadRequest("Username " + u.Username + " sudah dipakai")
	} else if !errors.Is(err, sql.ErrNoRows) {
		return 0, apperror.Internal(err)
	}

	hash, err := auth.HashPassword(u.Password)
	if err != nil {
		return 0, apperror.Internal(err)
	}

	id, err := users.Create(ctx, u.Nama, u.Username, hash, roleID)
	if err != nil {
		return 0, apperror.Internal(err)
	}
	return id, nil
}
