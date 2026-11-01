-- name: ListJamu :many
-- Termasuk data produk sachet: nama, stok (peringatan stok), dan HPP
-- (HPP per gelas = HPP sachet karena satu gelas = satu sachet).
SELECT
    j.id, j.nama, j.harga_jual, j.aktif, j.created_at,
    j.produk_id, p.nama AS produk_nama, p.stok AS produk_stok, p.hpp_rata_rata AS produk_hpp
FROM jamu j
JOIN produk p ON p.id = j.produk_id
WHERE (sqlc.narg(aktif) IS NULL OR j.aktif = sqlc.narg(aktif))
ORDER BY j.nama, j.id;

-- name: GetJamuByID :one
SELECT
    j.id, j.nama, j.harga_jual, j.aktif, j.created_at,
    j.produk_id, p.nama AS produk_nama, p.stok AS produk_stok, p.hpp_rata_rata AS produk_hpp
FROM jamu j
JOIN produk p ON p.id = j.produk_id
WHERE j.id = ?;

-- name: CreateJamu :execlastid
INSERT INTO jamu (nama, harga_jual, produk_id) VALUES (?, ?, ?);

-- name: UpdateJamu :exec
UPDATE jamu SET nama = ?, harga_jual = ?, produk_id = ?, aktif = ? WHERE id = ?;
