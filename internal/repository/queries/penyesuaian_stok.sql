-- name: CreatePenyesuaianStok :execlastid
INSERT INTO penyesuaian_stok (user_id, produk_id, jenis, qty, stok_sebelum, hpp, keterangan)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetPenyesuaianStokByID :one
SELECT ps.id, ps.tanggal, ps.user_id, u.nama AS user_nama, ps.produk_id, ps.jenis,
       ps.qty, ps.stok_sebelum, ps.hpp, ps.keterangan
FROM penyesuaian_stok ps
JOIN users u ON u.id = ps.user_id
WHERE ps.id = ?;

-- name: ListPenyesuaianStokByProduk :many
SELECT ps.id, ps.tanggal, ps.user_id, u.nama AS user_nama, ps.produk_id, ps.jenis,
       ps.qty, ps.stok_sebelum, ps.hpp, ps.keterangan
FROM penyesuaian_stok ps
JOIN users u ON u.id = ps.user_id
WHERE ps.produk_id = ?
ORDER BY ps.tanggal DESC, ps.id DESC
LIMIT ?;
