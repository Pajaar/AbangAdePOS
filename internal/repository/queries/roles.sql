-- name: GetRoleByKode :one
SELECT id, kode, nama FROM roles WHERE kode = ?;
