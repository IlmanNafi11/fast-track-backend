-- name: CreateTransaksi :exec
INSERT INTO transaksis (user_id, kantong_id, tanggal, jenis, jumlah, catatan)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetTransaksiByID :one
SELECT 
    t.id,
    t.user_id,
    t.kantong_id,
    t.tanggal,
    t.jenis,
    t.jumlah,
    t.catatan,
    t.created_at,
    t.updated_at,
    k.nama as kantong_nama
FROM transaksis t
LEFT JOIN kantongs k ON t.kantong_id = k.id
WHERE t.id = $1 AND t.user_id = $2;

-- name: GetTransaksiListByUserID :many
SELECT 
    t.id,
    t.user_id,
    t.kantong_id,
    t.tanggal,
    t.jenis,
    t.jumlah,
    t.catatan,
    t.created_at,
    t.updated_at,
    k.nama as kantong_nama
FROM transaksis t
LEFT JOIN kantongs k ON t.kantong_id = k.id
WHERE t.user_id = $1
ORDER BY t.tanggal DESC;

-- name: SearchTransaksi :many
SELECT 
    t.id,
    t.user_id,
    t.kantong_id,
    t.tanggal,
    t.jenis,
    t.jumlah,
    t.catatan,
    t.created_at,
    t.updated_at,
    k.nama as kantong_nama
FROM transaksis t
LEFT JOIN kantongs k ON t.kantong_id = k.id
WHERE t.user_id = $1
  AND ($2::text IS NULL OR k.nama ILIKE '%' || $2 || '%' OR t.catatan ILIKE '%' || $2 || '%')
  AND ($3::text IS NULL OR t.jenis = $3)
  AND ($4::text IS NULL OR k.nama ILIKE '%' || $4 || '%')
  AND ($5::date IS NULL OR t.tanggal >= $5)
  AND ($6::date IS NULL OR t.tanggal <= $6)
ORDER BY 
  CASE 
    WHEN $7 = 'tanggal' AND $8 = 'asc' THEN t.tanggal
  END ASC,
  CASE 
    WHEN $7 = 'tanggal' AND $8 = 'desc' THEN t.tanggal
  END DESC,
  CASE 
    WHEN $7 = 'jumlah' AND $8 = 'asc' THEN t.jumlah
  END ASC,
  CASE 
    WHEN $7 = 'jumlah' AND $8 = 'desc' THEN t.jumlah
  END DESC
LIMIT $9 OFFSET $10;

-- name: CountTransaksiByUserID :one
SELECT COUNT(*)
FROM transaksis t
LEFT JOIN kantongs k ON t.kantong_id = k.id
WHERE t.user_id = $1
  AND ($2::text IS NULL OR k.nama ILIKE '%' || $2 || '%' OR t.catatan ILIKE '%' || $2 || '%')
  AND ($3::text IS NULL OR t.jenis = $3)
  AND ($4::text IS NULL OR k.nama ILIKE '%' || $4 || '%')
  AND ($5::date IS NULL OR t.tanggal >= $5)
  AND ($6::date IS NULL OR t.tanggal <= $6);

-- name: UpdateTransaksi :exec
UPDATE transaksis
SET kantong_id = $3,
    tanggal = $4,
    jenis = $5,
    jumlah = $6,
    catatan = $7,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND user_id = $2;

-- name: DeleteTransaksi :exec
DELETE FROM transaksis
WHERE id = $1 AND user_id = $2;

-- name: GetTransaksiSumByKantongAndJenis :one
SELECT COALESCE(SUM(jumlah), 0) as total
FROM transaksis
WHERE kantong_id = $1 AND jenis = $2;

-- name: GetMonthlyTransaksiSummary :many
SELECT 
    DATE_TRUNC('month', tanggal) as bulan,
    jenis,
    SUM(jumlah) as total_jumlah,
    COUNT(*) as jumlah_transaksi
FROM transaksis
WHERE user_id = $1
  AND tanggal >= $2
  AND tanggal <= $3
GROUP BY DATE_TRUNC('month', tanggal), jenis
ORDER BY bulan DESC, jenis;