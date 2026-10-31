-- ============================================================
-- 1. Daftar akun awal (usulan default, belum dikonfirmasi pemilik).
--    Kode dan nama boleh diubah lewat data; kode Go mencari akun lewat kode.
-- 2. Endpoint cek harga di kasir (peringatan di bawah modal, tanpa angka).
-- ============================================================

SET NAMES utf8mb4 COLLATE utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- AKUN
-- ------------------------------------------------------------
-- Potongan Penjualan bertipe pendapatan tapi saldo normalnya debit
-- (pengurang penjualan di laporan laba rugi).
INSERT INTO akun (kode, nama, tipe) VALUES
  ('1-101', 'Kas',                'aset'),
  ('1-102', 'Piutang Usaha',      'aset'),
  ('1-103', 'Persediaan',         'aset'),
  ('2-101', 'Utang Usaha',        'kewajiban'),
  ('3-101', 'Modal Pemilik',      'modal'),
  ('4-101', 'Penjualan',          'pendapatan'),
  ('4-102', 'Potongan Penjualan', 'pendapatan'),
  ('5-101', 'HPP',                'beban'),
  ('6-101', 'Beban Operasional',  'beban');

-- ------------------------------------------------------------
-- ENDPOINT CEK HARGA (menu Kasir)
-- ------------------------------------------------------------
INSERT INTO api_endpoints (method, path, keterangan) VALUES
  ('POST', '/api/v1/kasir/cek-harga', 'Cek harga keranjang di bawah modal (hanya boolean, tanpa angka)');

INSERT INTO menu_endpoints (menu_id, endpoint_id)
SELECT m.id, e.id
FROM menus m
JOIN api_endpoints e ON e.method = 'POST' AND e.path = '/api/v1/kasir/cek-harga'
WHERE m.kode = 'kasir';
