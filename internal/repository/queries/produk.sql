-- name: ListProduk :many
-- q mencocokkan sebagian nama, atau barcode persis.
SELECT
    p.id, p.kategori_id, k.nama AS kategori_nama, p.nama, p.harga_normal,
    p.bisa_ditawar, p.margin_min_persen, p.stok, p.hpp_rata_rata,
    p.aturan_pakai, p.status, p.aktif, p.created_at
FROM produk p
JOIN kategori k ON k.id = p.kategori_id
WHERE (sqlc.narg(q) IS NULL
       OR p.nama LIKE CONCAT('%', sqlc.narg(q), '%')
       OR EXISTS (SELECT 1 FROM produk_barcode b WHERE b.produk_id = p.id AND b.barcode = sqlc.narg(q)))
  AND (sqlc.narg(kategori_id) IS NULL OR p.kategori_id = sqlc.narg(kategori_id))
  AND (sqlc.narg(status) IS NULL OR p.status = sqlc.narg(status))
  AND (sqlc.narg(aktif) IS NULL OR p.aktif = sqlc.narg(aktif))
ORDER BY p.nama, p.id
LIMIT ? OFFSET ?;

-- name: CountProduk :one
SELECT COUNT(*)
FROM produk p
WHERE (sqlc.narg(q) IS NULL
       OR p.nama LIKE CONCAT('%', sqlc.narg(q), '%')
       OR EXISTS (SELECT 1 FROM produk_barcode b WHERE b.produk_id = p.id AND b.barcode = sqlc.narg(q)))
  AND (sqlc.narg(kategori_id) IS NULL OR p.kategori_id = sqlc.narg(kategori_id))
  AND (sqlc.narg(status) IS NULL OR p.status = sqlc.narg(status))
  AND (sqlc.narg(aktif) IS NULL OR p.aktif = sqlc.narg(aktif));

-- name: GetProdukByID :one
SELECT
    p.id, p.kategori_id, k.nama AS kategori_nama, p.nama, p.harga_normal,
    p.bisa_ditawar, p.margin_min_persen, p.stok, p.hpp_rata_rata,
    p.aturan_pakai, p.status, p.aktif, p.created_at
FROM produk p
JOIN kategori k ON k.id = p.kategori_id
WHERE p.id = ?;

-- name: GetProdukStokForUpdate :one
-- Wajib di dalam transaksi: mengunci baris produk sebelum stok/HPP diubah.
SELECT id, stok, hpp_rata_rata FROM produk WHERE id = ? FOR UPDATE;

-- name: CreateProduk :execlastid
INSERT INTO produk (kategori_id, nama, harga_normal, bisa_ditawar, margin_min_persen, aturan_pakai, status)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: UpdateProduk :exec
-- Stok dan HPP tidak diubah di sini (hanya lewat pembelian, penjualan, penyesuaian stok).
UPDATE produk
SET kategori_id = ?, nama = ?, harga_normal = ?, bisa_ditawar = ?,
    margin_min_persen = ?, aturan_pakai = ?, status = ?, aktif = ?
WHERE id = ?;

-- name: UpdateProdukStokHPP :exec
UPDATE produk SET stok = ?, hpp_rata_rata = ? WHERE id = ?;
