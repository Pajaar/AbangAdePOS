package inservice

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"pos-jamu-backend/internal/auth"
	"pos-jamu-backend/internal/database"
	"pos-jamu-backend/internal/module/auth/service/repository"
	"pos-jamu-backend/internal/module/auth/tck/domain"
	tck "pos-jamu-backend/internal/module/auth/tck/inservice"
	"pos-jamu-backend/internal/util/apperror"
)

type SessionInService struct {
	db              *sql.DB
	refreshTokenTTL time.Duration
}

var _ tck.ISessionInService = (*SessionInService)(nil)

func NewSessionInService(db *sql.DB, refreshTokenTTL time.Duration) *SessionInService {
	return &SessionInService{db: db, refreshTokenTTL: refreshTokenTTL}
}

func (s *SessionInService) Start(ctx context.Context, user domain.User, client domain.ClientInfo) (domain.Session, error) {
	plain, err := auth.GenerateRefreshToken()
	if err != nil {
		return domain.Session{}, apperror.Internal(err)
	}
	expiresAt := time.Now().Add(s.refreshTokenTTL)

	err = database.WithTx(ctx, s.db, func(tx *sql.Tx) error {
		if _, err := repository.NewRefreshTokenRepository(tx).Create(ctx, user.ID, auth.HashRefreshToken(plain), expiresAt, client); err != nil {
			return err
		}
		return repository.NewUserRepository(tx).UpdateLastLogin(ctx, user.ID)
	})
	if err != nil {
		return domain.Session{}, apperror.Internal(err)
	}

	return domain.Session{
		UserID:           user.ID,
		RoleKode:         user.RoleKode,
		RefreshToken:     plain,
		RefreshExpiresAt: expiresAt,
	}, nil
}

// Rotate menjalankan cek dan rotasi dalam SATU transaksi dengan baris token
// dikunci (SELECT ... FOR UPDATE). Request kedua yang membawa token sama
// akan menunggu, lalu melihat token sudah di-revoke dan masuk jalur reuse.
func (s *SessionInService) Rotate(ctx context.Context, oldPlain string, client domain.ClientInfo) (domain.Session, error) {
	if oldPlain == "" {
		return domain.Session{}, errSessionInvalid
	}

	newPlain, err := auth.GenerateRefreshToken()
	if err != nil {
		return domain.Session{}, apperror.Internal(err)
	}
	newExpiresAt := time.Now().Add(s.refreshTokenTTL)

	var (
		session domain.Session
		reused  bool
	)
	err = database.WithTx(ctx, s.db, func(tx *sql.Tx) error {
		tokens := repository.NewRefreshTokenRepository(tx)

		old, err := tokens.GetByHashForUpdate(ctx, auth.HashRefreshToken(oldPlain))
		if errors.Is(err, sql.ErrNoRows) {
			return errSessionInvalid
		}
		if err != nil {
			return err
		}

		if old.Revoked {
			// Return nil supaya revoke massal ini ikut ter-commit;
			// errornya dikembalikan setelah transaksi selesai.
			reused = true
			return tokens.RevokeAllByUser(ctx, old.UserID)
		}
		if time.Now().After(old.ExpiresAt) {
			return errSessionInvalid
		}

		user, err := repository.NewUserRepository(tx).GetByID(ctx, old.UserID)
		if errors.Is(err, sql.ErrNoRows) {
			return errSessionInvalid
		}
		if err != nil {
			return err
		}
		if !user.Aktif {
			return errSessionInvalid
		}

		newID, err := tokens.Create(ctx, old.UserID, auth.HashRefreshToken(newPlain), newExpiresAt, client)
		if err != nil {
			return err
		}
		rotated, err := tokens.Rotate(ctx, old.ID, newID)
		if err != nil {
			return err
		}
		if !rotated {
			// Pengaman tambahan; dengan FOR UPDATE seharusnya tidak terjadi.
			return errSessionInvalid
		}

		session = domain.Session{
			UserID:           user.ID,
			RoleKode:         user.RoleKode,
			RefreshToken:     newPlain,
			RefreshExpiresAt: newExpiresAt,
		}
		return nil
	})
	if err != nil {
		return domain.Session{}, asAppError(err)
	}
	if reused {
		return domain.Session{}, errSessionInvalid
	}
	return session, nil
}

// End idempoten: token kosong, tidak dikenal, atau sudah di-revoke tetap
// dianggap sukses supaya tidak membocorkan status token ke klien.
func (s *SessionInService) End(ctx context.Context, plain string) error {
	if plain == "" {
		return nil
	}
	if err := repository.NewRefreshTokenRepository(s.db).RevokeByHash(ctx, auth.HashRefreshToken(plain)); err != nil {
		return apperror.Internal(err)
	}
	return nil
}
