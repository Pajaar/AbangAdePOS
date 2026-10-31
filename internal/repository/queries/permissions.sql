-- name: ListAllPermissions :many
-- Dipanggil sekali saat startup (dan saat reload) untuk membangun
-- map role_kode -> set("METHOD path") di internal/authorize.Store.
SELECT r.kode AS role_kode, e.method, e.path
FROM role_menus rm
JOIN menu_endpoints me ON me.menu_id = rm.menu_id
JOIN api_endpoints e ON e.id = me.endpoint_id
JOIN roles r ON r.id = rm.role_id;

-- name: ListAllRoleActions :many
-- Izin tanpa halaman (menus.tipe = 'aksi', misalnya lihat_hpp) per role.
SELECT r.kode AS role_kode, m.kode AS aksi_kode
FROM role_menus rm
JOIN menus m ON m.id = rm.menu_id
JOIN roles r ON r.id = rm.role_id
WHERE m.tipe = 'aksi';
