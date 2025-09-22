CREATE TABLE IF NOT EXISTS transaksis (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    kantong_id UUID NOT NULL REFERENCES kantongs(id) ON DELETE CASCADE,
    tanggal DATE NOT NULL,
    jenis VARCHAR(20) NOT NULL CHECK (jenis IN ('Pemasukan', 'Pengeluaran')),
    jumlah DECIMAL(15,2) NOT NULL CHECK (jumlah > 0),
    catatan VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_transaksis_user_id ON transaksis(user_id);
CREATE INDEX IF NOT EXISTS idx_transaksis_kantong_id ON transaksis(kantong_id);
CREATE INDEX IF NOT EXISTS idx_transaksis_tanggal ON transaksis(tanggal);
CREATE INDEX IF NOT EXISTS idx_transaksis_jenis ON transaksis(jenis);
CREATE INDEX IF NOT EXISTS idx_transaksis_user_tanggal ON transaksis(user_id, tanggal);
CREATE INDEX IF NOT EXISTS idx_transaksis_user_jenis ON transaksis(user_id, jenis);