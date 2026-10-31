-- Membatalkan 000001_init_schema.up.sql.
-- FOREIGN_KEY_CHECKS dimatikan sementara supaya urutan DROP tidak perlu
-- persis kebalikan dari urutan CREATE (lebih tahan terhadap perubahan skema).
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS pengeluaran;
DROP TABLE IF EXISTS jurnal_detail;
DROP TABLE IF EXISTS jurnal;
DROP TABLE IF EXISTS akun;
DROP TABLE IF EXISTS pembayaran_utang;
DROP TABLE IF EXISTS pembelian_detail;
DROP TABLE IF EXISTS pembelian;
DROP TABLE IF EXISTS pembayaran_piutang;
DROP TABLE IF EXISTS penjualan_detail_additional;
DROP TABLE IF EXISTS penjualan_detail;
DROP TABLE IF EXISTS penjualan;
DROP TABLE IF EXISTS additional_items;
DROP TABLE IF EXISTS jamu;
DROP TABLE IF EXISTS barcode_lookup_cache;
DROP TABLE IF EXISTS produk_barcode;
DROP TABLE IF EXISTS produk;
DROP TABLE IF EXISTS kategori;
DROP TABLE IF EXISTS supplier;
DROP TABLE IF EXISTS pelanggan;
DROP TABLE IF EXISTS role_menus;
DROP TABLE IF EXISTS menu_endpoints;
DROP TABLE IF EXISTS api_endpoints;
DROP TABLE IF EXISTS menus;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;

SET FOREIGN_KEY_CHECKS = 1;
