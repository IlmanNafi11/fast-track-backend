CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nama VARCHAR(100) UNIQUE NOT NULL,
    kategori VARCHAR(20) NOT NULL CHECK (kategori IN ('admin', 'aplikasi')),
    deskripsi VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_permissions_nama ON permissions(nama);
CREATE INDEX IF NOT EXISTS idx_permissions_kategori ON permissions(kategori);