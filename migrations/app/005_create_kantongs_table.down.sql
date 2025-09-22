DROP TRIGGER IF EXISTS update_kantongs_updated_at ON kantongs;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP INDEX IF EXISTS idx_kantongs_kategori;
DROP INDEX IF EXISTS idx_kantongs_id_kartu;
DROP INDEX IF EXISTS idx_kantongs_nama;
DROP INDEX IF EXISTS idx_kantongs_user_id;
DROP TABLE IF EXISTS kantongs;