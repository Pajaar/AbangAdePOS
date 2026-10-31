-- Membatalkan 000005_endpoint_penyesuaian_stok.up.sql.
DELETE me FROM menu_endpoints me
JOIN api_endpoints e ON e.id = me.endpoint_id
WHERE e.path = '/api/v1/produk/:id/penyesuaian-stok';

DELETE FROM api_endpoints WHERE path = '/api/v1/produk/:id/penyesuaian-stok';
