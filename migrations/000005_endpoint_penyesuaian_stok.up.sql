-- ============================================================
-- Endpoint penyesuaian stok (stok awal, opname, rusak) masuk menu Produk,
-- sehingga admin dan headadmin sama-sama boleh memakainya.
-- ============================================================

SET NAMES utf8mb4 COLLATE utf8mb4_unicode_ci;

INSERT INTO api_endpoints (method, path, keterangan) VALUES
  ('POST', '/api/v1/produk/:id/penyesuaian-stok', 'Input stok awal, hasil stock opname, atau barang rusak'),
  ('GET',  '/api/v1/produk/:id/penyesuaian-stok', 'Riwayat penyesuaian stok satu produk');

INSERT INTO menu_endpoints (menu_id, endpoint_id)
SELECT m.id, e.id
FROM menus m
JOIN api_endpoints e ON e.path = '/api/v1/produk/:id/penyesuaian-stok'
WHERE m.kode = 'produk';
