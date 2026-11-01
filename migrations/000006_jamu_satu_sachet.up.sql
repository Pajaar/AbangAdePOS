-- ============================================================
-- Jamu seduh = tepat satu sachet dari satu produk.
--   1. jamu.produk_id menggantikan tabel jamu_sachet (yang dihapus).
--   2. Sachet tambahan di satu gelas boleh produk apa saja, harganya dari
--      produk.harga_normal. Karena itu baris tambahan di nota
--      (penjualan_detail_additional) berisi additional_item_id ATAU produk_id.
-- ============================================================

-- ------------------------------------------------------------
-- JAMU: satu produk sachet per jamu
-- ------------------------------------------------------------
ALTER TABLE jamu ADD COLUMN produk_id BIGINT UNSIGNED NULL AFTER harga_jual;

-- Pindahkan data lama: ambil sachet pertama tiap jamu.
UPDATE jamu j
JOIN jamu_sachet js ON js.id = (
  SELECT MIN(x.id) FROM jamu_sachet x WHERE x.jamu_id = j.id
)
SET j.produk_id = js.produk_id;

ALTER TABLE jamu MODIFY produk_id BIGINT UNSIGNED NOT NULL;
ALTER TABLE jamu ADD CONSTRAINT fk_jamu_produk FOREIGN KEY (produk_id) REFERENCES produk(id);

DROP TABLE jamu_sachet;

-- ------------------------------------------------------------
-- TAMBAHAN DI NOTA: item tambahan tetap ATAU sachet produk
-- ------------------------------------------------------------
--   additional_item_id terisi -> tambahan harga tetap (telur, madu)
--   produk_id terisi          -> sachet tambahan; harga = snapshot produk.harga_normal,
--                                stok produk itu dipotong, hpp = snapshot hpp_rata_rata
ALTER TABLE penjualan_detail_additional MODIFY additional_item_id BIGINT UNSIGNED NULL;
ALTER TABLE penjualan_detail_additional ADD COLUMN produk_id BIGINT UNSIGNED NULL AFTER additional_item_id;
ALTER TABLE penjualan_detail_additional
  ADD CONSTRAINT fk_pjadd_produk FOREIGN KEY (produk_id) REFERENCES produk(id);
ALTER TABLE penjualan_detail_additional
  ADD CONSTRAINT chk_pjadd_item CHECK (
    (additional_item_id IS NOT NULL AND produk_id IS NULL) OR
    (additional_item_id IS NULL AND produk_id IS NOT NULL)
  );
