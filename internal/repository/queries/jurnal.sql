-- name: GetAkunIDByKode :one
SELECT id FROM akun WHERE kode = ?;

-- name: CreateJurnal :execlastid
INSERT INTO jurnal (user_id, keterangan, sumber_tipe, sumber_id)
VALUES (?, ?, ?, ?);

-- name: CreateJurnalDetail :exec
INSERT INTO jurnal_detail (jurnal_id, akun_id, debit, kredit)
VALUES (?, ?, ?, ?);
