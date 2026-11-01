-- Membatalkan 000006_jamu_satu_sachet.up.sql.
-- Baris tambahan nota yang berisi produk_id (sachet tambahan) membuat langkah
-- NOT NULL additional_item_id gagal; hapus/ubah baris itu dulu bila ada.

ALTER TABLE penjualan_detail_additional DROP CHECK chk_pjadd_item;
ALTER TABLE penjualan_detail_additional DROP FOREIGN KEY fk_pjadd_produk;
ALTER TABLE penjualan_detail_additional DROP COLUMN produk_id;
ALTER TABLE penjualan_detail_additional MODIFY additional_item_id BIGINT UNSIGNED NOT NULL;

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

INSERT INTO jamu_sachet (jamu_id, produk_id, jumlah)
SELECT id, produk_id, 1 FROM jamu;

ALTER TABLE jamu DROP FOREIGN KEY fk_jamu_produk;
ALTER TABLE jamu DROP COLUMN produk_id;
