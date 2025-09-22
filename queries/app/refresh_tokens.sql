-- Refresh token queries

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token, expires_at)
VALUES ($1, $2, $3)
RETURNING id, user_id, token, expires_at, is_revoked, created_at, updated_at;

-- name: GetRefreshToken :one
SELECT id, user_id, token, expires_at, is_revoked, created_at, updated_at
FROM refresh_tokens
WHERE token = $1 AND is_revoked = false AND expires_at > CURRENT_TIMESTAMP;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens 
SET is_revoked = true, updated_at = CURRENT_TIMESTAMP
WHERE token = $1;

-- name: RevokeAllUserRefreshTokens :exec
UPDATE refresh_tokens 
SET is_revoked = true, updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: CleanupExpiredRefreshTokens :exec
DELETE FROM refresh_tokens 
WHERE expires_at < CURRENT_TIMESTAMP OR is_revoked = true;