CREATE TABLE IF NOT EXISTS invoices (
    id VARCHAR(50) PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    jumlah DECIMAL(15,2) NOT NULL CHECK (jumlah > 0),
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('sukses', 'gagal', 'pending')),
    dibayar_pada TIMESTAMP NULL,
    metode_pembayaran VARCHAR(100) NULL,
    keterangan TEXT NULL,
    subscription_plan_id CHAR(36) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX idx_invoices_user_id ON invoices(user_id);
CREATE INDEX idx_invoices_status ON invoices(status);
CREATE INDEX idx_invoices_dibayar_pada ON invoices(dibayar_pada);
CREATE INDEX idx_invoices_subscription_plan_id ON invoices(subscription_plan_id);
CREATE INDEX idx_invoices_created_at ON invoices(created_at);

-- Add foreign key constraints
ALTER TABLE invoices 
ADD CONSTRAINT fk_invoices_user_id 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE invoices 
ADD CONSTRAINT fk_invoices_subscription_plan_id 
FOREIGN KEY (subscription_plan_id) REFERENCES subscription_plans(id) ON DELETE SET NULL ON UPDATE CASCADE;

-- Create trigger for updated_at
DELIMITER //
CREATE TRIGGER invoices_updated_at 
    BEFORE UPDATE ON invoices 
    FOR EACH ROW 
BEGIN 
    SET NEW.updated_at = CURRENT_TIMESTAMP; 
END//
DELIMITER ;