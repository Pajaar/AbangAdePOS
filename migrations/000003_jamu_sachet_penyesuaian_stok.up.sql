-- ============================================================
-- 1. Jamu seduh boleh memakai lebih dari satu jenis sachet.
--    jamu.sachet_id + jamu.sachet_per_gelas dipindah ke tabel jamu_sachet.
-- 2. Tabel penyesuaian_stok untuk stok awal, stock opname, dan barang rusak.
-- ============================================================

-- ------------------------------------------------------------
-- JAMU_SACHET
-- ------------------------------------------------------------

-- Komposisi sachet per gelas jamu seduh.
-- Saat jamu terjual, stok setiap produk_id dipotong sebanyak jumlah x qty gelas.
-- HPP jamu per gelas = SUM(jumlah x produk.hpp_rata_rata).
-- Setiap jamu wajib punya minimal satu baris (dijaga di service).
CREATE TABLE jamu_sachet (
  id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  jamu_id    BIGINT UNSIGNED NOT NULL,
  produk_id  BIGINT UNSIGNED NOT NULL,
  jumlah     INT NOT NULL DEFAULT 1,
  CONSTRAINT fk_jamusachet_jamu   FOREIGN KEY (jamu_id)   REFERENCES jamu(id),
  CONSTRAINT fk_jamusachet_produk FOREIGN KEY (produk_id) REFERENCES produk(id),
  CONSTRAINT chk_jamusachet_jumlah CHECK (jumlah > 0),
  UNIQUE KEY uq_jamusachet (jamu_id, produk_id)
) ENGINE=InnoDB;

-- Pindahkan data lama: satu sachet per jamu.
INSERT INTO jamu_sachet (jamu_id, produk_id, jumlah)
SELECT id, sachet_id, sachet_per_gelas FROM jamu;

-- CHECK dan FK harus dilepas dulu sebelum kolomnya dihapus.
ALTER TABLE jamu DROP CHECK chk_jamu_sachet;
ALTER TABLE jamu DROP FOREIGN KEY fk_jamu_sachet;
ALTER TABLE jamu DROP COLUMN sachet_id, DROP COLUMN sachet_per_gelas;

-- ------------------------------------------------------------
-- PENYESUAIAN_STOK
-- ------------------------------------------------------------

-- Perubahan stok di luar penjualan dan pembelian.
--   stok_awal : stok barang lama saat sistem mulai dipakai. qty > 0.
--               hpp = harga modal per unit yang diinput; ikut rumus moving average.
--               Jurnal: Persediaan / Modal Pemilik.
--   opname    : koreksi hasil hitung fisik. qty = stok fisik - stok_sebelum (boleh + atau -).
--               Jurnal: stok berkurang HPP / Persediaan, stok bertambah Persediaan / HPP.
--   rusak     : barang rusak / kedaluwarsa dibuang. qty < 0.
--               Jurnal: HPP / Persediaan.
-- hpp          = nilai per unit saat penyesuaian (snapshot). Untuk opname dan rusak
--                = produk.hpp_rata_rata saat itu. Nilai jurnal = ABS(qty) x hpp.
-- stok_sebelum = produk.stok sesaat sebelum penyesuaian (untuk jejak audit).
CREATE TABLE penyesuaian_stok (
  id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  tanggal       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  user_id       BIGINT UNSIGNED NOT NULL,
  produk_id     BIGINT UNSIGNED NOT NULL,
  jenis         ENUM('stok_awal','opname','rusak') NOT NULL,
  qty           INT NOT NULL,
  stok_sebelum  INT NOT NULL,
  hpp           DECIMAL(14,2) NOT NULL DEFAULT 0,
  keterangan    VARCHAR(255) NULL,
  CONSTRAINT fk_penyesuaian_user   FOREIGN KEY (user_id)   REFERENCES users(id),
  CONSTRAINT fk_penyesuaian_produk FOREIGN KEY (produk_id) REFERENCES produk(id),
  CONSTRAINT chk_penyesuaian_hpp CHECK (hpp >= 0),
  CONSTRAINT chk_penyesuaian_qty CHECK (
    (jenis = 'stok_awal' AND qty > 0) OR
    (jenis = 'opname'    AND qty <> 0) OR
    (jenis = 'rusak'     AND qty < 0)
  ),
  INDEX idx_penyesuaian_produk (produk_id),
  INDEX idx_penyesuaian_tanggal (tanggal)
) ENGINE=InnoDB;
