DROP TRIGGER IF EXISTS trigger_user_subscriptions_updated_at ON user_subscriptions;
DROP FUNCTION IF EXISTS update_user_subscriptions_updated_at();
DROP TABLE IF EXISTS user_subscriptions;