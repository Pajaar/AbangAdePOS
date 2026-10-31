-- ============================================================
-- Data awal akses POS Jamu: role, menu, endpoint API, hak akses
-- Jalankan SETELAH 000001_init_schema.up.sql.
--
-- Path endpoint memakai gaya route Gin (':id').
-- Middleware harus mencocokkan POLA route (di Gin: c.FullPath()),
-- bukan URL asli yang berisi angka.
--
-- Endpoint berikut sengaja TIDAK didaftarkan karena bebas untuk semua user yang login:
--   POST /api/v1/auth/login   (publik)
--   GET  /api/v1/auth/me      (data user yang sedang login)
--   GET  /api/v1/auth/menus   (menu sidebar sesuai role)
--
-- User headadmin pertama TIDAK dibuat di sini, karena password harus di-hash bcrypt.
-- Buat lewat perintah Go terpisah: cmd/seed-admin.
-- ============================================================

SET NAMES utf8mb4 COLLATE utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- ROLE
-- ------------------------------------------------------------
INSERT INTO roles (kode, nama) VALUES
  ('headadmin', 'Head Admin'),
  ('admin',     'Admin'),
  ('cashier',   'Kasir');

-- ------------------------------------------------------------
-- MENU (sidebar) DAN AKSI
-- ------------------------------------------------------------
INSERT INTO menus (kode, nama, tipe, path, icon, urutan) VALUES
  ('kasir',       'Kasir',                 'menu', '/kasir',       'cash',     1),
  ('piutang',     'Piutang',               'menu', '/piutang',     'wallet',   2),
  ('produk',      'Produk',                'menu', '/produk',      'box',      3),
  ('menu_jamu',   'Menu Jamu',             'menu', '/menu-jamu',   'cup',      4),
  ('pembelian',   'Pembelian',             'menu', '/pembelian',   'truck',    5),
  ('pengeluaran', 'Pengeluaran',           'menu', '/pengeluaran', 'receipt',  6),
  ('laporan',     'Laporan',               'menu', '/laporan',     'chart',    7),
  ('pengguna',    'Pengguna',              'menu', '/pengguna',    'users',    8),
  ('lihat_hpp',   'Lihat HPP dan margin',  'aksi', NULL,           NULL,      90);

-- ------------------------------------------------------------
-- ENDPOINT API
-- ------------------------------------------------------------
INSERT INTO api_endpoints (method, path, keterangan) VALUES
  -- kasir
  ('GET',    '/api/v1/katalog',                     'Grid produk, jamu seduh, dan tambahan untuk kasir'),
  ('GET',    '/api/v1/barcode/:kode',               'Cari produk dari barcode'),
  ('POST',   '/api/v1/penjualan',                   'Simpan transaksi penjualan'),
  ('GET',    '/api/v1/penjualan/:id',               'Detail nota / cetak struk'),
  ('GET',    '/api/v1/pelanggan',                   'Cari nama pelanggan (autocomplete)'),
  ('POST',   '/api/v1/pelanggan',                   'Tambah nama pelanggan baru'),
  -- piutang
  ('GET',    '/api/v1/piutang',                     'Daftar piutang per pelanggan'),
  ('GET',    '/api/v1/piutang/:pelanggan_id',       'Kartu piutang satu pelanggan'),
  ('POST',   '/api/v1/piutang/pembayaran',          'Terima pembayaran / cicilan piutang'),
  -- produk
  ('GET',    '/api/v1/produk',                      'Daftar produk'),
  ('GET',    '/api/v1/produk/:id',                  'Detail produk'),
  ('POST',   '/api/v1/produk',                      'Tambah produk'),
  ('PUT',    '/api/v1/produk/:id',                  'Ubah produk'),
  ('POST',   '/api/v1/produk/:id/barcode',          'Tambah barcode ke produk'),
  ('DELETE', '/api/v1/produk-barcode/:id',          'Hapus barcode dari produk'),
  ('POST',   '/api/v1/barcode/lookup',              'Cari nama produk di API publik'),
  ('GET',    '/api/v1/kategori',                    'Daftar kategori'),
  ('POST',   '/api/v1/kategori',                    'Tambah kategori'),
  ('PUT',    '/api/v1/kategori/:id',                'Ubah kategori'),
  -- menu jamu
  ('GET',    '/api/v1/jamu',                        'Daftar jamu seduh'),
  ('POST',   '/api/v1/jamu',                        'Tambah jamu seduh'),
  ('PUT',    '/api/v1/jamu/:id',                    'Ubah jamu seduh'),
  ('GET',    '/api/v1/additional-items',            'Daftar tambahan'),
  ('POST',   '/api/v1/additional-items',            'Tambah item tambahan'),
  ('PUT',    '/api/v1/additional-items/:id',        'Ubah item tambahan'),
  -- pembelian
  ('GET',    '/api/v1/supplier',                    'Daftar supplier'),
  ('POST',   '/api/v1/supplier',                    'Tambah supplier'),
  ('GET',    '/api/v1/pembelian',                   'Daftar pembelian'),
  ('GET',    '/api/v1/pembelian/:id',               'Detail pembelian'),
  ('POST',   '/api/v1/pembelian',                   'Simpan pembelian stok'),
  ('POST',   '/api/v1/pembelian/:id/pembayaran',    'Bayar utang ke supplier'),
  -- pengeluaran
  ('GET',    '/api/v1/pengeluaran',                 'Daftar pengeluaran'),
  ('POST',   '/api/v1/pengeluaran',                 'Catat pengeluaran'),
  ('GET',    '/api/v1/akun',                        'Daftar akun'),
  -- laporan
  ('GET',    '/api/v1/laporan/laba-rugi',           'Laporan laba rugi'),
  ('GET',    '/api/v1/laporan/margin',              'Margin per produk'),
  ('GET',    '/api/v1/laporan/diskon',              'Total diskon dan potongan tawar'),
  ('GET',    '/api/v1/laporan/buku-besar',          'Buku besar per akun'),
  -- pengguna
  ('GET',    '/api/v1/users',                       'Daftar pengguna'),
  ('POST',   '/api/v1/users',                       'Tambah pengguna'),
  ('PUT',    '/api/v1/users/:id',                   'Ubah pengguna / reset password'),
  ('GET',    '/api/v1/roles',                       'Daftar role'),
  ('GET',    '/api/v1/roles/:id/akses',             'Lihat hak akses role'),
  ('PUT',    '/api/v1/roles/:id/akses',             'Ubah hak akses role'),
  ('GET',    '/api/v1/menus',                       'Daftar semua menu dan aksi');

-- ------------------------------------------------------------
-- ENDPOINT MILIK MENU
-- (satu endpoint boleh muncul di beberapa menu)
-- ------------------------------------------------------------
INSERT INTO menu_endpoints (menu_id, endpoint_id)
SELECT m.id, e.id
FROM (
            SELECT 'kasir' AS menu, 'GET' AS method, '/api/v1/katalog' AS path
  UNION ALL SELECT 'kasir',       'GET',    '/api/v1/barcode/:kode'
  UNION ALL SELECT 'kasir',       'POST',   '/api/v1/penjualan'
  UNION ALL SELECT 'kasir',       'GET',    '/api/v1/penjualan/:id'
  UNION ALL SELECT 'kasir',       'GET',    '/api/v1/pelanggan'
  UNION ALL SELECT 'kasir',       'POST',   '/api/v1/pelanggan'

  UNION ALL SELECT 'piutang',     'GET',    '/api/v1/piutang'
  UNION ALL SELECT 'piutang',     'GET',    '/api/v1/piutang/:pelanggan_id'
  UNION ALL SELECT 'piutang',     'POST',   '/api/v1/piutang/pembayaran'
  UNION ALL SELECT 'piutang',     'GET',    '/api/v1/pelanggan'
  UNION ALL SELECT 'piutang',     'GET',    '/api/v1/penjualan/:id'

  UNION ALL SELECT 'produk',      'GET',    '/api/v1/produk'
  UNION ALL SELECT 'produk',      'GET',    '/api/v1/produk/:id'
  UNION ALL SELECT 'produk',      'POST',   '/api/v1/produk'
  UNION ALL SELECT 'produk',      'PUT',    '/api/v1/produk/:id'
  UNION ALL SELECT 'produk',      'POST',   '/api/v1/produk/:id/barcode'
  UNION ALL SELECT 'produk',      'DELETE', '/api/v1/produk-barcode/:id'
  UNION ALL SELECT 'produk',      'GET',    '/api/v1/barcode/:kode'
  UNION ALL SELECT 'produk',      'POST',   '/api/v1/barcode/lookup'
  UNION ALL SELECT 'produk',      'GET',    '/api/v1/kategori'
  UNION ALL SELECT 'produk',      'POST',   '/api/v1/kategori'
  UNION ALL SELECT 'produk',      'PUT',    '/api/v1/kategori/:id'

  UNION ALL SELECT 'menu_jamu',   'GET',    '/api/v1/jamu'
  UNION ALL SELECT 'menu_jamu',   'POST',   '/api/v1/jamu'
  UNION ALL SELECT 'menu_jamu',   'PUT',    '/api/v1/jamu/:id'
  UNION ALL SELECT 'menu_jamu',   'GET',    '/api/v1/additional-items'
  UNION ALL SELECT 'menu_jamu',   'POST',   '/api/v1/additional-items'
  UNION ALL SELECT 'menu_jamu',   'PUT',    '/api/v1/additional-items/:id'
  UNION ALL SELECT 'menu_jamu',   'GET',    '/api/v1/produk'

  UNION ALL SELECT 'pembelian',   'GET',    '/api/v1/supplier'
  UNION ALL SELECT 'pembelian',   'POST',   '/api/v1/supplier'
  UNION ALL SELECT 'pembelian',   'GET',    '/api/v1/pembelian'
  UNION ALL SELECT 'pembelian',   'GET',    '/api/v1/pembelian/:id'
  UNION ALL SELECT 'pembelian',   'POST',   '/api/v1/pembelian'
  UNION ALL SELECT 'pembelian',   'POST',   '/api/v1/pembelian/:id/pembayaran'
  UNION ALL SELECT 'pembelian',   'GET',    '/api/v1/produk'
  UNION ALL SELECT 'pembelian',   'GET',    '/api/v1/barcode/:kode'
  UNION ALL SELECT 'pembelian',   'POST',   '/api/v1/barcode/lookup'

  UNION ALL SELECT 'pengeluaran', 'GET',    '/api/v1/pengeluaran'
  UNION ALL SELECT 'pengeluaran', 'POST',   '/api/v1/pengeluaran'
  UNION ALL SELECT 'pengeluaran', 'GET',    '/api/v1/akun'

  UNION ALL SELECT 'laporan',     'GET',    '/api/v1/laporan/laba-rugi'
  UNION ALL SELECT 'laporan',     'GET',    '/api/v1/laporan/margin'
  UNION ALL SELECT 'laporan',     'GET',    '/api/v1/laporan/diskon'
  UNION ALL SELECT 'laporan',     'GET',    '/api/v1/laporan/buku-besar'
  UNION ALL SELECT 'laporan',     'GET',    '/api/v1/akun'

  UNION ALL SELECT 'pengguna',    'GET',    '/api/v1/users'
  UNION ALL SELECT 'pengguna',    'POST',   '/api/v1/users'
  UNION ALL SELECT 'pengguna',    'PUT',    '/api/v1/users/:id'
  UNION ALL SELECT 'pengguna',    'GET',    '/api/v1/roles'
  UNION ALL SELECT 'pengguna',    'GET',    '/api/v1/roles/:id/akses'
  UNION ALL SELECT 'pengguna',    'PUT',    '/api/v1/roles/:id/akses'
  UNION ALL SELECT 'pengguna',    'GET',    '/api/v1/menus'
) x
JOIN menus m         ON m.kode = x.menu
JOIN api_endpoints e ON e.method = x.method AND e.path = x.path;

-- ------------------------------------------------------------
-- HAK AKSES ROLE
-- ------------------------------------------------------------

-- headadmin: semua menu dan aksi, diberikan eksplisit (bukan bypass di kode)
INSERT INTO role_menus (role_id, menu_id)
SELECT r.id, m.id
FROM roles r CROSS JOIN menus m
WHERE r.kode = 'headadmin';

-- admin: fokus input barang. Tanpa laporan, tanpa lihat_hpp.
INSERT INTO role_menus (role_id, menu_id)
SELECT r.id, m.id
FROM roles r JOIN menus m ON m.kode IN ('produk', 'menu_jamu', 'pembelian')
WHERE r.kode = 'admin';

-- cashier: kasir dan terima pembayaran piutang di meja kasir.
INSERT INTO role_menus (role_id, menu_id)
SELECT r.id, m.id
FROM roles r JOIN menus m ON m.kode IN ('kasir', 'piutang')
WHERE r.kode = 'cashier';
