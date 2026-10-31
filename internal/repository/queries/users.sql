-- name: GetUserByUsername :one
SELECT
    u.id,
    u.nama,
    u.username,
    u.password_hash,
    u.role_id,
    u.aktif,
    u.last_login_at,
    u.created_at,
    r.kode AS role_kode
FROM users u
JOIN roles r ON r.id = u.role_id
WHERE u.username = ?;

-- name: GetUserByID :one
SELECT
    u.id,
    u.nama,
    u.username,
    u.password_hash,
    u.role_id,
    u.aktif,
    u.last_login_at,
    u.created_at,
    r.kode AS role_kode
FROM users u
JOIN roles r ON r.id = u.role_id
WHERE u.id = ?;

-- name: UpdateLastLogin :exec
UPDATE users SET last_login_at = NOW() WHERE id = ?;

-- name: CreateUser :execlastid
INSERT INTO users (nama, username, password_hash, role_id)
VALUES (?, ?, ?, ?);
