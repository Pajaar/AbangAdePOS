package repository

import (
	"context"
	"database/sql"
	"time"

	"pos-jamu-backend/internal/module/auth/tck/domain"
	sqlcgen "pos-jamu-backend/internal/repository/sqlc"
)

type IRefreshTokenRepository interface {
	Create(ctx context.Context, userID uint64, tokenHash string, expiresAt time.Time, client domain.ClientInfo) (uint64, error)
	// GetByHashForUpdate mengunci baris sampai transaksi selesai.
	// Wajib dipanggil dengan *sql.Tx.
	GetByHashForUpdate(ctx context.Context, tokenHash string) (domain.RefreshToken, error)
	// Rotate mengembalikan false bila token lama sudah di-revoke lebih dulu.
	Rotate(ctx context.Context, oldID, newID uint64) (bool, error)
	RevokeByHash(ctx context.Context, tokenHash string) error
	RevokeAllByUser(ctx context.Context, userID uint64) error
}

type refreshTokenRepository struct {
	q *sqlcgen.Queries
}

func NewRefreshTokenRepository(db sqlcgen.DBTX) IRefreshTokenRepository {
	return &refreshTokenRepository{q: sqlcgen.New(db)}
}

func (r *refreshTokenRepository) Create(ctx context.Context, userID uint64, tokenHash string, expiresAt time.Time, client domain.ClientInfo) (uint64, error) {
	id, err := r.q.CreateRefreshToken(ctx, sqlcgen.CreateRefreshTokenParams{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		UserAgent: nullString(client.UserAgent),
		IpAddress: nullString(client.IPAddress),
	})
	return uint64(id), err
}

func (r *refreshTokenRepository) GetByHashForUpdate(ctx context.Context, tokenHash string) (domain.RefreshToken, error) {
	row, err := r.q.GetRefreshTokenByHashForUpdate(ctx, tokenHash)
	if err != nil {
		return domain.RefreshToken{}, err
	}
	return domain.RefreshToken{
		ID:        row.ID,
		UserID:    row.UserID,
		ExpiresAt: row.ExpiresAt,
		Revoked:   row.RevokedAt.Valid,
	}, nil
}

func (r *refreshTokenRepository) Rotate(ctx context.Context, oldID, newID uint64) (bool, error) {
	n, err := r.q.RotateRefreshToken(ctx, sqlcgen.RotateRefreshTokenParams{
		ID:           oldID,
		ReplacedByID: sql.NullInt64{Int64: int64(newID), Valid: true},
	})
	return n == 1, err
}

func (r *refreshTokenRepository) RevokeByHash(ctx context.Context, tokenHash string) error {
	return r.q.RevokeRefreshTokenByHash(ctx, tokenHash)
}

func (r *refreshTokenRepository) RevokeAllByUser(ctx context.Context, userID uint64) error {
	return r.q.RevokeAllRefreshTokensByUser(ctx, userID)
}

func nullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}
