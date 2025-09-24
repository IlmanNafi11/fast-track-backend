DROP TRIGGER IF EXISTS invoices_updated_at;
DROP INDEX IF EXISTS idx_invoices_created_at ON invoices;
DROP INDEX IF EXISTS idx_invoices_subscription_plan_id ON invoices;
DROP INDEX IF EXISTS idx_invoices_dibayar_pada ON invoices;
DROP INDEX IF EXISTS idx_invoices_status ON invoices;
DROP INDEX IF EXISTS idx_invoices_user_id ON invoices;
DROP TABLE IF EXISTS invoices;