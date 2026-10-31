-- Membatalkan 000002_seed_akses.up.sql.
-- Urutan DELETE: child dulu baru parent, sesuai FK.
DELETE FROM role_menus;
DELETE FROM menu_endpoints;
DELETE FROM api_endpoints;
DELETE FROM menus;
DELETE FROM roles;
