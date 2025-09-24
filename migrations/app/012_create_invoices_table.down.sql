DROP TRIGGER IF EXISTS trigger_invoices_updated_at ON invoices;
DROP FUNCTION IF EXISTS update_invoices_updated_at();
DROP INDEX IF EXISTS idx_invoices_created_at;
DROP INDEX IF EXISTS idx_invoices_subscription_plan_id;
DROP INDEX IF EXISTS idx_invoices_dibayar_pada;
DROP INDEX IF EXISTS idx_invoices_status;
DROP INDEX IF EXISTS idx_invoices_user_id;
DROP TABLE IF EXISTS invoices;