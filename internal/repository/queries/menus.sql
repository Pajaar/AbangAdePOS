-- name: ListMenusByRole :many
-- Mengembalikan baris tipe='menu' DAN tipe='aksi' milik role tersebut.
-- Service yang memisahkan: 'menu' dibangun jadi tree, 'aksi' jadi daftar kode.
SELECT m.id, m.parent_id, m.kode, m.nama, m.tipe, m.path, m.icon, m.urutan
FROM role_menus rm
JOIN menus m ON m.id = rm.menu_id
JOIN roles r ON r.id = rm.role_id
WHERE r.kode = ?
ORDER BY m.urutan;
