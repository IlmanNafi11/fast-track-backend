CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS kantongs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_kartu VARCHAR(6) UNIQUE NOT NULL,
    user_id INTEGER NOT NULL,
    nama VARCHAR(100) NOT NULL,
    kategori VARCHAR(20) NOT NULL CHECK (kategori IN ('Pengeluaran','Tabungan','Darurat','Transport','Tidak Spesifik')),
    deskripsi VARCHAR(500),
    limit_amount DECIMAL(15,2) CHECK (limit_amount >= 0),
    saldo DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (saldo >= 0),
    warna VARCHAR(10) NOT NULL CHECK (warna IN ('Navy','Glass','Purple','Green','Red')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_kantongs_user_id ON kantongs(user_id);
CREATE INDEX IF NOT EXISTS idx_kantongs_nama ON kantongs(nama);
CREATE INDEX IF NOT EXISTS idx_kantongs_id_kartu ON kantongs(id_kartu);
CREATE INDEX IF NOT EXISTS idx_kantongs_kategori ON kantongs(kategori);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_kantongs_updated_at BEFORE UPDATE
    ON kantongs FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();