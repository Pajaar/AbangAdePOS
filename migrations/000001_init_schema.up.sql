-- ============================================================
-- Skema database POS Jamu (versi 3)
-- Target: MySQL 8.0.16+ (InnoDB, CHECK constraint aktif)
-- Urutan tabel sudah disusun agar foreign key bisa dibuat langsung.
--
-- Perubahan dari versi 2:
--   - Akses dan role: roles, users, refresh_tokens, menus, api_endpoints,
--     menu_endpoints, role_menus
--   - Barcode: produk_barcode, barcode_lookup_cache, produk.status
--   - Kolom user_id (siapa yang mencatat) di penjualan, pembayaran_piutang,
--     pembelian, pembayaran_utang, jurnal, pengeluaran
--
-- Data awal role, menu, dan endpoint ada di migrasi 000002_seed_akses.up.sql
--
-- Catatan migrasi: CREATE DATABASE/USE sengaja TIDAK ada di sini.
-- DSN yang dipakai migrate CLI maupun aplikasi sudah menunjuk ke database
-- target (DB_NAME), apapun nama yang disediakan oleh provider (Railway, dll).
-- ============================================================

-- ------------------------------------------------------------
-- AKSES DAN ROLE
-- ------------------------------------------------------------

-- headadmin, admin, cashier
CREATE TABLE roles (
  id    BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  kode  VARCHAR(30)  NOT NULL UNIQUE,
  nama  VARCHAR(100) NOT NULL
) ENGINE=InnoDB;

CREATE TABLE users (
  id             BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  nama           VARCHAR(100) NOT NULL,
  username       VARCHAR(50)  NOT NULL UNIQUE,
  password_hash  VARCHAR(255) NOT NULL,          -- bcrypt, jangan simpan password asli
  role_id        BIGINT UNSIGNED NOT NULL,
  aktif          BOOLEAN NOT NULL DEFAULT TRUE,
  last_login_at  DATETIME NULL,
  created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_users_role FOREIGN KEY (role_id) REFERENCES roles(id)
) ENGINE=InnoDB;

-- Refresh token untuk memperpanjang sesi login.
-- Access token (JWT) TIDAK disimpan di database: umurnya pendek dan diverifikasi lewat tanda tangan.
-- token_hash = SHA-256 dari refresh token. Token aslinya hanya ada di cookie user.
-- Rotasi: setiap refresh, token lama di-revoke dan replaced_by_id menunjuk token baru.
-- Jika token yang sudah di-revoke dipakai lagi, anggap dicuri: revoke semua token user itu.
CREATE TABLE refresh_tokens (
  id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  user_id         BIGINT UNSIGNED NOT NULL,
  token_hash      CHAR(64) NOT NULL UNIQUE,
  expires_at      DATETIME NOT NULL,
  revoked_at      DATETIME NULL,
  replaced_by_id  BIGINT UNSIGNED NULL,
  user_agent      VARCHAR(255) NULL,
  ip_address      VARCHAR(45)  NULL,
  created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_refresh_user     FOREIGN KEY (user_id)        REFERENCES users(id),
  CONSTRAINT fk_refresh_replaced FOREIGN KEY (replaced_by_id) REFERENCES refresh_tokens(id),
  INDEX idx_refresh_user (user_id)
) ENGINE=InnoDB;

-- tipe 'menu' : tampil di sidebar frontend (path = route frontend)
-- tipe 'aksi' : izin tanpa halaman sendiri, misalnya lihat_hpp
--               (menentukan apakah kolom HPP dan margin ikut dikirim di response)
CREATE TABLE menus (
  id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  parent_id  BIGINT UNSIGNED NULL,
  kode       VARCHAR(50)  NOT NULL UNIQUE,
  nama       VARCHAR(100) NOT NULL,
  tipe       ENUM('menu','aksi') NOT NULL DEFAULT 'menu',
  path       VARCHAR(100) NULL,
  icon       VARCHAR(50)  NULL,
  urutan     INT NOT NULL DEFAULT 0,
  CONSTRAINT fk_menus_parent FOREIGN KEY (parent_id) REFERENCES menus(id)
) ENGINE=InnoDB;

-- Satu baris = satu route API.
-- path disimpan sebagai POLA route, bukan URL asli:
--   '/api/v1/produk/:id'  (bukan '/api/v1/produk/15')
-- Endpoint login dan /auth/* TIDAK didaftarkan di sini (bebas untuk user yang login).
CREATE TABLE api_endpoints (
  id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  method      ENUM('GET','POST','PUT','PATCH','DELETE') NOT NULL,
  path        VARCHAR(150) NOT NULL,
  keterangan  VARCHAR(255) NULL,
  UNIQUE KEY uq_endpoint (method, path)
) ENGINE=InnoDB;

-- Satu endpoint bisa dipakai beberapa menu.
-- Contoh: GET /api/v1/produk dipakai menu Produk dan menu Pembelian.
CREATE TABLE menu_endpoints (
  menu_id      BIGINT UNSIGNED NOT NULL,
  endpoint_id  BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (menu_id, endpoint_id),
  CONSTRAINT fk_menuep_menu     FOREIGN KEY (menu_id)     REFERENCES menus(id),
  CONSTRAINT fk_menuep_endpoint FOREIGN KEY (endpoint_id) REFERENCES api_endpoints(id)
) ENGINE=InnoDB;

-- Hak akses: role boleh membuka menu apa saja.
-- Endpoint yang boleh dipanggil = semua endpoint milik menu-menu tersebut.
CREATE TABLE role_menus (
  role_id  BIGINT UNSIGNED NOT NULL,
  menu_id  BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (role_id, menu_id),
  CONSTRAINT fk_rolemenu_role FOREIGN KEY (role_id) REFERENCES roles(id),
  CONSTRAINT fk_rolemenu_menu FOREIGN KEY (menu_id) REFERENCES menus(id)
) ENGINE=InnoDB;

-- ------------------------------------------------------------
-- MASTER
-- ------------------------------------------------------------

-- Hanya diisi untuk orang yang ngutang. Pembeli tunai tidak dicatat.
CREATE TABLE pelanggan (
  id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  nama        VARCHAR(100) NOT NULL,
  no_hp       VARCHAR(20)  NULL,
  aktif       BOOLEAN NOT NULL DEFAULT TRUE,
  created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

CREATE TABLE supplier (
  id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  nama        VARCHAR(100) NOT NULL,
  no_hp       VARCHAR(20)  NULL,
  created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- Murni pengelompokan barang (Obat, Sachet Jamu, Minyak Urut, dll).
CREATE TABLE kategori (
  id     BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  nama   VARCHAR(100) NOT NULL UNIQUE,
  aktif  BOOLEAN NOT NULL DEFAULT TRUE
) ENGINE=InnoDB;

-- Semua barang yang distok. Stok SELALU dipotong saat terjual.
-- Sachet jamu juga ada di sini: harga_normal = harga jual sachet utuh.
-- Harga tengah dan nett tidak disimpan: dihitung dari hpp_rata_rata dan margin_min_persen.
-- status 'draft' = produk hasil scan barcode yang belum dilengkapi, tidak tampil di kasir.
CREATE TABLE produk (
  id                 BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  kategori_id        BIGINT UNSIGNED NOT NULL,
  nama               VARCHAR(150) NOT NULL,
  harga_normal       DECIMAL(14,2) NOT NULL DEFAULT 0,
  bisa_ditawar       BOOLEAN NOT NULL DEFAULT FALSE,
  margin_min_persen  DECIMAL(5,2) NULL,
  stok               INT NOT NULL DEFAULT 0,
  hpp_rata_rata      DECIMAL(14,2) NOT NULL DEFAULT 0,
  aturan_pakai       VARCHAR(255) NULL,   -- contoh: "3x sehari 1 kapsul, sebelum makan"
  status             ENUM('aktif','draft') NOT NULL DEFAULT 'aktif',
  aktif              BOOLEAN NOT NULL DEFAULT TRUE,
  created_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_produk_kategori FOREIGN KEY (kategori_id) REFERENCES kategori(id),
  CONSTRAINT chk_produk_harga CHECK (harga_normal >= 0)
) ENGINE=InnoDB;

-- Satu produk bisa punya beberapa barcode (per sachet, per kotak, kemasan lama/baru).
-- barcode disimpan sudah dinormalisasi (UPC-A 12 digit -> EAN-13 dengan awalan 0).
-- isi_per_kemasan: scan barcode kotak isi 10 = stok berkurang 10.
CREATE TABLE produk_barcode (
  id               BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  produk_id        BIGINT UNSIGNED NOT NULL,
  barcode          VARCHAR(14) NOT NULL UNIQUE,
  jenis            ENUM('EAN13','EAN8','UPCA','INTERNAL') NOT NULL,
  isi_per_kemasan  INT NOT NULL DEFAULT 1,
  sumber           ENUM('manual','openfoodfacts','upcitemdb') NOT NULL DEFAULT 'manual',
  created_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_barcode_produk FOREIGN KEY (produk_id) REFERENCES produk(id),
  CONSTRAINT chk_barcode_isi CHECK (isi_per_kemasan > 0)
) ENGINE=InnoDB;

-- Cache hasil API publik, termasuk yang TIDAK ditemukan.
-- Error/timeout API jangan disimpan di sini.
CREATE TABLE barcode_lookup_cache (
  barcode     VARCHAR(14) PRIMARY KEY,
  sumber      VARCHAR(30) NOT NULL,
  ditemukan   BOOLEAN NOT NULL,
  nama        VARCHAR(255) NULL,
  merek       VARCHAR(150) NULL,
  raw_json    JSON NULL,
  fetched_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- Menu jamu seduh per gelas. Tidak punya stok sendiri:
-- saat terjual, stok sachet_id dipotong sebanyak sachet_per_gelas.
-- harga_jual di sini = harga seduh, terpisah dari harga sachet utuh di produk.
CREATE TABLE jamu (
  id                BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  nama              VARCHAR(100) NOT NULL,
  harga_jual        DECIMAL(14,2) NOT NULL,
  sachet_id         BIGINT UNSIGNED NOT NULL,
  sachet_per_gelas  INT NOT NULL DEFAULT 1,
  aktif             BOOLEAN NOT NULL DEFAULT TRUE,
  created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_jamu_sachet FOREIGN KEY (sachet_id) REFERENCES produk(id),
  CONSTRAINT chk_jamu_harga CHECK (harga_jual >= 0),
  CONSTRAINT chk_jamu_sachet CHECK (sachet_per_gelas > 0)
) ENGINE=InnoDB;

-- Tambahan untuk jamu, harga fix.
-- produk_id diisi jika tambahannya berupa sachet (stok dipotong saat terjual),
-- dikosongkan jika tidak dilacak stoknya (madu, telur, dll).
CREATE TABLE additional_items (
  id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  nama        VARCHAR(100) NOT NULL,
  harga_jual  DECIMAL(14,2) NOT NULL,
  produk_id   BIGINT UNSIGNED NULL,
  aktif       BOOLEAN NOT NULL DEFAULT TRUE,
  CONSTRAINT fk_additional_produk FOREIGN KEY (produk_id) REFERENCES produk(id),
  CONSTRAINT chk_additional_harga CHECK (harga_jual >= 0)
) ENGINE=InnoDB;

-- ------------------------------------------------------------
-- PENJUALAN DAN PIUTANG
-- ------------------------------------------------------------

-- total        = jumlah semua baris (termasuk additional) dikurangi diskon
-- dibayar      = uang yang dibayar saat transaksi
-- sisa_piutang = total - dibayar, lalu dikurangi setiap ada pembayaran_piutang
-- user_id      = kasir yang mencatat
CREATE TABLE penjualan (
  id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  no_nota       VARCHAR(30) NOT NULL UNIQUE,
  tanggal       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  user_id       BIGINT UNSIGNED NOT NULL,
  pelanggan_id  BIGINT UNSIGNED NULL,
  diskon        DECIMAL(14,2) NOT NULL DEFAULT 0,
  total         DECIMAL(14,2) NOT NULL,
  dibayar       DECIMAL(14,2) NOT NULL,
  sisa_piutang  DECIMAL(14,2) NOT NULL DEFAULT 0,
  CONSTRAINT fk_penjualan_user      FOREIGN KEY (user_id)      REFERENCES users(id),
  CONSTRAINT fk_penjualan_pelanggan FOREIGN KEY (pelanggan_id) REFERENCES pelanggan(id),
  CONSTRAINT chk_penjualan_angka CHECK (diskon >= 0 AND total >= 0 AND dibayar >= 0 AND sisa_piutang >= 0),
  -- Utang wajib punya nama: tidak boleh ada sisa piutang tanpa pelanggan.
  CONSTRAINT chk_penjualan_utang CHECK (sisa_piutang = 0 OR pelanggan_id IS NOT NULL),
  INDEX idx_penjualan_tanggal (tanggal),
  INDEX idx_penjualan_piutang (pelanggan_id, sisa_piutang)
) ENGINE=InnoDB;

-- Satu baris berisi produk ATAU jamu, tidak keduanya.
--   produk_id terisi -> potong stok produk itu
--   jamu_id terisi   -> potong stok sachet milik jamu itu
-- harga_normal, harga_akhir, dan hpp adalah snapshot per unit pada saat transaksi.
CREATE TABLE penjualan_detail (
  id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  penjualan_id  BIGINT UNSIGNED NOT NULL,
  produk_id     BIGINT UNSIGNED NULL,
  jamu_id       BIGINT UNSIGNED NULL,
  qty           INT NOT NULL,
  harga_normal  DECIMAL(14,2) NOT NULL,
  harga_akhir   DECIMAL(14,2) NOT NULL,
  hpp           DECIMAL(14,2) NOT NULL DEFAULT 0,
  CONSTRAINT fk_pjdetail_penjualan FOREIGN KEY (penjualan_id) REFERENCES penjualan(id),
  CONSTRAINT fk_pjdetail_produk    FOREIGN KEY (produk_id)    REFERENCES produk(id),
  CONSTRAINT fk_pjdetail_jamu      FOREIGN KEY (jamu_id)      REFERENCES jamu(id),
  CONSTRAINT chk_pjdetail_qty CHECK (qty > 0),
  CONSTRAINT chk_pjdetail_item CHECK (
    (produk_id IS NOT NULL AND jamu_id IS NULL) OR
    (produk_id IS NULL AND jamu_id IS NOT NULL)
  )
) ENGINE=InnoDB;

CREATE TABLE penjualan_detail_additional (
  id                   BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  penjualan_detail_id  BIGINT UNSIGNED NOT NULL,
  additional_item_id   BIGINT UNSIGNED NOT NULL,
  qty                  INT NOT NULL DEFAULT 1,
  harga                DECIMAL(14,2) NOT NULL,
  hpp                  DECIMAL(14,2) NOT NULL DEFAULT 0,
  CONSTRAINT fk_pjadd_detail FOREIGN KEY (penjualan_detail_id) REFERENCES penjualan_detail(id),
  CONSTRAINT fk_pjadd_item   FOREIGN KEY (additional_item_id)  REFERENCES additional_items(id),
  CONSTRAINT chk_pjadd_qty CHECK (qty > 0)
) ENGINE=InnoDB;

CREATE TABLE pembayaran_piutang (
  id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  penjualan_id  BIGINT UNSIGNED NOT NULL,
  user_id       BIGINT UNSIGNED NOT NULL,
  tanggal       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  jumlah        DECIMAL(14,2) NOT NULL,
  CONSTRAINT fk_bayarpiutang_penjualan FOREIGN KEY (penjualan_id) REFERENCES penjualan(id),
  CONSTRAINT fk_bayarpiutang_user      FOREIGN KEY (user_id)      REFERENCES users(id),
  CONSTRAINT chk_bayarpiutang_jumlah CHECK (jumlah > 0)
) ENGINE=InnoDB;

-- ------------------------------------------------------------
-- PEMBELIAN
-- ------------------------------------------------------------

CREATE TABLE pembelian (
  id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  supplier_id  BIGINT UNSIGNED NULL,
  user_id      BIGINT UNSIGNED NOT NULL,
  tanggal      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  total        DECIMAL(14,2) NOT NULL,
  dibayar      DECIMAL(14,2) NOT NULL,
  sisa_utang   DECIMAL(14,2) NOT NULL DEFAULT 0,
  CONSTRAINT fk_pembelian_supplier FOREIGN KEY (supplier_id) REFERENCES supplier(id),
  CONSTRAINT fk_pembelian_user     FOREIGN KEY (user_id)     REFERENCES users(id),
  CONSTRAINT chk_pembelian_angka CHECK (total >= 0 AND dibayar >= 0 AND sisa_utang >= 0),
  CONSTRAINT chk_pembelian_utang CHECK (sisa_utang = 0 OR supplier_id IS NOT NULL),
  INDEX idx_pembelian_tanggal (tanggal)
) ENGINE=InnoDB;

-- harga_beli per unit. Setiap baris baru memicu hitung ulang produk.hpp_rata_rata.
CREATE TABLE pembelian_detail (
  id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  pembelian_id  BIGINT UNSIGNED NOT NULL,
  produk_id     BIGINT UNSIGNED NOT NULL,
  qty           INT NOT NULL,
  harga_beli    DECIMAL(14,2) NOT NULL,
  CONSTRAINT fk_pbdetail_pembelian FOREIGN KEY (pembelian_id) REFERENCES pembelian(id),
  CONSTRAINT fk_pbdetail_produk    FOREIGN KEY (produk_id)    REFERENCES produk(id),
  CONSTRAINT chk_pbdetail_qty CHECK (qty > 0)
) ENGINE=InnoDB;

CREATE TABLE pembayaran_utang (
  id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  pembelian_id  BIGINT UNSIGNED NOT NULL,
  user_id       BIGINT UNSIGNED NOT NULL,
  tanggal       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  jumlah        DECIMAL(14,2) NOT NULL,
  CONSTRAINT fk_bayarutang_pembelian FOREIGN KEY (pembelian_id) REFERENCES pembelian(id),
  CONSTRAINT fk_bayarutang_user      FOREIGN KEY (user_id)      REFERENCES users(id),
  CONSTRAINT chk_bayarutang_jumlah CHECK (jumlah > 0)
) ENGINE=InnoDB;

-- ------------------------------------------------------------
-- AKUNTANSI
-- ------------------------------------------------------------

CREATE TABLE akun (
  id    BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  kode  VARCHAR(10)  NOT NULL UNIQUE,
  nama  VARCHAR(100) NOT NULL,
  tipe  ENUM('aset','kewajiban','modal','pendapatan','beban') NOT NULL
) ENGINE=InnoDB;

-- sumber_tipe + sumber_id menunjuk transaksi asal jurnal,
-- misalnya ('penjualan', 15) atau ('pembelian', 4).
CREATE TABLE jurnal (
  id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  tanggal      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  user_id      BIGINT UNSIGNED NOT NULL,
  keterangan   VARCHAR(255) NULL,
  sumber_tipe  VARCHAR(30) NULL,
  sumber_id    BIGINT UNSIGNED NULL,
  CONSTRAINT fk_jurnal_user FOREIGN KEY (user_id) REFERENCES users(id),
  INDEX idx_jurnal_tanggal (tanggal),
  INDEX idx_jurnal_sumber (sumber_tipe, sumber_id)
) ENGINE=InnoDB;

-- Satu baris hanya boleh berisi debit ATAU kredit.
-- Buku besar = query tabel ini per akun.
CREATE TABLE jurnal_detail (
  id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  jurnal_id  BIGINT UNSIGNED NOT NULL,
  akun_id    BIGINT UNSIGNED NOT NULL,
  debit      DECIMAL(14,2) NOT NULL DEFAULT 0,
  kredit     DECIMAL(14,2) NOT NULL DEFAULT 0,
  CONSTRAINT fk_jdetail_jurnal FOREIGN KEY (jurnal_id) REFERENCES jurnal(id),
  CONSTRAINT fk_jdetail_akun   FOREIGN KEY (akun_id)   REFERENCES akun(id),
  CONSTRAINT chk_jdetail_sisi CHECK (
    (debit > 0 AND kredit = 0) OR (kredit > 0 AND debit = 0)
  )
) ENGINE=InnoDB;

-- Biaya operasional dan bahan pelengkap (gula, gelas plastik, listrik, dll).
CREATE TABLE pengeluaran (
  id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  tanggal     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  user_id     BIGINT UNSIGNED NOT NULL,
  akun_id     BIGINT UNSIGNED NOT NULL,   -- akun beban yang dibebani
  keterangan  VARCHAR(255) NULL,
  jumlah      DECIMAL(14,2) NOT NULL,
  CONSTRAINT fk_pengeluaran_user FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT fk_pengeluaran_akun FOREIGN KEY (akun_id) REFERENCES akun(id),
  CONSTRAINT chk_pengeluaran_jumlah CHECK (jumlah > 0),
  INDEX idx_pengeluaran_tanggal (tanggal)
) ENGINE=InnoDB;
