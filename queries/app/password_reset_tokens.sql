-- Password reset token queries

-- name: CreatePasswordResetToken :one
INSERT INTO password_reset_tokens (email, token, expires_at)
VALUES ($1, $2, $3)
RETURNING id, email, token, expires_at, is_used, created_at;

-- name: GetPasswordResetToken :one
SELECT id, email, token, expires_at, is_used, created_at
FROM password_reset_tokens
WHERE token = $1 AND is_used = false AND expires_at > CURRENT_TIMESTAMP;

-- name: MarkPasswordResetTokenAsUsed :exec
UPDATE password_reset_tokens 
SET is_used = true
WHERE token = $1;

-- name: CleanupExpiredPasswordResetTokens :exec
DELETE FROM password_reset_tokens 
WHERE expires_at < CURRENT_TIMESTAMP OR is_used = true;