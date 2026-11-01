-- name: ListAdditionalItems :many
-- produk_id kosong = tambahan yang tidak dilacak stoknya (telur, madu), HPP 0.
SELECT
    a.id, a.nama, a.harga_jual, a.produk_id, a.aktif,
    p.nama AS produk_nama, p.hpp_rata_rata AS produk_hpp
FROM additional_items a
LEFT JOIN produk p ON p.id = a.produk_id
WHERE (sqlc.narg(aktif) IS NULL OR a.aktif = sqlc.narg(aktif))
ORDER BY a.nama, a.id;

-- name: GetAdditionalItemByID :one
SELECT
    a.id, a.nama, a.harga_jual, a.produk_id, a.aktif,
    p.nama AS produk_nama, p.hpp_rata_rata AS produk_hpp
FROM additional_items a
LEFT JOIN produk p ON p.id = a.produk_id
WHERE a.id = ?;

-- name: CreateAdditionalItem :execlastid
INSERT INTO additional_items (nama, harga_jual, produk_id) VALUES (?, ?, ?);

-- name: UpdateAdditionalItem :exec
UPDATE additional_items SET nama = ?, harga_jual = ?, produk_id = ?, aktif = ? WHERE id = ?;
