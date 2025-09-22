CREATE TABLE IF NOT EXISTS anggarans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kantong_id UUID NOT NULL,
    user_id INTEGER NOT NULL,
    bulan INTEGER NOT NULL CHECK (bulan >= 1 AND bulan <= 12),
    tahun INTEGER NOT NULL CHECK (tahun >= 2020),
    rencana DECIMAL(15,2),
    carry_in DECIMAL(15,2) NOT NULL DEFAULT 0,
    penyesuaian DECIMAL(15,2) NOT NULL DEFAULT 0,
    terpakai DECIMAL(15,2) NOT NULL DEFAULT 0,
    sisa DECIMAL(15,2) NOT NULL DEFAULT 0,
    progres DECIMAL(5,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (kantong_id) REFERENCES kantongs(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    
    UNIQUE(kantong_id, bulan, tahun)
);

CREATE TABLE IF NOT EXISTS penyesuaian_anggarans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    anggaran_id UUID NOT NULL,
    jenis VARCHAR(10) NOT NULL CHECK (jenis IN ('tambah', 'kurangi')),
    jumlah DECIMAL(15,2) NOT NULL CHECK (jumlah >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (anggaran_id) REFERENCES anggarans(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_anggarans_kantong_user ON anggarans(kantong_id, user_id);
CREATE INDEX IF NOT EXISTS idx_anggarans_user_bulan_tahun ON anggarans(user_id, bulan, tahun);
CREATE INDEX IF NOT EXISTS idx_penyesuaian_anggarans_anggaran ON penyesuaian_anggarans(anggaran_id);

CREATE OR REPLACE FUNCTION update_anggaran_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE OR REPLACE TRIGGER update_anggarans_updated_at
    BEFORE UPDATE ON anggarans
    FOR EACH ROW
    EXECUTE FUNCTION update_anggaran_updated_at();