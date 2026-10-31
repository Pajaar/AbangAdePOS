-- name: ListBarcodeByProduk :many
SELECT id, produk_id, barcode, jenis, isi_per_kemasan, sumber, created_at
FROM produk_barcode
WHERE produk_id = ?
ORDER BY id;

-- name: GetBarcodeByID :one
SELECT id, produk_id, barcode, jenis, isi_per_kemasan, sumber, created_at
FROM produk_barcode
WHERE id = ?;

-- name: GetBarcodeByKode :one
SELECT id, produk_id, barcode, jenis, isi_per_kemasan, sumber, created_at
FROM produk_barcode
WHERE barcode = ?;

-- name: CreateBarcode :execlastid
INSERT INTO produk_barcode (produk_id, barcode, jenis, isi_per_kemasan, sumber)
VALUES (?, ?, ?, ?, ?);

-- name: DeleteBarcode :exec
DELETE FROM produk_barcode WHERE id = ?;

-- name: GetBarcodeLookupCache :one
SELECT barcode, sumber, ditemukan, nama, merek, fetched_at
FROM barcode_lookup_cache
WHERE barcode = ?;

-- name: UpsertBarcodeLookupCache :exec
INSERT INTO barcode_lookup_cache (barcode, sumber, ditemukan, nama, merek, raw_json)
VALUES (?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    sumber = VALUES(sumber), ditemukan = VALUES(ditemukan), nama = VALUES(nama),
    merek = VALUES(merek), raw_json = VALUES(raw_json), fetched_at = CURRENT_TIMESTAMP;
