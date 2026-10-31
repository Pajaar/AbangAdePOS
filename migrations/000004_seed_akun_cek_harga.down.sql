-- Membatalkan 000004_seed_akun_cek_harga.up.sql.
-- Penghapusan akun gagal (FK) bila sudah ada jurnal atau pengeluaran yang memakainya.
DELETE me FROM menu_endpoints me
JOIN api_endpoints e ON e.id = me.endpoint_id
WHERE e.method = 'POST' AND e.path = '/api/v1/kasir/cek-harga';

DELETE FROM api_endpoints WHERE method = 'POST' AND path = '/api/v1/kasir/cek-harga';

DELETE FROM akun WHERE kode IN
  ('1-101', '1-102', '1-103', '2-101', '3-101', '4-101', '4-102', '5-101', '6-101');
