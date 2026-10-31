-- Membatalkan 000003_jamu_sachet_penyesuaian_stok.up.sql.
-- Data jamu_sachet dikembalikan ke jamu.sachet_id: jika satu jamu punya lebih
-- dari satu sachet, hanya sachet dengan id jamu_sachet terkecil yang dipakai
-- (komposisi lainnya hilang). Jamu tanpa sachet membuat langkah NOT NULL gagal.

DROP TABLE IF EXISTS penyesuaian_stok;

ALTER TABLE jamu
  ADD COLUMN sachet_id        BIGINT UNSIGNED NULL AFTER harga_jual,
  ADD COLUMN sachet_per_gelas INT NOT NULL DEFAULT 1 AFTER sachet_id;

UPDATE jamu j
JOIN jamu_sachet js ON js.id = (
  SELECT MIN(x.id) FROM jamu_sachet x WHERE x.jamu_id = j.id
)
SET j.sachet_id = js.produk_id,
    j.sachet_per_gelas = js.jumlah;

ALTER TABLE jamu MODIFY sachet_id BIGINT UNSIGNED NOT NULL;
ALTER TABLE jamu ADD CONSTRAINT fk_jamu_sachet FOREIGN KEY (sachet_id) REFERENCES produk(id);
ALTER TABLE jamu ADD CONSTRAINT chk_jamu_sachet CHECK (sachet_per_gelas > 0);

DROP TABLE IF EXISTS jamu_sachet;
