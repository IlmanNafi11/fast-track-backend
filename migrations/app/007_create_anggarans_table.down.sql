DROP TRIGGER IF EXISTS update_anggarans_updated_at ON anggarans;
DROP FUNCTION IF EXISTS update_anggaran_updated_at();
DROP INDEX IF EXISTS idx_penyesuaian_anggarans_anggaran;
DROP INDEX IF EXISTS idx_anggarans_user_bulan_tahun;
DROP INDEX IF EXISTS idx_anggarans_kantong_user;
DROP TABLE IF EXISTS penyesuaian_anggarans;
DROP TABLE IF EXISTS anggarans;