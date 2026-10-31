-- name: ListKategori :many
SELECT id, nama, aktif
FROM kategori
WHERE (sqlc.narg(aktif) IS NULL OR aktif = sqlc.narg(aktif))
ORDER BY nama;

-- name: GetKategoriByID :one
SELECT id, nama, aktif FROM kategori WHERE id = ?;

-- name: GetKategoriByNama :one
-- Perbandingan mengikuti collation tabel (tidak membedakan huruf besar/kecil).
SELECT id, nama, aktif FROM kategori WHERE nama = ?;

-- name: CreateKategori :execlastid
INSERT INTO kategori (nama) VALUES (?);

-- name: UpdateKategori :exec
UPDATE kategori SET nama = ?, aktif = ? WHERE id = ?;
