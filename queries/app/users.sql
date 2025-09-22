-- User queries

-- name: GetUserByEmail :one
SELECT id, email, password, name, is_active, created_at, updated_at
FROM users 
WHERE email = $1 AND is_active = true;

-- name: GetUserByID :one
SELECT id, email, password, name, is_active, created_at, updated_at
FROM users 
WHERE id = $1 AND is_active = true;

-- name: CreateUser :one
INSERT INTO users (email, password, name, is_active)
VALUES ($1, $2, $3, $4)
RETURNING id, email, name, is_active, created_at, updated_at;

-- name: UpdateUserPassword :exec
UPDATE users 
SET password = $2, updated_at = CURRENT_TIMESTAMP
WHERE email = $1;

-- name: DeactivateUser :exec
UPDATE users 
SET is_active = false, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;