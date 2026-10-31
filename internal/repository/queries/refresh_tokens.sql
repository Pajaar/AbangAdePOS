-- name: CreateRefreshToken :execlastid
INSERT INTO refresh_tokens (user_id, token_hash, expires_at, user_agent, ip_address)
VALUES (?, ?, ?, ?, ?);

-- name: GetRefreshTokenByHashForUpdate :one
-- Wajib dipanggil di dalam transaksi. Baris dikunci sampai commit supaya
-- dua refresh dengan token yang sama tidak bisa lolos bersamaan.
SELECT id, user_id, token_hash, expires_at, revoked_at, replaced_by_id, user_agent, ip_address, created_at
FROM refresh_tokens
WHERE token_hash = ?
FOR UPDATE;

-- name: RevokeRefreshTokenByHash :exec
UPDATE refresh_tokens
SET revoked_at = NOW()
WHERE token_hash = ? AND revoked_at IS NULL;

-- name: RotateRefreshToken :execrows
-- Hanya berhasil (1 baris) bila token lama belum di-revoke.
UPDATE refresh_tokens
SET revoked_at = NOW(), replaced_by_id = ?
WHERE id = ? AND revoked_at IS NULL;

-- name: RevokeAllRefreshTokensByUser :exec
UPDATE refresh_tokens
SET revoked_at = NOW()
WHERE user_id = ? AND revoked_at IS NULL;
